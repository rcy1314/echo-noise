package controllers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/storage"
	"github.com/lin-snow/ech0/internal/syncmanager"
)

func isAdmin(c *gin.Context) bool {
	session := sessions.Default(c)
	isAdmin := session.Get("is_admin")
	if isAdmin == nil {
		return false
	}
	return isAdmin.(bool)
}

// 通过预签名URL上传备份到云存储（R2/S3）
func HandleBackupUploadToURL(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
		return
	}

	var req struct {
		UploadURL string `json:"uploadURL" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.UploadURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "上传URL无效"})
		return
	}
	u, err := url.Parse(req.UploadURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "上传URL格式错误"})
		return
	}

	tempDir := fmt.Sprintf("/tmp/ech0_backup_%s", time.Now().Format("20060102150405"))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建临时目录失败"})
		return
	}
	defer os.RemoveAll(tempDir)

	if err := backupImages(tempDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份图片失败: " + err.Error()})
		return
	}
	if err := backupVideos(tempDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份视频失败: " + err.Error()})
		return
	}
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}
	switch dbType {
	case "postgres":
		if err := backupPostgres(tempDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "PostgreSQL备份失败: " + err.Error()})
			return
		}
	case "mysql":
		if err := backupMySQL(tempDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "MySQL备份失败: " + err.Error()})
			return
		}
	default:
		if err := backupSQLite(tempDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "SQLite备份失败: " + err.Error()})
			return
		}
	}
	zipFile := filepath.Join(tempDir, "backup.zip")
	if err := createBackupZip(tempDir, zipFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建备份文件失败: " + err.Error()})
		return
	}

	f, err := os.Open(zipFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "读取备份文件失败"})
		return
	}
	defer f.Close()
	stat, _ := f.Stat()
	reqHttp, err := http.NewRequest(http.MethodPut, req.UploadURL, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建请求失败"})
		return
	}
	reqHttp.Header.Set("Content-Type", "application/zip")
	reqHttp.ContentLength = stat.Size()
	client := &http.Client{Timeout: 120 * time.Second}
	if u.Host != "" {
		if host, _, _ := net.SplitHostPort(u.Host); host != "" {
			_ = host
		}
	}
	resp, err := client.Do(reqHttp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "上传失败: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "云备份上传成功"})
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": fmt.Sprintf("上传失败(%d): %s", resp.StatusCode, string(bodyBytes))})
}

// 通过预签名URL从云存储恢复
func HandleBackupRestoreFromURL(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
		return
	}
	var req struct {
		DownloadURL string `json:"downloadURL" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.DownloadURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "下载URL无效"})
		return
	}
	u, err := url.Parse(req.DownloadURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "下载URL格式错误"})
		return
	}
	tempDir := fmt.Sprintf("/tmp/ech0_restore_%s", time.Now().Format("20060102150405"))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建临时目录失败"})
		return
	}
	defer os.RemoveAll(tempDir)

	resp, err := http.Get(req.DownloadURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "下载失败: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": fmt.Sprintf("下载失败(%d)", resp.StatusCode)})
		return
	}
	restorePath := filepath.Join(tempDir, "cloud_backup.zip")
	out, err := os.Create(restorePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "保存备份文件失败"})
		return
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "写入备份文件失败"})
		return
	}
	out.Close()

	if err := unzipBackup(restorePath, tempDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "解压失败"})
		return
	}
	if err := backupCurrentImages(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份当前图片失败"})
		return
	}
	if err := backupCurrentVideos(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份当前视频失败"})
		return
	}
	if err := restoreImages(tempDir); err != nil {
		restoreCurrentImages()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "恢复图片失败"})
		return
	}
	if err := restoreVideos(tempDir); err != nil {
		restoreCurrentVideos()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "恢复视频失败"})
		return
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}
	switch dbType {
	case "postgres":
		if err := restorePostgres(tempDir); err != nil {
			restoreCurrentImages()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "PostgreSQL恢复失败: " + err.Error()})
			return
		}
	case "mysql":
		if err := restoreMySQL(tempDir); err != nil {
			restoreCurrentImages()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "MySQL恢复失败: " + err.Error()})
			return
		}
	default:
		if err := restoreSQLite(tempDir); err != nil {
			restoreCurrentImages()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "SQLite恢复失败: " + err.Error()})
			return
		}
	}
	if err := database.ReconnectDB(); err != nil {
		restoreCurrentImages()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "数据库重连失败"})
		return
	}
	cleanupImageBackup()
	cleanupVideoBackup()
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "云备份恢复成功", "shouldRefresh": true})
}

