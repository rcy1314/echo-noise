package syncmanager

import (
    "bytes"
    "archive/zip"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/lin-snow/ech0/internal/database"
    "github.com/lin-snow/ech0/internal/models"
    "github.com/lin-snow/ech0/internal/storage"
)

var (
    mu sync.Mutex
    configured models.SiteConfig
    scheduledStop chan struct{}
    debounceTimer *time.Timer
)

func Configure(cfg models.SiteConfig) {
    mu.Lock()
    configured = cfg
    if scheduledStop != nil {
        close(scheduledStop)
        scheduledStop = nil
    }
    if cfg.StorageEnabled && cfg.StorageAutoSyncEnabled && cfg.StorageSyncMode == "scheduled" && cfg.StorageSyncIntervalMinute > 0 {
        scheduledStop = make(chan struct{})
        interval := time.Duration(cfg.StorageSyncIntervalMinute) * time.Minute
        go func(stop <-chan struct{}) {
            ticker := time.NewTicker(interval)
            defer ticker.Stop()
            for {
                select {
                case <-ticker.C:
                    _ = SyncNow()
                case <-stop:
                    return
                }
            }
        }(scheduledStop)
    }
    mu.Unlock()
}

func Trigger() {
    mu.Lock()
    defer mu.Unlock()
    if !(configured.StorageEnabled && configured.StorageAutoSyncEnabled && configured.StorageSyncMode == "instant") {
        return
    }
    if debounceTimer != nil {
        debounceTimer.Stop()
    }
    debounceTimer = time.AfterFunc(15*time.Second, func() { _ = SyncNow() })
}

func SyncNow() error {
    mu.Lock()
    cfg := configured
    mu.Unlock()
    if !cfg.StorageEnabled { return nil }

    // 构建临时备份文件
    tmpDir := os.TempDir()
    backupPath := filepath.Join(tmpDir, "backup.zip")
    _ = os.Remove(backupPath)
    f, err := os.Create(backupPath)
    if err != nil { return err }
    zw := zip.NewWriter(f)

    // 数据库文件
    dbType := os.Getenv("DB_TYPE")
    if dbType == "" { dbType = "sqlite" }
    if dbType == "sqlite" {
        dbPath := os.Getenv("DB_PATH")
        if dbPath == "" { dbPath = "/app/data/noise.db" }
        if err := addFileToZip(zw, dbPath, "database.db"); err != nil { /* ignore missing */ }
    }
    // images
    _ = addDirToZip(zw, "./data/images", "images")
    // video
    _ = addDirToZip(zw, "./data/video", "video")

    _ = zw.Close()
    f.Close()

    // 预签名上传
    url, err := storage.PresignUpload(cfg, cfg.StorageBucket, "backup.zip", 1*time.Hour, "application/zip")
    if err != nil { return err }
    // PUT 上传
    content, err := os.ReadFile(backupPath)
    if err != nil { return err }
    req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(content))
    req.Header.Set("Content-Type", "application/zip")
    client := &http.Client{ Timeout: 120 * time.Second }
    resp, err := client.Do(req)
    if err != nil { return err }
    io.Copy(io.Discard, resp.Body)
    resp.Body.Close()

    now := time.Now()
    db, _ := database.GetDB()
    _ = db.Table("site_configs").Where("id = ?", cfg.ID).Updates(map[string]any{ "storage_last_sync_time": &now }).Error
    return nil
}

func addFileToZip(zw *zip.Writer, srcPath, zipName string) error {
    fi, err := os.Stat(srcPath)
    if err != nil || fi.IsDir() { return err }
    r, err := os.Open(srcPath)
    if err != nil { return err }
    defer r.Close()
    w, err := zw.Create(zipName)
    if err != nil { return err }
    _, err = io.Copy(w, r)
    return err
}

func addDirToZip(zw *zip.Writer, dir, prefix string) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil { return nil }
        if info.IsDir() { return nil }
        rel, _ := filepath.Rel(dir, path)
        zipName := filepath.ToSlash(filepath.Join(prefix, rel))
        return addFileToZip(zw, path, zipName)
    })
}

