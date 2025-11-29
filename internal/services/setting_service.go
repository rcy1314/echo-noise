package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/syncmanager"
)

// GetFrontendConfig 获取前端配置
func GetFrontendConfig() (map[string]interface{}, error) {
	db, err := database.GetDB()
	if err != nil {
		return getDefaultConfig(), nil
	}

	var config models.SiteConfig
	if err := db.Table("site_configs").First(&config).Error; err != nil {
		return getDefaultConfig(), nil
	}

	// 新增：读取Setting表的AllowRegistration
	var setting models.Setting
	allowReg := true
	if err := db.Table("settings").First(&setting).Error; err == nil {
		allowReg = setting.AllowRegistration
	}

	// 读取 DB 类型
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	configMap := map[string]interface{}{
		"allowRegistration": allowReg,
		"dbType":            dbType,
		"frontendSettings": map[string]interface{}{
			"siteTitle":        config.SiteTitle,
			"subtitleText":     config.SubtitleText,
			"avatarURL":        config.AvatarURL,
			"username":         config.Username,
			"description":      config.Description,
			"backgrounds":      config.GetBackgroundsList(),
			"cardFooterTitle":  config.CardFooterTitle,
			"cardFooterLink":   config.CardFooterLink,
			"pageFooterHTML":   config.PageFooterHTML,
			"rssTitle":         config.RSSTitle,
			"rssDescription":   config.RSSDescription,
			"rssAuthorName":    config.RSSAuthorName,
			"rssFaviconURL":    config.RSSFaviconURL,
			"walineServerURL":  config.WalineServerURL,
			"enableGithubCard": config.EnableGithubCard,
			// GitHub OAuth
			"githubOAuthEnabled": config.GithubOAuthEnabled,
			"githubClientId":     config.GithubClientId,
			"githubClientSecret": config.GithubClientSecret,
			"githubCallbackURL":  config.GithubCallbackURL,
			// PWA 设置
			"pwaEnabled":     config.PwaEnabled,
			"pwaTitle":       choose(config.PwaTitle, config.SiteTitle),
			"pwaDescription": choose(config.PwaDescription, config.Description),
			"pwaIconURL":     choose(config.PwaIconURL, config.RSSFaviconURL),
			// 默认内容主题
			"defaultContentTheme": choose(config.ContentThemeDefault, "dark"),
			// 公告栏
			"announcementText":    choose(config.AnnouncementText, "欢迎访问我的说说笔记！"),
			"announcementEnabled": config.AnnouncementEnabled,
			// 音乐播放器
			"musicEnabled":          config.MusicEnabled,
			"musicPlaylistId":       choose(config.MusicPlaylistId, ""),
			"musicSongId":           choose(config.MusicSongId, ""),
			"musicPosition":         choose(config.MusicPosition, "bottom-left"),
			"musicTheme":            choose(config.MusicTheme, "auto"),
			"musicLyric":            config.MusicLyric,
			"musicAutoplay":         config.MusicAutoplay,
			"musicDefaultMinimized": config.MusicDefaultMinimized,
			"musicEmbed":            config.MusicEmbed,
			"musicCssCdnURL":        choose(config.MusicCssCdnURL, ""),
			"musicJsCdnURL":         choose(config.MusicJsCdnURL, ""),
			// 评论系统
			"commentEnabled":      config.CommentEnabled,
			"commentSystem":       choose(config.CommentSystem, "waline"),
			"commentEmailEnabled": config.CommentEmailEnabled,
		},
		"storageEnabled": config.StorageEnabled,
		"storageConfig": map[string]interface{}{
			"provider":        choose(config.StorageProvider, ""),
			"endpoint":        choose(config.StorageEndpoint, ""),
			"region":          choose(config.StorageRegion, ""),
			"bucket":          choose(config.StorageBucket, ""),
			"accessKey":       choose(config.StorageAccessKey, ""),
			"secretKey":       choose(config.StorageSecretKey, ""),
			"usePathStyle":    config.StorageUsePathStyle,
			"publicBaseURL":   choose(config.StoragePublicBaseURL, ""),
			"autoSyncEnabled": config.StorageAutoSyncEnabled,
			"syncMode":        choose(config.StorageSyncMode, "instant"),
			"syncIntervalMinute": func() int {
				if config.StorageSyncIntervalMinute > 0 {
					return config.StorageSyncIntervalMinute
				}
				return 15
			}(),
			"lastSyncTime": func() string {
				if config.StorageLastSyncTime != nil {
					return config.StorageLastSyncTime.Format(time.RFC3339)
				}
				return ""
			}(),
		},
		"smtpEnabled":    config.SmtpEnabled,
		"smtpDriver":     config.SmtpDriver,
		"smtpHost":       config.SmtpHost,
		"smtpPort":       config.SmtpPort,
		"smtpUser":       config.SmtpUser,
		"smtpPass":       config.SmtpPass,
		"smtpFrom":       config.SmtpFrom,
		"smtpEncryption": config.SmtpEncryption,
		"smtpTLS":        config.SmtpTLS,
	}
	return configMap, nil
}