func HandleBackupDownload(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
		return
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	tempDir := fmt.Sprintf("/tmp/ech0_backup_%s", time.Now().Format("20060102150405"))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建临时目录失败"})
		return
	}
	defer os.RemoveAll(tempDir)

	// 备份图片文件
	if err := backupImages(tempDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份图片失败: " + err.Error()})
		return
	}
	// 备份视频文件
	if err := backupVideos(tempDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份视频失败: " + err.Error()})
		return
	}
	// 根据数据库类型执行不同的备份逻辑
	switch dbType {
	case "postgres":
		if err := backupPostgres(tempDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "PostgreSQL备份失败: " + err.Error()})
			return
		}
	case "mysql":
		if err := backupMySQL(tempDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "MySQL备份失败: " + err.Error()})
			return
		}
	default:
		if err := backupSQLite(tempDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "SQLite备份失败: " + err.Error()})
			return
		}
	}

	// 创建zip文件
	zipFile := filepath.Join(tempDir, "backup.zip")
	if err := createBackupZip(tempDir, zipFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建备份文件失败: " + err.Error()})
		return
	}

	// 设置响应头并发送文件
	backupName := fmt.Sprintf("ech0_backup_%s_%s.zip", dbType, time.Now().Format("20060102150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+backupName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(zipFile)
}

func HandleBackupRestore(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
		return
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	file, err := c.FormFile("database")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "请选择有效的备份文件"})
		return
	}

	// 检查文件大小
	if file.Size > 500*1024*1024 { // 500MB
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "备份文件过大"})
		return
	}

	tempDir := fmt.Sprintf("/tmp/ech0_restore_%s", time.Now().Format("20060102150405"))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "创建临时目录失败"})
		return
	}
	defer os.RemoveAll(tempDir)

	// 保存并解压备份文件
	backupPath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, backupPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "保存备份文件失败"})
		return
	}

	if err := unzipBackup(backupPath, tempDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "解压备份文件失败"})
		return
	}

	// 备份当前图片
	if err := backupCurrentImages(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份当前图片失败"})
		return
	}
	// 备份当前视频
	if err := backupCurrentVideos(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "备份当前视频失败"})
		return
	}

	// 恢复图片
	if err := restoreImages(tempDir); err != nil {
		restoreCurrentImages() // 尝试恢复原图片
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "恢复图片失败"})
		return
	}
	// 恢复视频
	if err := restoreVideos(tempDir); err != nil {
		restoreCurrentVideos() // 尝试恢复原视频
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "恢复视频失败"})
		return
	}
	// 根据数据库类型执行不同的恢复逻辑
	switch dbType {
	case "postgres":
		if err := restorePostgres(tempDir); err != nil {
			restoreCurrentImages() // 恢复原图片
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "PostgreSQL恢复失败: " + err.Error()})
			return
		}
	case "mysql":
		if err := restoreMySQL(tempDir); err != nil {
			restoreCurrentImages() // 恢复原图片
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "MySQL恢复失败: " + err.Error()})
			return
		}
	default:
		if err := restoreSQLite(tempDir); err != nil {
			restoreCurrentImages() // 恢复原图片
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "SQLite恢复失败: " + err.Error()})
			return
		}
	}

	// 重连数据库
	if err := database.ReconnectDB(); err != nil {
		restoreCurrentImages() // 恢复原图片
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "数据库重连失败"})
		return
	}

	// 清理原图片备份
	cleanupImageBackup()
	// 清理原视频备份
	cleanupVideoBackup()

	c.JSON(http.StatusOK, gin.H{
		"code":          1,
		"msg":           "数据恢复成功",
		"shouldRefresh": true,
	})
}

