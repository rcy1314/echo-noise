package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings" // 添加 strings 包
	"sync"    // 添加 sync 包
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
)

func checkUser(c *gin.Context) (*models.User, error) {
	userID, exists := c.Get("user_id") // 修改 userid 为 user_id
	if !exists {
		return nil, fmt.Errorf(models.UserNotFoundMessage)
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		return nil, fmt.Errorf(models.UserNotFoundMessage)
	}
	return user, nil
}

func Login(c *gin.Context) {
	var loginDto dto.LoginDto
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[any]("参数错误"))
		return
	}

	user, err := services.Login(loginDto)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[any](err.Error()))
		return
	}

	session := sessions.Default(c)
	session.Clear()
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	session.Set("is_admin", user.IsAdmin)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusOK, dto.Fail[any]("Session 保存失败"))
		return
	}

	c.JSON(http.StatusOK, dto.OK(user, "登录成功"))
}

// 添加登出功能
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, dto.OK[any](nil, "登出成功"))
}
func Register(c *gin.Context) {
	// 新增：注册前判断是否允许注册
	db, _ := database.GetDB()
	var setting models.Setting
	allowReg := true
	if err := db.Table("settings").First(&setting).Error; err == nil {
		allowReg = setting.AllowRegistration
	}

	if !allowReg {
		c.JSON(http.StatusOK, dto.Fail[string]("当前不允许注册新用户"))
		return
	}

	var user dto.RegisterDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	if err := services.Register(user); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.RegisterSuccessMessage))
}

// GetMessages 处理 GET /messages 请求，返回所有留言
func GetMessages(c *gin.Context) {
	showPrivate := false
	userID, exists := c.Get("user_id")
	if exists {
		user, err := services.GetUserByID(userID.(uint))
		if err == nil && user.IsAdmin {
			showPrivate = true
		}
	}

	messages, err := services.GetAllMessages(showPrivate)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetAllMessagesFailMessage))
		return
	}
	c.JSON(http.StatusOK, dto.OK(messages, models.GetAllMessagesSuccess))
}

// GetMessage 处理 GET /messages/:id 请求，获取留言详情
func GetMessage(c *gin.Context) {
	// 从 URL 参数获取留言 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
		return
	}

	// 检查是否显示私密消息
	showPrivate := false
	userID, exists := c.Get("user_id")
	if exists {
		user, err := services.GetUserByID(userID.(uint))
		if err == nil && user.IsAdmin {
			showPrivate = true
		}
	}

	// 调用 Service 层根据 ID 获取留言
	message, err := services.GetMessageByID(uint(id), showPrivate)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetMessageByIDFailMessage))
		return
	}

	if message == nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.MessageNotFoundMessage))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.OK(message, models.GetMessageByIDSuccess))
}