// UpdateSetting 更新站点配置
func UpdateFrontendSetting(userID uint, settingMap map[string]interface{}) error {
	db, err := database.GetDB()
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	frontendSettings, ok := settingMap["frontendSettings"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("无效的前端配置格式")
	}

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var config models.SiteConfig
	// 先尝试获取现有配置
	if err := tx.Table("site_configs").First(&config).Error; err != nil {
		config.ID = 1 // 设置默认ID
	}

	// 更新配置字段
	if v, ok := frontendSettings["siteTitle"].(string); ok {
		config.SiteTitle = v
	}
	if v, ok := frontendSettings["subtitleText"].(string); ok {
		config.SubtitleText = v
	}
	if v, ok := frontendSettings["avatarURL"].(string); ok {
		config.AvatarURL = v
	}
	if v, ok := frontendSettings["username"].(string); ok {
		config.Username = v
	}
	if v, ok := frontendSettings["description"].(string); ok {
		config.Description = v
	}
	if v, ok := frontendSettings["cardFooterTitle"].(string); ok {
		config.CardFooterTitle = v
	}
	if v, ok := frontendSettings["cardFooterLink"].(string); ok {
		config.CardFooterLink = v
	}
	if v, ok := frontendSettings["pageFooterHTML"].(string); ok {
		config.PageFooterHTML = v
	}
	if v, ok := frontendSettings["rssTitle"].(string); ok {
		config.RSSTitle = v
	}
	if v, ok := frontendSettings["rssDescription"].(string); ok {
		config.RSSDescription = v
	}
	if v, ok := frontendSettings["rssAuthorName"].(string); ok {
		config.RSSAuthorName = v
	}
	if v, ok := frontendSettings["rssFaviconURL"].(string); ok {
		config.RSSFaviconURL = v
	}
	if v, ok := frontendSettings["walineServerURL"].(string); ok {
		config.WalineServerURL = v
	}
	// 评论系统设置
	if vb, ok := frontendSettings["commentEnabled"].(bool); ok {
		config.CommentEnabled = vb
	} else if vs, ok := frontendSettings["commentEnabled"].(string); ok {
		if vs == "true" {
			config.CommentEnabled = true
		} else if vs == "false" {
			config.CommentEnabled = false
		}
	}

	// 音乐播放器设置
	if vb, ok := frontendSettings["musicEnabled"].(bool); ok {
		config.MusicEnabled = vb
	} else if vs, ok := frontendSettings["musicEnabled"].(string); ok {
		config.MusicEnabled = (vs == "true")
	}
	if v, ok := frontendSettings["musicPlaylistId"].(string); ok {
		config.MusicPlaylistId = v
	}
	if v, ok := frontendSettings["musicSongId"].(string); ok {
		config.MusicSongId = v
	}
	if v, ok := frontendSettings["musicPosition"].(string); ok {
		config.MusicPosition = v
	}
	if v, ok := frontendSettings["musicTheme"].(string); ok {
		config.MusicTheme = v
	}
	if vb, ok := frontendSettings["musicLyric"].(bool); ok {
		config.MusicLyric = vb
	} else if vs, ok := frontendSettings["musicLyric"].(string); ok {
		config.MusicLyric = (vs == "true")
	}
	if vb, ok := frontendSettings["musicAutoplay"].(bool); ok {
		config.MusicAutoplay = vb
	} else if vs, ok := frontendSettings["musicAutoplay"].(string); ok {
		config.MusicAutoplay = (vs == "true")
	}
	if vb, ok := frontendSettings["musicDefaultMinimized"].(bool); ok {
		config.MusicDefaultMinimized = vb
	} else if vs, ok := frontendSettings["musicDefaultMinimized"].(string); ok {
		config.MusicDefaultMinimized = (vs == "true")
	}
	if vb, ok := frontendSettings["musicEmbed"].(bool); ok {
		config.MusicEmbed = vb
	} else if vs, ok := frontendSettings["musicEmbed"].(string); ok {
		config.MusicEmbed = (vs == "true")
	}
	if v, ok := frontendSettings["musicCssCdnURL"].(string); ok {
		config.MusicCssCdnURL = v
	}
	if v, ok := frontendSettings["musicJsCdnURL"].(string); ok {
		config.MusicJsCdnURL = v
	}
	if v, ok := frontendSettings["commentSystem"].(string); ok {
		config.CommentSystem = v
	}
	if vb, ok := frontendSettings["commentEmailEnabled"].(bool); ok {
		config.CommentEmailEnabled = vb
	} else if vs, ok := frontendSettings["commentEmailEnabled"].(string); ok {
		if vs == "true" {
			config.CommentEmailEnabled = true
		} else if vs == "false" {
			config.CommentEmailEnabled = false
		}
	}
	// GitHub OAuth 设置
	if vb, ok := frontendSettings["githubOAuthEnabled"].(bool); ok {
		config.GithubOAuthEnabled = vb
	} else if vs, ok := frontendSettings["githubOAuthEnabled"].(string); ok {
		if vs == "true" {
			config.GithubOAuthEnabled = true
		} else if vs == "false" {
			config.GithubOAuthEnabled = false
		}
	}
	if v, ok := frontendSettings["githubClientId"].(string); ok {
		config.GithubClientId = v
	}
	if v, ok := frontendSettings["githubClientSecret"].(string); ok {
		config.GithubClientSecret = v
	}
	if v, ok := frontendSettings["githubCallbackURL"].(string); ok {
		config.GithubCallbackURL = v
	}
	if v, ok := frontendSettings["enableGithubCard"].(bool); ok {
		config.EnableGithubCard = v
	} else if vs, ok := frontendSettings["enableGithubCard"].(string); ok {
		if vs == "true" {
			config.EnableGithubCard = true
		} else if vs == "false" {
			config.EnableGithubCard = false
		}
	}
	// 公告栏
	if v, ok := frontendSettings["announcementText"].(string); ok {
		config.AnnouncementText = v
	}
	if vb, ok := frontendSettings["announcementEnabled"].(bool); ok {
		config.AnnouncementEnabled = vb
	} else if vs, ok := frontendSettings["announcementEnabled"].(string); ok {
		if vs == "true" {
			config.AnnouncementEnabled = true
		} else if vs == "false" {
			config.AnnouncementEnabled = false
		}
	}
	// PWA 设置
	if v, ok := frontendSettings["pwaEnabled"].(bool); ok {
		config.PwaEnabled = v
	}
	if v, ok := frontendSettings["pwaTitle"].(string); ok {
		config.PwaTitle = v
	}
	if v, ok := frontendSettings["pwaDescription"].(string); ok {
		config.PwaDescription = v
	}
	if v, ok := frontendSettings["pwaIconURL"].(string); ok {
		config.PwaIconURL = v
	}

	// 默认内容主题
	if v, ok := frontendSettings["defaultContentTheme"].(string); ok {
		if v == "dark" || v == "light" {
			config.ContentThemeDefault = v
		}
	}

	// 处理背景图片列表
	if backgrounds, ok := frontendSettings["backgrounds"].([]interface{}); ok {
		backgroundsList := make([]string, 0, len(backgrounds))
		for _, bg := range backgrounds {
			if bgStr, ok := bg.(string); ok && bgStr != "" {
				backgroundsList = append(backgroundsList, bgStr)
			}
		}
		// 确保至少保留一个默认背景
		if len(backgroundsList) == 0 {
			backgroundsList = getDefaultConfig()["frontendSettings"].(map[string]interface{})["backgrounds"].([]string)
		}
		backgroundsJSON, err := json.Marshal(backgroundsList)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("背景图片列表序列化失败: %v", err)
		}
		config.Backgrounds = string(backgroundsJSON)
	} else if backgrounds, ok := frontendSettings["backgrounds"].([]string); ok {
		// 直接处理字符串数组
		if len(backgrounds) == 0 {
			backgrounds = getDefaultConfig()["frontendSettings"].(map[string]interface{})["backgrounds"].([]string)
		}
		backgroundsJSON, err := json.Marshal(backgrounds)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("背景图片列表序列化失败: %v", err)
		}
		config.Backgrounds = string(backgroundsJSON)
	}

	// 保存或更新配置
	if config.ID == 0 {
		if err := tx.Table("site_configs").Create(&config).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建配置失败: %v", err)
		}
	} else {
		if err := tx.Table("site_configs").Save(&config).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新配置失败: %v", err)
		}
	}

	if v, ok := settingMap["storageEnabled"].(bool); ok {
		config.StorageEnabled = v
	}
	if sc, ok := settingMap["storageConfig"].(map[string]interface{}); ok {
		if pv, ok := sc["provider"].(string); ok {
			config.StorageProvider = pv
		}
		if v, ok := sc["endpoint"].(string); ok {
			config.StorageEndpoint = v
		}
		if v, ok := sc["region"].(string); ok {
			config.StorageRegion = v
		}
		if v, ok := sc["bucket"].(string); ok {
			config.StorageBucket = v
		}
		if v, ok := sc["accessKey"].(string); ok {
			config.StorageAccessKey = v
		}
		if v, ok := sc["secretKey"].(string); ok {
			config.StorageSecretKey = v
		}
		if v, ok := sc["usePathStyle"].(bool); ok {
			config.StorageUsePathStyle = v
		}
		if v, ok := sc["publicBaseURL"].(string); ok {
			config.StoragePublicBaseURL = v
		}
		if vb, ok := sc["autoSyncEnabled"].(bool); ok {
			config.StorageAutoSyncEnabled = vb
		} else if vs, ok := sc["autoSyncEnabled"].(string); ok {
			config.StorageAutoSyncEnabled = (vs == "true")
		}
		if v, ok := sc["syncMode"].(string); ok {
			if v == "instant" || v == "scheduled" {
				config.StorageSyncMode = v
			}
		}
		if vi, ok := sc["syncIntervalMinute"].(float64); ok {
			config.StorageSyncIntervalMinute = int(vi)
		} else if vi2, ok := sc["syncIntervalMinute"].(int); ok {
			config.StorageSyncIntervalMinute = vi2
		} else if vs, ok := sc["syncIntervalMinute"].(string); ok {
			if n, err := strconv.Atoi(vs); err == nil {
				config.StorageSyncIntervalMinute = n
			}
		}
	}

	if config.StorageProvider == "r2" {
		config.StorageUsePathStyle = true
	}
	if config.StorageEnabled {
		if config.StorageProvider == "" || config.StorageEndpoint == "" || config.StorageBucket == "" || config.StorageAccessKey == "" || config.StorageSecretKey == "" {
			tx.Rollback()
			return fmt.Errorf("云存储配置不完整")
		}
	}

	if v, ok := settingMap["smtpEnabled"].(bool); ok {
		config.SmtpEnabled = v
	}
	if v, ok := settingMap["smtpDriver"].(string); ok {
		config.SmtpDriver = v
	}
	if v, ok := settingMap["smtpHost"].(string); ok {
		config.SmtpHost = v
	}
	if v, ok := settingMap["smtpPort"].(float64); ok {
		config.SmtpPort = int(v)
	} else if vi, ok := settingMap["smtpPort"].(int); ok {
		config.SmtpPort = vi
	} else if vs, ok := settingMap["smtpPort"].(string); ok {
		if p, err := strconv.Atoi(vs); err == nil {
			config.SmtpPort = p
		}
	}
	if v, ok := settingMap["smtpUser"].(string); ok {
		config.SmtpUser = v
	}
	if v, ok := settingMap["smtpPass"].(string); ok {
		config.SmtpPass = v
	}
	if v, ok := settingMap["smtpFrom"].(string); ok {
		config.SmtpFrom = v
	}
	if v, ok := settingMap["smtpEncryption"].(string); ok {
		config.SmtpEncryption = v
	}
	if v, ok := settingMap["smtpTLS"].(bool); ok {
		config.SmtpTLS = v
	}

	// 自动启用：当必填项齐全时，强制启用
	if !config.SmtpEnabled {
		if config.SmtpHost != "" && config.SmtpPort > 0 && config.SmtpUser != "" && config.SmtpPass != "" &&
			(config.SmtpEncryption == "ssl" || config.SmtpEncryption == "tls") {
			config.SmtpEnabled = true
		}
	}

	// 基础校验：开启时必填项必须完整
	if config.SmtpEnabled {
		if config.SmtpHost == "" || config.SmtpPort <= 0 || config.SmtpUser == "" || config.SmtpPass == "" ||
			(config.SmtpEncryption != "ssl" && config.SmtpEncryption != "tls") {
			tx.Rollback()
			return fmt.Errorf("邮件设置错误")
		}
	}

	if err := tx.Table("site_configs").Save(&config).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新配置失败: %v", err)
	}

	// 同步配置到自动同步管理器
	// 注意：仅在服务进程内触发，不影响数据库事务
	// 读取最新配置并传入管理器
	db.Table("site_configs").First(&config)
	// 调用同步管理器进行配置
	syncmanager.Configure(config)

	if config.StorageEnabled {
		dbType := os.Getenv("DB_TYPE")
		if dbType == "" {
			dbType = "sqlite"
		}
		if dbType == "sqlite" {
			base := strings.TrimSpace(config.StoragePublicBaseURL)
			if base != "" {
				url := strings.TrimRight(base, "/") + "/database.db"
				client := &http.Client{Timeout: 60 * time.Second}
				resp, err := client.Get(url)
				if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
					defer resp.Body.Close()
					tempFile := filepath.Join(os.TempDir(), "cloud_database.db")
					out, err := os.Create(tempFile)
					if err == nil {
						_, _ = io.Copy(out, resp.Body)
						out.Close()
						dbPath := os.Getenv("DB_PATH")
						if dbPath == "" {
							dbPath = "/app/data/noise.db"
						}
						_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
						_ = copyFile(tempFile, dbPath)
						_ = os.Remove(tempFile)
						_ = database.ReconnectDB()
					}
				}
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交配置更新失败: %v", err)
	}

	return nil
}