func HandleBackupPresignUpload(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
		return
	}
	var req struct {
		ObjectKey      string `json:"objectKey"`
		ExpiresSeconds int    `json:"expiresSeconds"`
		ContentType    string `json:"contentType"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "请求参数错误"})
		return
	}
	if req.ObjectKey == "" {
		req.ObjectKey = "backup.zip"
	}
	if req.ContentType == "" {
		req.ContentType = "application/zip"
	}
	if req.ExpiresSeconds <= 0 {
		req.ExpiresSeconds = 3600
	}

	db, _ := database.GetDB()
	var cfg models.SiteConfig
	if err := db.Table("site_configs").First(&cfg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "读取站点配置失败"})
		return
	}
	if !cfg.StorageEnabled || cfg.StorageBucket == "" || cfg.StorageAccessKey == "" || cfg.StorageSecretKey == "" || cfg.StorageEndpoint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "云存储未启用或配置不完整"})
		return
	}
	url, err := storage.PresignUpload(cfg, cfg.StorageBucket, req.ObjectKey, time.Duration(req.ExpiresSeconds)*time.Second, req.ContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "生成预签名上传URL失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "data": gin.H{"url": url}})
}

func HandleBackupPresignDownload(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
		return
	}
	var req struct {
		ObjectKey      string `json:"objectKey"`
		ExpiresSeconds int    `json:"expiresSeconds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "请求参数错误"})
		return
	}
	if req.ObjectKey == "" {
		req.ObjectKey = "backup.zip"
	}
	if req.ExpiresSeconds <= 0 {
		req.ExpiresSeconds = 3600
	}

	db, _ := database.GetDB()
	var cfg models.SiteConfig
	if err := db.Table("site_configs").First(&cfg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "读取站点配置失败"})
		return
	}
	if !cfg.StorageEnabled || cfg.StorageBucket == "" || cfg.StorageAccessKey == "" || cfg.StorageSecretKey == "" || cfg.StorageEndpoint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "云存储未启用或配置不完整"})
		return
	}
	url, err := storage.PresignDownload(cfg, cfg.StorageBucket, req.ObjectKey, time.Duration(req.ExpiresSeconds)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "生成预签名下载URL失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "data": gin.H{"url": url}})
}

