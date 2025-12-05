package controllers

import (
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/lin-snow/ech0/config"
    "github.com/lin-snow/ech0/internal/database"
    "github.com/lin-snow/ech0/internal/models"
)

type AttachmentInfo struct {
    Name       string        `json:"name"`
    URL        string        `json:"url"`
    Size       int64         `json:"size"`
    ModifiedAt time.Time     `json:"modified_at"`
    Belongs    []BelongItem  `json:"belongs"`
}

type BelongItem struct {
    ID        uint      `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    Snippet   string    `json:"snippet"`
}

func ListImageAttachments(c *gin.Context) {
    wd, _ := os.Getwd()
    exePath, _ := os.Executable()
    exeDir := filepath.Dir(exePath)
    sp := strings.TrimRight(config.Config.Upload.SavePath, "/")
    dir := pickDir([]string{
        sp,
        "./" + sp,
        filepath.Join(wd, sp),
        filepath.Join(exeDir, sp),
        "./data/images",
        filepath.Join(wd, "data/images"),
        filepath.Join(exeDir, "data/images"),
        "/data/images",
        "/app/data/images",
    }, "./data/images")
    entries, err := os.ReadDir(dir)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"code": 1, "data": []AttachmentInfo{}})
        return
    }

    var messages []models.Message
    database.DB.Select("id", "content", "image_url", "created_at").Order("created_at DESC").Find(&messages)

    var list []AttachmentInfo
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        name := e.Name()
        p := filepath.Join(dir, name)
        fi, err := os.Stat(p)
        if err != nil {
            continue
        }
        url := "/api/images/" + name
        belongs := findBelongs(messages, name, "/images/", "/api/images/")
        list = append(list, AttachmentInfo{ Name: name, URL: url, Size: fi.Size(), ModifiedAt: fi.ModTime(), Belongs: belongs })
    }

    c.JSON(http.StatusOK, gin.H{"code": 1, "data": list})
}

func ListVideoAttachments(c *gin.Context) {
    wd, _ := os.Getwd()
    exePath, _ := os.Executable()
    exeDir := filepath.Dir(exePath)
    dir := pickDir([]string{
        "./data/video",
        filepath.Join(wd, "data/video"),
        filepath.Join(exeDir, "data/video"),
        "/data/video",
        "/app/data/video",
    }, "./data/video")
    entries, err := os.ReadDir(dir)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"code": 1, "data": []AttachmentInfo{}})
        return
    }

    var messages []models.Message
    database.DB.Select("id", "content", "image_url", "created_at").Order("created_at DESC").Find(&messages)

    var list []AttachmentInfo
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        name := e.Name()
        p := filepath.Join(dir, name)
        fi, err := os.Stat(p)
        if err != nil {
            continue
        }
        url := "/video/" + name
        belongs := findBelongs(messages, name, "/video/", "/api/video/")
        list = append(list, AttachmentInfo{ Name: name, URL: url, Size: fi.Size(), ModifiedAt: fi.ModTime(), Belongs: belongs })
    }

    c.JSON(http.StatusOK, gin.H{"code": 1, "data": list})
}

func findBelongs(messages []models.Message, name, p1, p2 string) []BelongItem {
    var out []BelongItem
    needle1 := p1 + name
    needle2 := p2 + name
    for _, m := range messages {
        has := false
        if strings.Contains(m.Content, needle1) || strings.Contains(m.Content, needle2) {
            has = true
        }
        if !has {
            if strings.Contains(m.ImageURL, needle1) || strings.Contains(m.ImageURL, needle2) {
                has = true
            }
        }
        if has {
            snip := m.Content
            if len(snip) > 80 {
                snip = snip[:80]
            }
            out = append(out, BelongItem{ ID: m.ID, CreatedAt: m.CreatedAt, Snippet: snip })
        }
    }
    return out
}

func DeleteImageAttachment(c *gin.Context) {
    name := c.Param("name")
    base := filepath.Base(name)
    wd, _ := os.Getwd()
    exePath, _ := os.Executable()
    exeDir := filepath.Dir(exePath)
    sp := strings.TrimRight(config.Config.Upload.SavePath, "/")
    imgDir := pickDir([]string{
        sp,
        "./" + sp,
        filepath.Join(wd, sp),
        filepath.Join(exeDir, sp),
        "./data/images",
        filepath.Join(wd, "data/images"),
        filepath.Join(exeDir, "data/images"),
        "/data/images",
        "/app/data/images",
    }, "./data/images")
    p := filepath.Join(imgDir, base)
    if _, err := os.Stat(p); err != nil {
        c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "文件不存在"})
        return
    }
    if err := os.Remove(p); err != nil {
        c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除失败"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 1, "data": true})
}

func DeleteVideoAttachment(c *gin.Context) {
    name := c.Param("name")
    base := filepath.Base(name)
    wd, _ := os.Getwd()
    exePath, _ := os.Executable()
    exeDir := filepath.Dir(exePath)
    vidDir := pickDir([]string{
        "./data/video",
        filepath.Join(wd, "data/video"),
        filepath.Join(exeDir, "data/video"),
        "/data/video",
        "/app/data/video",
    }, "./data/video")
    p := filepath.Join(vidDir, base)
    if _, err := os.Stat(p); err != nil {
        c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "文件不存在"})
        return
    }
    if err := os.Remove(p); err != nil {
        c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除失败"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": 1, "data": true})
}

func pickDir(candidates []string, fallback string) string {
    for _, d := range candidates {
        if d == "" { continue }
        info, err := os.Stat(d)
        if err == nil && info.IsDir() {
            return d
        }
    }
    return fallback
}