// 获取默认配置
func getDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"allowRegistration": true,
		"frontendSettings": map[string]interface{}{
			"siteTitle":    "Noise的说说笔记",
			"subtitleText": "欢迎访问，点击头像可更换封面背景！",
			"avatarURL":    "https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png",
			"username":     "Noise",
			"description":  "执迷不悟",
			"backgrounds": []string{
				"https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg",
				"https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg",
				"https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg",
				"https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg",
				"https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg",
				"https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg",
				"https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg",
				"https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg",
				"https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg",
				"https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg",
				"https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg",
				"https://s2.loli.net/2025/03/27/wYy12qDMH6bGJOI.jpg",
				"https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg",
			},
			"cardFooterTitle":  "Noise·说说·笔记~",
			"cardFooterLink":   "note.noisework.cn",
			"pageFooterHTML":   `<div class="text-center text-xs text-gray-400 py-4">来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 使用<a href="https://github.com/rcy1314/echo-noise" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0-Noise</a>发布</div>`,
			"rssTitle":         "Noise的说说笔记",
			"rssDescription":   "一个说说笔记~",
			"rssAuthorName":    "Noise",
			"rssFaviconURL":    "/favicon-32x32.png",
			"walineServerURL":  "请前往waline官网https://waline.js.org查看部署配置",
			"enableGithubCard": false,
			// GitHub OAuth 默认关闭
			"githubOAuthEnabled": false,
			"githubClientId":     "",
			"githubClientSecret": "",
			"githubCallbackURL":  "",
			// PWA 设置默认值
			"pwaEnabled":          true,
			"pwaTitle":            "",
			"pwaDescription":      "",
			"pwaIconURL":          "",
			"defaultContentTheme": "light",
			"announcementText":    "欢迎访问我的说说笔记！",
			"announcementEnabled": true,
			// 音乐播放器默认关闭，但默认位置为左下角并最小化
			"musicEnabled":          false,
			"musicPlaylistId":       "",
			"musicSongId":           "",
			"musicPosition":         "bottom-left",
			"musicTheme":            "auto",
			"musicLyric":            true,
			"musicAutoplay":         false,
			"musicDefaultMinimized": true,
			"musicEmbed":            false,
			"musicCssCdnURL":        "",
			"musicJsCdnURL":         "",
			// 评论系统默认值
			"commentEnabled":                false,
			"commentSystem":                 "waline",
			"commentEmailEnabled":           false,
			"commentEmailReplyName":         "",
			"commentEmailAdminPrefix":       "",
			"commentEmailReplyPrefix":       "",
			"commentEmailReplyTemplate":     "- 您在{site}主页上的内容有了新的评论\n- {nick} 回复说：\n- {content}\n- 您可以点击查看回复的完整内容：{url}",
			"commentEmailAdminTemplate":     "",
			"commentEmailSiteURL":           "",
			"commentEmailReplyTemplateHTML": "",
			"commentEmailAdminTemplateHTML": "",
		},
		"storageEnabled": false,
		"storageConfig": map[string]interface{}{
			"provider":      "",
			"endpoint":      "",
			"region":        "",
			"bucket":        "",
			"accessKey":     "",
			"secretKey":     "",
			"usePathStyle":  true,
			"publicBaseURL": "",
		},
	}
}

// 选择第一个非空字符串
func choose(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