// 立即执行云端同步（备份到 R2/S3）
func HandleBackupSyncNow(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "需要管理员权限"})
		return
	}
	db, _ := database.GetDB()
	var cfg models.SiteConfig
	if err := db.Table("site_configs").First(&cfg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "读取站点配置失败"})
		return
	}
	if !cfg.StorageEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "云存储未启用"})
		return
	}
	if err := syncmanager.SyncNow(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "同步失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "已同步到云端"})
}
func backupPostgres(tempDir string) error {
	dumpFile := filepath.Join(tempDir, "database.sql")

	args := []string{
		"-h", os.Getenv("DB_HOST"),
		"-p", os.Getenv("DB_PORT"),
		"-U", os.Getenv("DB_USER"),
		"-d", os.Getenv("DB_NAME"),
		"-f", dumpFile,
		"--no-owner",
		"--no-privileges",
		"--no-password",
		"--clean",
		"--if-exists",
		"--no-tablespaces",
		"--schema=public",      // 只备份 public schema
		"--no-comments",        // 跳过注释
		"--no-publications",    // 跳过发布
		"--no-subscriptions",   // 跳过订阅
		"--no-security-labels", // 跳过安全标签
	}

	cmd := exec.Command("pg_dump", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", os.Getenv("DB_PASSWORD")))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump 执行失败: %v, 错误输出: %s", err, stderr.String())
	}

	return nil
}
func backupMySQL(tempDir string) error {
	dumpFile := filepath.Join(tempDir, "database.sql")

	args := []string{
		"-h", os.Getenv("DB_HOST"),
		"-P", os.Getenv("DB_PORT"),
		"-u", os.Getenv("DB_USER"),
		fmt.Sprintf("-p%s", os.Getenv("DB_PASSWORD")),
		"--skip-opt",             // 禁用所有优化选项
		"--skip-comments",        // 跳过注释
		"--skip-triggers",        // 跳过触发器
		"--skip-extended-insert", // 单行插入
		"--compact",              // 最简输出
		"--skip-ssl",
		os.Getenv("DB_NAME"),
	}

	// 先尝试使用 mariadb-dump
	cmd := exec.Command("mariadb-dump", args...)
	outFile, err := os.Create(dumpFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// 如果失败，尝试最基础的备份方式
		args = []string{
			"-h", os.Getenv("DB_HOST"),
			"-P", os.Getenv("DB_PORT"),
			"-u", os.Getenv("DB_USER"),
			fmt.Sprintf("-p%s", os.Getenv("DB_PASSWORD")),
			"--skip-opt",
			"--compact",
			"--skip-ssl",
			os.Getenv("DB_NAME"),
		}
		cmd = exec.Command("mariadb-dump", args...)
		outFile.Seek(0, 0)
		outFile.Truncate(0)
		cmd.Stdout = outFile
		cmd.Stderr = &stderr
		err = cmd.Run()

		if err != nil {
			// 最后尝试使用 mysqldump
			cmd = exec.Command("mysqldump", args...)
			outFile.Seek(0, 0)
			outFile.Truncate(0)
			cmd.Stdout = outFile
			cmd.Stderr = &stderr
			err = cmd.Run()

			if err != nil {
				return fmt.Errorf("数据库备份失败: %v, 错误输出: %s", err, stderr.String())
			}
		}
	}

	return nil
}
func backupSQLite(tempDir string) error {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	cfgPath := strings.TrimSpace(config.Config.Database.Path)
	if cfgPath == "" {
		cfgPath = "data/noise.db"
	}
	dbEnv := strings.TrimSpace(os.Getenv("DB_PATH"))
	candidates := []string{
		dbEnv,
		cfgPath,
		filepath.Join(wd, cfgPath),
		filepath.Join(exeDir, cfgPath),
		"/app/" + cfgPath,
		"/app/data/noise.db",
	}
	var dbPath string
	for _, p := range candidates {
		if p == "" {
			continue
		}
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			dbPath = p
			break
		}
	}
	if dbPath == "" {
		dbPath = filepath.Join(wd, cfgPath)
	}
	return copyFile(dbPath, filepath.Join(tempDir, "database.db"))
}
func restorePostgres(tempDir string) error {
	dumpFile := filepath.Join(tempDir, "database.sql")

	// 先清理现有连接并重建数据库 - 分成单独的命令执行
	terminateCmd := exec.Command("psql",
		"-h", os.Getenv("DB_HOST"),
		"-p", os.Getenv("DB_PORT"),
		"-U", os.Getenv("DB_USER"),
		"-d", "postgres",
		"-c", fmt.Sprintf(`
            SELECT pg_terminate_backend(pg_stat_activity.pid) 
            FROM pg_stat_activity 
            WHERE pg_stat_activity.datname = '%s' 
            AND pg_stat_activity.pid <> pg_backend_pid();
        `, os.Getenv("DB_NAME")),
	)
	terminateCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", os.Getenv("DB_PASSWORD")))

	var terminateStderr bytes.Buffer
	terminateCmd.Stderr = &terminateStderr

	if err := terminateCmd.Run(); err != nil {
		return fmt.Errorf("终止数据库连接失败: %v, 错误输出: %s", err, terminateStderr.String())
	}

	// 单独执行DROP和CREATE命令
	recreateCmd := exec.Command("psql",
		"-h", os.Getenv("DB_HOST"),
		"-p", os.Getenv("DB_PORT"),
		"-U", os.Getenv("DB_USER"),
		"-d", "postgres",
		"-c", fmt.Sprintf("DROP DATABASE IF EXISTS %s; CREATE DATABASE %s WITH ENCODING='UTF8';",
			os.Getenv("DB_NAME"), os.Getenv("DB_NAME")),
	)
	recreateCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", os.Getenv("DB_PASSWORD")))

	var recreateStderr bytes.Buffer
	recreateCmd.Stderr = &recreateStderr

	if err := recreateCmd.Run(); err != nil {
		return fmt.Errorf("重建数据库失败: %v, 错误输出: %s", err, recreateStderr.String())
	}

	// 恢复数据
	args := []string{
		"-h", os.Getenv("DB_HOST"),
		"-p", os.Getenv("DB_PORT"),
		"-U", os.Getenv("DB_USER"),
		"-d", os.Getenv("DB_NAME"),
		"-f", dumpFile,
		"--single-transaction",
		"--no-owner",
		"--no-privileges",
	}

	cmd := exec.Command("psql", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", os.Getenv("DB_PASSWORD")))

	var restoreStderr bytes.Buffer
	cmd.Stderr = &restoreStderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql 执行失败: %v, 错误输出: %s", err, restoreStderr.String())
	}

	return nil
}
func restoreMySQL(tempDir string) error {
	dumpFile := filepath.Join(tempDir, "database.sql")

	// 先尝试使用 mariadb
	client := "mariadb"
	args := []string{
		"-h", os.Getenv("DB_HOST"),
		"-P", os.Getenv("DB_PORT"),
		"-u", os.Getenv("DB_USER"),
		fmt.Sprintf("-p%s", os.Getenv("DB_PASSWORD")),
		"--skip-ssl",
	}

	// 测试连接
	testCmd := exec.Command(client, append(args, "-e", "SELECT 1")...)
	if err := testCmd.Run(); err != nil {
		// 如果 mariadb 失败，尝试 mysql
		client = "mysql"
		testCmd = exec.Command(client, append(args, "-e", "SELECT 1")...)
		if err := testCmd.Run(); err != nil {
			return fmt.Errorf("数据库连接失败: %v", err)
		}
	}

	// 重置数据库
	resetCmd := exec.Command(client, append(args,
		"-e", fmt.Sprintf("DROP DATABASE IF EXISTS %s; CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;",
			os.Getenv("DB_NAME"), os.Getenv("DB_NAME")))...)
	if err := resetCmd.Run(); err != nil {
		return fmt.Errorf("重置数据库失败: %v", err)
	}

	// 恢复数据
	restoreArgs := append(args, os.Getenv("DB_NAME"))
	cmd := exec.Command(client, restoreArgs...)
	inFile, err := os.Open(dumpFile)
	if err != nil {
		return fmt.Errorf("打开备份文件失败: %v", err)
	}
	defer inFile.Close()

	cmd.Stdin = inFile
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("数据库恢复失败: %v, 错误输出: %s", err, stderr.String())
	}

	return nil
}
func restoreSQLite(tempDir string) error {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	cfgPath := strings.TrimSpace(config.Config.Database.Path)
	if cfgPath == "" {
		cfgPath = "data/noise.db"
	}
	dbEnv := strings.TrimSpace(os.Getenv("DB_PATH"))

	candidates := []string{
		dbEnv,
		cfgPath,
		filepath.Join(wd, cfgPath),
		filepath.Join(exeDir, cfgPath),
		"/app/" + cfgPath,
		"/app/data/noise.db",
	}
	var dbPath string
	for _, p := range candidates {
		if p == "" {
			continue
		}
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			dbPath = p
			break
		}
	}
	if dbPath == "" {
		dbPath = filepath.Join(wd, cfgPath)
		_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	}

	possibleFiles := []string{
		filepath.Join(tempDir, "database.db"),
		filepath.Join(tempDir, "noise.db"),
		filepath.Join(tempDir, "backup.db"),
	}
	var backupFile string
	for _, file := range possibleFiles {
		if _, err := os.Stat(file); err == nil {
			backupFile = file
			break
		}
	}
	if backupFile == "" {
		return fmt.Errorf("找不到有效的 SQLite 备份文件")
	}

	backupPath := dbPath + ".bak"
	_ = os.MkdirAll(filepath.Dir(backupPath), 0755)
	if err := copyFile(dbPath, backupPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("创建当前数据库备份失败: %v", err)
		}
	}
	if err := copyFile(backupFile, dbPath); err != nil {
		if _, statErr := os.Stat(backupPath); statErr == nil {
			_ = copyFile(backupPath, dbPath)
			_ = os.Remove(backupPath)
		}
		return fmt.Errorf("恢复数据库失败: %v", err)
	}
	_ = os.Remove(backupPath)
	return nil
}