func GetMessagesByPage(c *gin.Context) {
	var page, pageSize int = 1, 10

	// 尝试从 POST JSON 数据获取分页参数
	var pageRequest dto.PageQueryDto
	if err := c.ShouldBindJSON(&pageRequest); err == nil {
		page = pageRequest.Page
		pageSize = pageRequest.PageSize
	} else {
		// 如果不是 POST JSON，则尝试从 URL 查询参数获取
		if pageStr := c.Query("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		if sizeStr := c.Query("pageSize"); sizeStr != "" {
			if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
				pageSize = s
			}
		}
	}

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 检查权限并传递用户上下文
	var currentUserID *uint
	isAdmin := false
	// 优先从上下文获取（由中间件设置）
	if uid, exists := c.Get("user_id"); exists {
		u, err := services.GetUserByID(uid.(uint))
		if err == nil {
			id := u.ID
			currentUserID = &id
			isAdmin = u.IsAdmin
		}
	} else {
		// 兼容未使用鉴权中间件的场景：从 session 获取
		session := sessions.Default(c)
		if v := session.Get("user_id"); v != nil {
			switch val := v.(type) {
			case uint:
				id := val
				currentUserID = &id
			case int:
				id := uint(val)
				currentUserID = &id
			case int64:
				id := uint(val)
				currentUserID = &id
			case float64:
				id := uint(val)
				currentUserID = &id
			case string:
				if parsed, err := strconv.ParseUint(val, 10, 64); err == nil {
					id := uint(parsed)
					currentUserID = &id
				}
			}
		}
		if v := session.Get("is_admin"); v != nil {
			switch val := v.(type) {
			case bool:
				isAdmin = val
			case int:
				isAdmin = val != 0
			case int64:
				isAdmin = val != 0
			case float64:
				isAdmin = val != 0
			case string:
				isAdmin = val == "true" || val == "1"
			}
		}
		// 如果仅拿到 user_id，则再查一次用户，确保 isAdmin
		if currentUserID != nil && !isAdmin {
			u, err := services.GetUserByID(*currentUserID)
			if err == nil {
				isAdmin = u.IsAdmin
			}
		}
	}

	pageQueryResult, err := services.GetMessagesByPage(page, pageSize, currentUserID, isAdmin)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(pageQueryResult, models.GetMessagesByPageSuccess))
}
func GetStatus(c *gin.Context) {
	status, err := services.GetStatus()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetStatusFailMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(status, models.GetStatusSuccessMessage))
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	messageID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "无效的消息ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK, dto.Fail[string]("未授权访问"))
		return
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	if !user.IsAdmin {
		if err := services.DeleteMessage(uint(messageID), userID.(uint)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
			return
		}
	} else {
		if err := services.DeleteMessageByAdmin(uint(messageID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "删除成功"})
}

func GenerateRSS(c *gin.Context) {
	atom, err := services.GenerateRSS(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail[string](models.GenerateRSSFailMessage))
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(atom))
}

func UpdateUser(c *gin.Context) {
	user, err := checkUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	var userdto dto.UserInfoDto
	if err := c.ShouldBindJSON(&userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	if err := services.UpdateUser(user, userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateUserSuccessMessage))
}

func ChangePassword(c *gin.Context) {
	var userdto dto.UserInfoDto
	if err := c.ShouldBindJSON(&userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	user, err := checkUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	if err := services.ChangePassword(user, userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.ChangePasswordSuccessMessage))
}

// checkAdmin 函数需要重新添加
func checkAdmin(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("未授权访问")
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		return 0, err
	}

	if !user.IsAdmin {
		return 0, fmt.Errorf("需要管理员权限")
	}

	return userID.(uint), nil
}

func UpdateUserAdmin(c *gin.Context) {
	_, err := checkAdmin(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 1 {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
		return
	}

	if err := services.UpdateUserAdmin(uint(id)); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateUserSuccessMessage))
}

func GetUserInfo(c *gin.Context) {
	user, err := checkUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, dto.OK(user, models.QuerySuccessMessage))
}

