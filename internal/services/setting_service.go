package services

import (
	"encoding/json"
	"fmt"

	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
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

	configMap := map[string]interface{}{
		"allowRegistration": allowReg,
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
			// PWA 设置
			"pwaEnabled":     config.PwaEnabled,
			"pwaTitle":       choose(config.PwaTitle, config.SiteTitle),
			"pwaDescription": choose(config.PwaDescription, config.Description),
			"pwaIconURL":     choose(config.PwaIconURL, config.RSSFaviconURL),
		},
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
	if v, ok := frontendSettings["enableGithubCard"].(bool); ok {
		config.EnableGithubCard = v
	} else if vs, ok := frontendSettings["enableGithubCard"].(string); ok {
		if vs == "true" {
			config.EnableGithubCard = true
		} else if vs == "false" {
			config.EnableGithubCard = false
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
			"rssFaviconURL":    "/favicon.ico",
			"walineServerURL":  "请前往waline官网https://waline.js.org查看部署配置",
			"enableGithubCard": false,
			// PWA 设置默认值
			"pwaEnabled":     true,
			"pwaTitle":       "",
			"pwaDescription": "",
			"pwaIconURL":     "",
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