func backupImages(tempDir string) error {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	sp := strings.TrimRight(config.Config.Upload.SavePath, "/")
	if sp == "" {
		sp = "data/images"
	}
	candidates := []string{
		sp,
		"./" + sp,
		filepath.Join(wd, sp),
		filepath.Join(exeDir, sp),
		"./data/images",
		filepath.Join(wd, "data/images"),
		filepath.Join(exeDir, "data/images"),
		"/data/images",
		"/app/data/images",
	}
	var imagesDir string
	for _, d := range candidates {
		if d == "" {
			continue
		}
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			imagesDir = d
			break
		}
	}
	if imagesDir == "" {
		return nil
	}
	destDir := filepath.Join(tempDir, "images")
	return copyDir(imagesDir, destDir)
}

func backupCurrentImages() error {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	sp := strings.TrimRight(config.Config.Upload.SavePath, "/")
	if sp == "" {
		sp = "data/images"
	}
	candidates := []string{
		sp,
		"./" + sp,
		filepath.Join(wd, sp),
		filepath.Join(exeDir, sp),
		"./data/images",
		filepath.Join(wd, "data/images"),
		filepath.Join(exeDir, "data/images"),
		"/data/images",
		"/app/data/images",
	}
	var imagesDir string
	for _, d := range candidates {
		if d == "" {
			continue
		}
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			imagesDir = d
			break
		}
	}
	if imagesDir == "" {
		return nil
	}
	backupDir := imagesDir + "_backup"
	return copyDir(imagesDir, backupDir)
}