func UpdateSetting(c *gin.Context) {
	_, err := checkAdmin(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	db, _ := database.GetDB()
	var oldSetting models.Setting
	if err := db.Table("settings").First(&oldSetting).Error; err != nil {
		c.JSON(http.StatusOK, dto.Fail[string]("读取原有配置失败"))
		return
	}

	// 解析前端传来的配置
	var setting dto.SettingDto
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 只更新传递的字段，未传递的字段保持原值
	oldSetting.AllowRegistration = setting.AllowRegistration

	// 合并前端配置
	var frontendSettings map[string]interface{}
	if setting.FrontendSettings.SiteTitle == "" &&
		setting.FrontendSettings.SubtitleText == "" &&
		setting.FrontendSettings.AvatarURL == "" &&
		setting.FrontendSettings.Username == "" &&
		setting.FrontendSettings.Description == "" &&
		len(setting.FrontendSettings.Backgrounds) == 0 &&
		setting.FrontendSettings.CardFooterTitle == "" &&
		setting.FrontendSettings.CardFooterLink == "" &&
		setting.FrontendSettings.PageFooterHTML == "" &&
		setting.FrontendSettings.RSSTitle == "" &&
		setting.FrontendSettings.RSSDescription == "" &&
		setting.FrontendSettings.RSSAuthorName == "" &&
		setting.FrontendSettings.RSSFaviconURL == "" &&
		setting.FrontendSettings.WalineServerURL == "" &&
		setting.FrontendSettings.EnableGithubCard == nil {
		// 前端未传递 frontendSettings，自动从数据库读取
		var config models.SiteConfig
		if err := db.Table("site_configs").First(&config).Error; err == nil {
			frontendSettings = map[string]interface{}{
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
				// PWA 设置（直接从数据库配置）
				"pwaEnabled":     config.PwaEnabled,
				"pwaTitle":       config.PwaTitle,
				"pwaDescription": config.PwaDescription,
				"pwaIconURL":     config.PwaIconURL,
			}
		}
	} else {
		frontendSettings = map[string]interface{}{}
		if setting.FrontendSettings.SiteTitle != "" {
			frontendSettings["siteTitle"] = setting.FrontendSettings.SiteTitle
		}
		if setting.FrontendSettings.SubtitleText != "" {
			frontendSettings["subtitleText"] = setting.FrontendSettings.SubtitleText
		}
		if setting.FrontendSettings.AvatarURL != "" {
			frontendSettings["avatarURL"] = setting.FrontendSettings.AvatarURL
		}
		if setting.FrontendSettings.Username != "" {
			frontendSettings["username"] = setting.FrontendSettings.Username
		}
		if setting.FrontendSettings.Description != "" {
			frontendSettings["description"] = setting.FrontendSettings.Description
		}
		if len(setting.FrontendSettings.Backgrounds) > 0 {
			frontendSettings["backgrounds"] = setting.FrontendSettings.Backgrounds
		}
		if setting.FrontendSettings.CardFooterTitle != "" {
			frontendSettings["cardFooterTitle"] = setting.FrontendSettings.CardFooterTitle
		}
		if setting.FrontendSettings.CardFooterLink != "" {
			frontendSettings["cardFooterLink"] = setting.FrontendSettings.CardFooterLink
		}
		if setting.FrontendSettings.PageFooterHTML != "" {
			frontendSettings["pageFooterHTML"] = setting.FrontendSettings.PageFooterHTML
		}
		if setting.FrontendSettings.RSSTitle != "" {
			frontendSettings["rssTitle"] = setting.FrontendSettings.RSSTitle
		}
		if setting.FrontendSettings.RSSDescription != "" {
			frontendSettings["rssDescription"] = setting.FrontendSettings.RSSDescription
		}
		if setting.FrontendSettings.RSSAuthorName != "" {
			frontendSettings["rssAuthorName"] = setting.FrontendSettings.RSSAuthorName
		}
		if setting.FrontendSettings.RSSFaviconURL != "" {
			frontendSettings["rssFaviconURL"] = setting.FrontendSettings.RSSFaviconURL
		}
		if setting.FrontendSettings.WalineServerURL != "" {
			frontendSettings["walineServerURL"] = setting.FrontendSettings.WalineServerURL
		}
		if setting.FrontendSettings.EnableGithubCard != nil {
			frontendSettings["enableGithubCard"] = *setting.FrontendSettings.EnableGithubCard
		}
		if setting.FrontendSettings.PwaEnabled != nil {
			frontendSettings["pwaEnabled"] = *setting.FrontendSettings.PwaEnabled
		}
		if setting.FrontendSettings.PwaTitle != nil {
			frontendSettings["pwaTitle"] = *setting.FrontendSettings.PwaTitle
		}
		if setting.FrontendSettings.PwaDescription != nil {
			frontendSettings["pwaDescription"] = *setting.FrontendSettings.PwaDescription
		}
		if setting.FrontendSettings.PwaIconURL != nil {
			frontendSettings["pwaIconURL"] = *setting.FrontendSettings.PwaIconURL
		}
	}

	settingMap := map[string]interface{}{
		"allowRegistration": setting.AllowRegistration,
		"frontendSettings":  frontendSettings,
	}
	if err := services.UpdateFrontendSetting(0, settingMap); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string]("保存前端配置失败: "+err.Error()))
		return
	}

	if err := db.Table("settings").Save(&oldSetting).Error; err != nil {
		c.JSON(http.StatusOK, dto.Fail[string]("保存配置失败"))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateSettingSuccessMessage))
}

func GetFrontendConfig(c *gin.Context) {
	config, err := services.GetFrontendConfig()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "获取配置失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "data": config})
}

// 动态生成 Web Manifest
func GetWebManifest(c *gin.Context) {
	configMap, _ := services.GetFrontendConfig()
	fs := map[string]interface{}{}
	if v, ok := configMap["frontendSettings"].(map[string]interface{}); ok {
		fs = v
	}

	// 读取 PWA 设置（优先用 PWA 字段，否则回退到站点字段）
	pwaEnabled := true
	if v, ok := fs["pwaEnabled"].(bool); ok {
		pwaEnabled = v
	}
	title := "说说笔记"
	description := ""
	icon := "/favicon.ico"

	if pwaEnabled {
		if v, ok := fs["pwaTitle"].(string); ok && v != "" {
			title = v
		}
		if v, ok := fs["pwaDescription"].(string); ok {
			description = v
		}
		if v, ok := fs["pwaIconURL"].(string); ok && v != "" {
			icon = v
		}
	}
	if title == "说说笔记" {
		if v, ok := fs["siteTitle"].(string); ok && v != "" {
			title = v
		}
	}
	if description == "" {
		if v, ok := fs["description"].(string); ok {
			description = v
		}
	}
	if icon == "/favicon.ico" {
		if v, ok := fs["rssFaviconURL"].(string); ok && v != "" {
			icon = v
		}
	}

	manifest := map[string]interface{}{
		"name":             title,
		"short_name":       title,
		"description":      description,
		"start_url":        "/",
		"display":          "standalone",
		"background_color": "#000000",
		"theme_color":      "#000000",
		"icons": []map[string]string{
			{"src": icon, "sizes": "any", "type": "image/x-icon"},
			{"src": "/android-chrome-192x192.png", "sizes": "192x192", "type": "image/png"},
			{"src": "/android-chrome-512x512.png", "sizes": "512x512", "type": "image/png"},
			{"src": "/apple-touch-icon.png", "sizes": "180x180", "type": "image/png"},
		},
	}

	c.Header("Content-Type", "application/manifest+json; charset=utf-8")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.JSON(http.StatusOK, manifest)
}

func UpdateMessage(c *gin.Context) {
	// 获取消息ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "消息ID不能为空"})
		return
	}

	// 检查用户权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 0, "msg": "未授权访问"})
		return
	}

	// 获取请求体
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "请求参数错误"})
		return
	}

	// 检查消息是否存在并且属于当前用户
	messageID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "无效的消息ID"})
		return
	}

	// 获取用户信息
	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "获取用户信息失败"})
		return
	}

	// 检查消息所有权或管理员权限
	message, err := services.GetMessageByID(uint(messageID), true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 0, "msg": "消息不存在"})
		return
	}

	if !user.IsAdmin && message.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "无权限修改此消息"})
		return
	}

	// 更新消息
	if err := services.UpdateMessage(uint(messageID), req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "更新成功"})
}
func GetMessagesCalendar(c *gin.Context) {
	// 改为调用 services 层方法
	calendarData, err := services.GetMessagesGroupByDate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": calendarData,
	})
}
func SearchMessages(c *gin.Context) {
	// 从查询参数获取数据
	keyword := c.Query("keyword")
	page := 1
	pageSize := 10

	// 尝试解析页码和每页数量
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr := c.Query("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	// 默认只搜索公开内容
	showPrivate := false

	// 检查用户是否为管理员，管理员可以看到私密内容
	userID, exists := c.Get("user_id")
	if exists {
		user, err := services.GetUserByID(userID.(uint))
		if err == nil && user.IsAdmin {
			showPrivate = true
		}
	}

	result, err := services.SearchMessages(keyword, page, pageSize, showPrivate)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	// 直接构造符合前端期望的JSON格式
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "搜索成功",
		"data": result,
	})
}
func GetUserToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK, dto.Fail[any]("未授权访问"))
		return
	}

	token, err := services.GetUserToken(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[any]("获取Token失败"))
		return
	}

	c.JSON(http.StatusOK, dto.OK(gin.H{
		"token": token,
	}, "获取成功"))
}

func RegenerateUserToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK, dto.Fail[any]("未授权访问"))
		return
	}

	token, err := services.RegenerateUserToken(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[any]("更新Token失败"))
		return
	}

	c.JSON(http.StatusOK, dto.OK(gin.H{
		"token": token,
	}, "更新成功"))
}

// RefreshRSS 刷新 RSS 内容
func RefreshRSS(c *gin.Context) {
	// 重新生成 RSS
	_, err := services.GenerateRSS(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  "RSS 刷新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "RSS 已刷新",
	})
}

// 检查版本更新
func CheckVersion(c *gin.Context) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://hub.docker.com/v2/repositories/noise233/echo-noise/tags")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "检查更新失败",
		})
		return
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Name        string `json:"name"`
			LastUpdated string `json:"last_updated"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "解析版本信息失败",
		})
		return
	}

	var lastUpdateTime string
	for _, tag := range result.Results {
		if tag.Name == "latest" {
			lastUpdateTime = tag.LastUpdated
			break
		}
	}

	if lastUpdateTime == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "未找到版本信息",
		})
		return
	}

	updateTime, err := time.Parse(time.RFC3339, lastUpdateTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "解析时间失败",
		})
		return
	}

	// 判断是否在24小时内更新
	hasUpdate := time.Since(updateTime) <= 24*time.Hour

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": gin.H{
			"hasUpdate":      hasUpdate,
			"lastUpdateTime": lastUpdateTime,
		},
	})
}

// GetNotifyConfig 获取推送配置
func GetNotifyConfig(c *gin.Context) {
	_, err := checkAdmin(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	config := models.GetNotifyConfig()
	if config == nil {
		// 如果配置不存在，返回空配置（所有字段默认值）
		config = &models.NotifyConfig{
			WebhookEnabled:           false,
			WebhookURL:               "",
			TelegramEnabled:          false,
			TelegramToken:            "",
			TelegramChatID:           "",
			WeworkEnabled:            false,
			WeworkKey:                "",
			FeishuEnabled:            false,
			FeishuWebhook:            "",
			FeishuSecret:             "",
			TwitterEnabled:           false,
			TwitterApiKey:            "",
			TwitterApiSecret:         "",
			TwitterAccessToken:       "",
			TwitterAccessTokenSecret: "",
			CustomHttpEnabled:        false,
			CustomHttpUrl:            "",
			CustomHttpMethod:         "",
			CustomHttpHeaders:        "",
			CustomHttpBody:           "",
		}
	}
	c.JSON(http.StatusOK, dto.OK(config, "获取成功"))
}

// SaveNotifyConfig 保存推送配置
func SaveNotifyConfig(c *gin.Context) {
	_, err := checkAdmin(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	var config models.NotifyConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string]("无效的配置数据"))
		return
	}
	// Twitter校验
	if config.TwitterEnabled {
		if config.TwitterApiKey == "" || config.TwitterApiSecret == "" || config.TwitterAccessToken == "" || config.TwitterAccessTokenSecret == "" {
			c.JSON(http.StatusOK, dto.Fail[string]("Twitter配置不完整"))
			return
		}
	}
	// 自定义HTTP校验
	if config.CustomHttpEnabled {
		if config.CustomHttpUrl == "" {
			c.JSON(http.StatusOK, dto.Fail[string]("自定义HTTP URL不能为空"))
			return
		}
	}
	// 根据启用状态验证配置
	if config.WebhookEnabled {
		if config.WebhookURL == "" {
			c.JSON(http.StatusOK, dto.Fail[string]("Webhook URL 不能为空"))
			return
		}
	}
	if config.TelegramEnabled {
		if config.TelegramToken == "" || config.TelegramChatID == "" {
			c.JSON(http.StatusOK, dto.Fail[string]("Telegram 配置不完整"))
			return
		}
	}
	if config.WeworkEnabled {
		if config.WeworkKey == "" {
			c.JSON(http.StatusOK, dto.Fail[string]("企业微信 Key 不能为空"))
			return
		}
	}
	if config.FeishuEnabled {
		if config.FeishuWebhook == "" {
			c.JSON(http.StatusOK, dto.Fail[string]("飞书 Webhook 不能为空"))
			return
		}
	}

	// 打印调试信息（临时添加）
	fmt.Printf("保存前的Twitter配置: Enabled=%v, Key=%s, Secret=%s\n",
		config.TwitterEnabled,
		config.TwitterApiKey,
		config.TwitterApiSecret)
	fmt.Printf("保存前的CustomHttp配置: Enabled=%v, Url=%s\n",
		config.CustomHttpEnabled,
		config.CustomHttpUrl)

	if err := models.SaveNotifyConfig(config); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string]("保存配置失败: "+err.Error()))
		return
	}

	// 获取保存后的配置（调试用）
	savedConfig := models.GetNotifyConfig()
	fmt.Printf("保存后的Twitter配置: Enabled=%v\n", savedConfig.TwitterEnabled)
	fmt.Printf("保存后的CustomHttp配置: Enabled=%v\n", savedConfig.CustomHttpEnabled)

	c.JSON(http.StatusOK, dto.OK[any](nil, "配置已更新"))
}

// TestNotify 测试推送
func TestNotify(c *gin.Context) {
	_, err := checkAdmin(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	var request struct {
		Type string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string]("无效的请求参数"))
		return
	}

	testMsg := "这是一条测试消息 - " + time.Now().Format("2006-01-02 15:04:05")
	var emptyImages []string

	var testErr error
	switch request.Type {
	case "webhook":
		testErr = models.SendWebhook(testMsg)
	case "telegram":
		testErr = models.SendTelegram(testMsg, emptyImages)
	case "wework":
		testErr = models.SendWework(testMsg, emptyImages)
	case "feishu":
		testErr = models.SendFeishu(testMsg)
	case "twitter":
		testErr = models.SendTwitter(testMsg)
	case "customHttp":
		testErr = models.SendCustomHttp(testMsg)
	default:
		c.JSON(http.StatusOK, dto.Fail[string]("不支持的推送类型"))
		return
	}

	if testErr != nil {
		c.JSON(http.StatusOK, dto.Fail[string](fmt.Sprintf("推送测试失败: %v", testErr)))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, "推送测试已发送"))
}

// 保留这个新版本的 PostMessage 函数
func PostMessage(c *gin.Context) {
	// 解析请求数据
	var request struct {
		Content  string `json:"content"`
		Private  bool   `json:"private"`
		ImageURL string `json:"image_url"`
		VideoURL string `json:"video_url"` // 新增视频字段
		Notify   bool   `json:"notify"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string]("内容不能为空"))
		return
	}

	// 验证用户身份
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK, dto.Fail[string]("未授权访问"))
		return
	}

	// 创建消息
	message := &models.Message{
		Content:  request.Content,
		Private:  request.Private,
		ImageURL: request.ImageURL,
		UserID:   userID.(uint),
	}

	if err := services.CreateMessage(message); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	// 处理推送逻辑
	if request.Notify {
		notifyConfig := models.GetNotifyConfig()
		if notifyConfig != nil {
			// 提取内容中的第一张图片链接
			var firstImageURL string
			var firstVideoURL string
			var formattedContent string

			// 如果已有上传的图片，优先使用
			if request.ImageURL != "" {
				firstImageURL = request.ImageURL
			}
			// 如果已有上传的视频，优先使用
			if request.VideoURL != "" {
				firstVideoURL = request.VideoURL
			}

			// 从 Markdown 内容中提取第一张图片
			imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
			matches := imageRegex.FindAllStringSubmatch(request.Content, -1)
			if firstImageURL == "" && len(matches) > 0 {
				firstImageURL = matches[0][2]
			}

			// 从 Markdown 内容中提取第一段视频（如 [video](url)）
			videoRegex := regexp.MustCompile(`\[video\]\(([^)]+)\)`)
			videoMatches := videoRegex.FindAllStringSubmatch(request.Content, -1)
			if firstVideoURL == "" && len(videoMatches) > 0 {
				firstVideoURL = videoMatches[0][1]
			}

			// 移除第一张图片的Markdown语法
			formattedContent = request.Content
			if len(matches) > 0 {
				formattedContent = imageRegex.ReplaceAllStringFunc(request.Content, func(match string) string {
					if match == matches[0][0] {
						return ""
					}
					return match
				})
			}

			// 处理长内容，如果超过4000字符，进行截断
			const maxContentLength = 4000
			var truncatedContent string
			if len(formattedContent) > maxContentLength {
				truncatedContent = formattedContent[:maxContentLength] + "...\n(内容过长，已截断)"
			} else {
				truncatedContent = formattedContent
			}

			// 格式化内容，处理Markdown语法
			headingRegex := regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
			truncatedContent = headingRegex.ReplaceAllString(truncatedContent, "$1 $2")

			// 准备图片和视频数组
			var images []string
			var videos []string
			if firstImageURL != "" {
				images = []string{firstImageURL}
			}
			if firstVideoURL != "" {
				videos = []string{firstVideoURL}
			}

			go func() {
				// Webhook
				if notifyConfig.WebhookEnabled && notifyConfig.WebhookURL != "" {
					models.SendWebhook(truncatedContent)
				}

				// Telegram
				if notifyConfig.TelegramEnabled && notifyConfig.TelegramToken != "" && notifyConfig.TelegramChatID != "" {
					const telegramMaxText = 4096
					const telegramMaxCaption = 1024

					isPublicURL := func(url string) bool {
						return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
					}

					// 推送图片
					if len(images) > 0 {
						if isPublicURL(images[0]) {
							caption := formattedContent
							if len(caption) > telegramMaxCaption {
								caption = caption[:telegramMaxCaption] + "...\n(内容过长，已截断)"
							}
							err := models.SendTelegramPhotoWithCaption(images[0], caption)
							if err != nil {
								sendTelegramErrorNotify(c, err)
							}
						} else {
							msg := formattedContent + "\n[图片] " + images[0]
							if len(msg) > telegramMaxText {
								msg = msg[:telegramMaxText] + "...\n(内容过长，已截断)"
							}
							err := models.SendTelegramMessage(msg)
							if err != nil {
								sendTelegramErrorNotify(c, err)
							}
						}
					}

					// 推送视频
					if len(videos) > 0 {
						if isPublicURL(videos[0]) {
							caption := formattedContent
							if len(caption) > telegramMaxCaption {
								caption = caption[:telegramMaxCaption] + "...\n(内容过长，已截断)"
							}
							err := models.SendTelegramVideoWithCaption(videos[0], caption)
							if err != nil {
								sendTelegramErrorNotify(c, err)
							}
						} else {
							msg := formattedContent + "\n[视频] " + videos[0]
							if len(msg) > telegramMaxText {
								msg = msg[:telegramMaxText] + "...\n(内容过长，已截断)"
							}
							err := models.SendTelegramMessage(msg)
							if err != nil {
								sendTelegramErrorNotify(c, err)
							}
						}
					}

					// 没有图片和视频，直接发文本
					if len(images) == 0 && len(videos) == 0 {
						if len(formattedContent) > telegramMaxText {
							sendTelegramErrorNotify(c, fmt.Errorf("Telegram 文本内容超出最大长度（%d 字符）", telegramMaxText))
						} else {
							err := models.SendTelegramMessage(formattedContent)
							if err != nil {
								sendTelegramErrorNotify(c, err)
							}
						}
					}
				}

				// 企业微信
				if notifyConfig.WeworkEnabled && notifyConfig.WeworkKey != "" {
					const weworkMaxLength = 2000
					var weworkContent string
					if len(formattedContent) > weworkMaxLength {
						weworkContent = formattedContent[:weworkMaxLength] + "...\n(内容过长，已截断)"
					} else {
						weworkContent = formattedContent
					}
					models.SendWework(weworkContent, images)
				}

				// 飞书
				if notifyConfig.FeishuEnabled && notifyConfig.FeishuWebhook != "" {
					const feishuMaxLength = 2000
					var feishuContent string
					if len(formattedContent) > feishuMaxLength {
						feishuContent = formattedContent[:feishuMaxLength] + "...\n(内容过长，已截断)"
					} else {
						feishuContent = formattedContent
					}
					models.SendFeishu(feishuContent)
				}
			}()
		}
	}

	c.JSON(http.StatusOK, dto.OK(message, "发布成功"))
}