func restoreImages(tempDir string) error {
	srcDir := filepath.Join(tempDir, "images")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return nil
	}
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	sp := strings.TrimRight(config.Config.Upload.SavePath, "/")
	if sp == "" {
		sp = "data/images"
	}
	candidates := []string{
		sp,
		"./" + sp,
		filepath.Join(wd, sp),
		filepath.Join(exeDir, sp),
		"./data/images",
		filepath.Join(wd, "data/images"),
		filepath.Join(exeDir, "data/images"),
		"/data/images",
		"/app/data/images",
	}
	var destDir string
	for _, d := range candidates {
		if d == "" {
			continue
		}
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			destDir = d
			break
		}
	}
	if destDir == "" {
		destDir = filepath.Join(wd, sp)
	}
	_ = os.MkdirAll(destDir, 0755)
	if err := os.RemoveAll(destDir); err != nil {
		return err
	}
	return copyDir(srcDir, destDir)
}

func restoreCurrentImages() error {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	sp := strings.TrimRight(config.Config.Upload.SavePath, "/")
	if sp == "" {
		sp = "data/images"
	}
	candidates := []string{
		sp,
		"./" + sp,
		filepath.Join(wd, sp),
		filepath.Join(exeDir, sp),
		"./data/images",
		filepath.Join(wd, "data/images"),
		filepath.Join(exeDir, "data/images"),
		"/data/images",
		"/app/data/images",
	}
	var imagesDir string
	for _, d := range candidates {
		if d == "" {
			continue
		}
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			imagesDir = d
			break
		}
	}
	if imagesDir == "" {
		imagesDir = filepath.Join(wd, sp)
	}
	backupDir := imagesDir + "_backup"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return nil
	}
	if err := os.RemoveAll(imagesDir); err != nil {
		return err
	}
	return copyDir(backupDir, imagesDir)
}

func cleanupImageBackup() {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	sp := strings.TrimRight(config.Config.Upload.SavePath, "/")
	if sp == "" {
		sp = "data/images"
	}
	candidates := []string{
		sp + "_backup",
		"./" + sp + "_backup",
		filepath.Join(wd, sp+"_backup"),
		filepath.Join(exeDir, sp+"_backup"),
		"./data/images_backup",
		filepath.Join(wd, "data/images_backup"),
		filepath.Join(exeDir, "data/images_backup"),
		"/data/images_backup",
		"/app/data/images_backup",
	}
	for _, d := range candidates {
		_ = os.RemoveAll(d)
	}
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func createBackupZip(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == zipPath {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func unzipBackup(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destDir, file.Name)

		// 安全检查：防止 zip slip 漏洞
		if !strings.HasPrefix(path, destDir) {
			return fmt.Errorf("非法的文件路径: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// 添加视频备份相关辅助函数
func backupCurrentVideos() error {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	candidates := []string{
		"./data/video",
		filepath.Join(wd, "data/video"),
		filepath.Join(exeDir, "data/video"),
		"/data/video",
		"/app/data/video",
	}
	var videoDir string
	for _, d := range candidates {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			videoDir = d
			break
		}
	}
	if videoDir == "" {
		return nil
	}
	backupDir := videoDir + "_backup"
	return copyDir(videoDir, backupDir)
}

func restoreCurrentVideos() error {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	candidates := []string{
		"./data/video",
		filepath.Join(wd, "data/video"),
		filepath.Join(exeDir, "data/video"),
		"/data/video",
		"/app/data/video",
	}
	var videoDir string
	for _, d := range candidates {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			videoDir = d
			break
		}
	}
	if videoDir == "" {
		videoDir = filepath.Join(wd, "data/video")
	}
	backupDir := videoDir + "_backup"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return nil
	}
	if err := os.RemoveAll(videoDir); err != nil {
		return err
	}
	return copyDir(backupDir, videoDir)
}

func cleanupVideoBackup() {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	candidates := []string{
		"./data/video_backup",
		filepath.Join(wd, "data/video_backup"),
		filepath.Join(exeDir, "data/video_backup"),
		"/data/video_backup",
		"/app/data/video_backup",
	}
	for _, d := range candidates {
		_ = os.RemoveAll(d)
	}
}
func backupVideos(tempDir string) error {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	candidates := []string{
		"./data/video",
		filepath.Join(wd, "data/video"),
		filepath.Join(exeDir, "data/video"),
		"/data/video",
		"/app/data/video",
	}
	var videoDir string
	for _, d := range candidates {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			videoDir = d
			break
		}
	}
	if videoDir == "" {
		return nil
	}
	destDir := filepath.Join(tempDir, "video")
	return copyDir(videoDir, destDir)
}

func restoreVideos(tempDir string) error {
	srcDir := filepath.Join(tempDir, "video")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return nil
	}
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	candidates := []string{
		"./data/video",
		filepath.Join(wd, "data/video"),
		filepath.Join(exeDir, "data/video"),
		"/data/video",
		"/app/data/video",
	}
	var destDir string
	for _, d := range candidates {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			destDir = d
			break
		}
	}
	if destDir == "" {
		destDir = filepath.Join(wd, "data/video")
	}
	_ = os.MkdirAll(destDir, 0755)
	if err := os.RemoveAll(destDir); err != nil {
		return err
	}
	return copyDir(srcDir, destDir)
}