// 上传视频
func UploadVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "未选择视频文件"})
		return
	}

	// 检查文件类型和大小
	allowedExts := map[string]bool{".mp4": true, ".webm": true, ".mov": true, ".avi": true}
	ext := strings.ToLower(filepath.Ext(file.Filename)) // 兼容大小写
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "仅支持 mp4/webm/mov/avi 格式"})
		return
	}
	if file.Size > 200*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "视频不能超过200MB"})
		return
	}

	// 保存到 data/video 目录
	saveDir := "./data/video"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "无法创建目录"})
		return
	}
	newName := uuid.New().String() + ext
	savePath := filepath.Join(saveDir, newName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "保存失败"})
		return
	}

	// 返回视频访问路径
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "上传成功",
		"data": "/video/" + newName,
	})
}
func SendNotify(c *gin.Context) {
	var request struct {
		Content string   `json:"content"`
		Images  []string `json:"images"`
		Format  string   `json:"format"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "请求参数错误"})
		return
	}

	// 验证内容不为空
	if request.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "推送内容不能为空"})
		return
	}

	// 获取推送配置
	config := models.GetNotifyConfig()
	if config == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "推送配置不存在"})
		return
	}

	// 并发处理所有启用的推送渠道
	var wg sync.WaitGroup
	errorChan := make(chan error, 6) // 增加通道容量

	// Telegram
	if config.TelegramEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := models.SendTelegram(request.Content, request.Images); err != nil {
				errorChan <- fmt.Errorf("Telegram: %v", err)
			}
		}()
	}

	// 企业微信
	if config.WeworkEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := models.SendWework(request.Content, request.Images); err != nil {
				errorChan <- fmt.Errorf("企业微信: %v", err)
			}
		}()
	}

	// 飞书
	if config.FeishuEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := models.SendFeishu(request.Content); err != nil {
				errorChan <- fmt.Errorf("飞书: %v", err)
			}
		}()
	}

	// Webhook
	if config.WebhookEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := models.SendWebhook(request.Content); err != nil {
				errorChan <- fmt.Errorf("Webhook: %v", err)
			}
		}()
	}
	// Twitter
	if config.TwitterEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Twitter 字数限制 280
			tweet := request.Content
			if len([]rune(tweet)) > 280 {
				tweet = string([]rune(tweet)[:280]) + "...(内容截断)"
			}
			if err := models.SendTwitter(tweet); err != nil {
				errorChan <- fmt.Errorf("Twitter: %v", err)
			}
		}()
	}

	// 自定义 HTTP
	if config.CustomHttpEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := models.SendCustomHttp(request.Content); err != nil {
				errorChan <- fmt.Errorf("CustomHttp: %v", err)
			}
		}()
	}

	// 等待所有推送完成
	wg.Wait()
	close(errorChan)

	// 收集错误
	var errors []string
	for err := range errorChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("部分推送失败: %s", strings.Join(errors, "; ")),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "推送成功",
	})
}
func sendTelegramErrorNotify(c *gin.Context, err error) {
	log.Printf("Telegram 推送失败: %v", err)
}
