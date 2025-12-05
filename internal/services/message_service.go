package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
)

// GetAllMessages 封装业务逻辑，获取所有笔记
func GetAllMessages(showPrivate bool) ([]models.Message, error) {
	var messages []models.Message
	var err error

	if showPrivate {
		// 如果是管理员或作者，则不需要过滤私密笔记
		messages, err = repository.GetAllMessages(true)
	} else {
		// 如果不是管理员或作者，则只查询公开的笔记
		messages, err = repository.GetAllMessages(false)
	}

	if err != nil {
		return nil, fmt.Errorf("获取消息失败: %v", err)
	}

	return messages, nil
}

// GetMessageByID 根据 ID 获取笔记
func GetMessageByID(id uint, showPrivate bool) (*models.Message, error) {
	message, err := repository.GetMessageByID(id, showPrivate)
	if err != nil {
		return nil, fmt.Errorf("获取消息失败: %v", err)
	}

	if message == nil {
		return nil, fmt.Errorf("消息不存在")
	}

	return message, nil
}

// GetMessagesByPage 分页获取笔记（管理员查看全部；普通用户可查看公开和自己的私密）
func GetMessagesByPage(page, pageSize int, userID *uint, isAdmin bool) (dto.PageQueryResult, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询数据库
	var messages []models.Message
	var total int64

    if isAdmin {
        // 管理员查看全部，置顶优先
        database.DB.Model(&models.Message{}).Count(&total)
        database.DB.Limit(pageSize).Offset(offset).Order("pinned DESC, created_at DESC").Find(&messages)
    } else if userID != nil {
        // 普通登录用户：公开的 + 自己的私密
        query := database.DB.Model(&models.Message{}).Where("private = ? OR user_id = ?", false, *userID)
        query.Count(&total)
        query.Limit(pageSize).Offset(offset).Order("pinned DESC, created_at DESC").Find(&messages)
    } else {
        // 未登录用户：仅公开
        database.DB.Model(&models.Message{}).Where("private = ?", false).Count(&total)
        database.DB.Limit(pageSize).Offset(offset).Where("private = ?", false).Order("pinned DESC, created_at DESC").Find(&messages)
    }

	// 返回结果
	return dto.PageQueryResult{Total: total, Items: messages}, nil
}

// CreateMessage 发布一条笔记
// 允许所有注册登陆用户发布信息
func CreateMessage(message *models.Message) error {
	user, err := GetUserByID(message.UserID)
	if err != nil {
		return err
	}

	// 删除管理员权限检查，允许所有登录用户发布信息
	message.Username = user.Username // 获取用户名
	return repository.CreateMessage(message)
}

// DeleteMessage 根据 ID 删除笔记
func DeleteMessage(id uint, userID uint) error {
	// 获取笔记信息
	message, err := repository.GetMessageByID(id, true)
	if err != nil {
		return err
	}

	// 验证是否为笔记作者
	if message.UserID != userID {
		return fmt.Errorf("无权删除他人的笔记")
	}

	return repository.DeleteMessage(id)
}

// DeleteMessageByAdmin 管理员删除笔记（无需验证作者）
func DeleteMessageByAdmin(id uint) error {
	return repository.DeleteMessage(id)
}

func GenerateRSS(c *gin.Context) (string, error) {
	messages, err := GetAllMessages(false)
	if err != nil {
		return "", fmt.Errorf("获取消息失败: %v", err)
	}

	// 获取前端配置
	config, err := GetFrontendConfig()
	if err != nil {
		return "", fmt.Errorf("获取配置失败: %v", err)
	}

	// 从配置中安全获取值
	getConfigValue := func(key string, defaultValue string) string {
		if settings, ok := config["frontendSettings"].(map[string]any); ok {
			if value, exists := settings[key].(string); exists && value != "" {
				return value
			}
		}
		return defaultValue
	}

	// 判断请求协议
	schema := "http"
	if c.Request.TLS != nil {
		schema = "https"
	}

	// 处理域名和端口
	requestHost := c.Request.Host
	var baseURL string

	// 从配置获取站点URL，如果没有则使用请求的host
	configURL := getConfigValue("siteURL", "")

	// 检查请求来源是否为反向代理域名
	if strings.Contains(requestHost, "note.noisework.cn") {
		// 如果是从反向代理域名访问的，使用完整的反向代理域名
		baseURL = "https://note.noisework.cn"
	} else if configURL != "" {
		// 使用配置的URL
		baseURL = configURL
		// 确保配置的URL包含正确的端口
		if !strings.Contains(baseURL, ":") && requestHost != "note.noisework.cn" {
			// 从请求中提取端口
			parts := strings.Split(requestHost, ":")
			if len(parts) == 2 && parts[1] == "1314" {
				baseURL = baseURL + ":1314"
			}
		}
	} else {
		// 如果没有配置，确保使用完整的请求地址（包括端口）
		if strings.Contains(requestHost, ":") {
			baseURL = schema + "://" + requestHost
		} else {
			// 默认添加1314端口，如果是直接IP或域名访问
			baseURL = schema + "://" + requestHost + ":1314"
		}
	}

	// 确保URL末尾没有斜杠
	baseURL = strings.TrimSuffix(baseURL, "/")

	feed := &feeds.Feed{
		Title: getConfigValue("rssTitle", "Noise的说说笔记"),
		Link: &feeds.Link{
			Href: baseURL + "/",
		},
		Image: &feeds.Image{
			Url: baseURL + getConfigValue("rssFaviconURL", "/favicon.ico"),
		},
		Description: getConfigValue("rssDescription", "一个说说笔记~"),
		Author: &feeds.Author{
			Name: getConfigValue("rssAuthorName", "Noise"),
		},
		Updated: time.Now(),
	}

	for _, msg := range messages {
		// 处理内容
		content := msg.Content
		if msg.ImageURL != "" {
			imageURL := baseURL + "/api" + msg.ImageURL
			content = fmt.Sprintf("![图片](%s)\n\n%s", imageURL, content)
		}

		// 渲染 Markdown
		htmlContent := pkg.MdToHTML([]byte(content))

		// 生成标题
		title := msg.Username
		if firstLine := pkg.GetFirstLine(msg.Content); firstLine != "" {
			title = firstLine
		}

		// 生成前端页面 URL
		pageURL := baseURL + "/#/messages/" + fmt.Sprintf("%d", msg.ID)

		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: pageURL},
			Description: string(htmlContent),
			Author:      &feeds.Author{Name: msg.Username},
			Created:     msg.CreatedAt,
			Id:          pageURL,
		}

		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", err
	}

	return rss, nil
}

// contains 辅助函数，用于检查字符串中是否包含指定子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// UpdateMessage 更新消息内容
func UpdateMessage(messageID uint, content string) error {
	// 获取消息
	message, err := repository.GetMessageByID(messageID, true)
	if err != nil {
		return fmt.Errorf("获取消息失败: %v", err)
	}
	if message == nil {
		return fmt.Errorf("消息不存在")
	}

	// 更新消息内容
	message.Content = content
	if err := database.DB.Save(message).Error; err != nil {
		return fmt.Errorf("更新消息失败: %v", err)
	}

    return nil
}

// UpdateMessagePinned 更新消息置顶状态
func UpdateMessagePinned(messageID uint, pinned bool) error {
    message, err := repository.GetMessageByID(messageID, true)
    if err != nil {
        return fmt.Errorf("获取消息失败: %v", err)
    }
    if message == nil {
        return fmt.Errorf("消息不存在")
    }
    message.Pinned = pinned
    if err := database.DB.Save(message).Error; err != nil {
        return fmt.Errorf("更新置顶状态失败: %v", err)
    }
    return nil
}

// IncrementLikeCount 点赞计数加一
func IncrementLikeCount(messageID uint) (int, error) {
    message, err := repository.GetMessageByID(messageID, false)
    if err != nil {
        return 0, fmt.Errorf("获取消息失败: %v", err)
    }
    if message == nil {
        return 0, fmt.Errorf("消息不存在")
    }
    message.LikeCount = message.LikeCount + 1
    if err := database.DB.Save(message).Error; err != nil {
        return 0, fmt.Errorf("更新点赞失败: %v", err)
    }
    return message.LikeCount, nil
}

// ToggleLike 根据 session 或用户切换点赞状态
func ToggleLike(messageID uint, userID *uint, sessionID string) (bool, int, error) {
    if sessionID == "" && (userID == nil || *userID == 0) {
        return false, 0, fmt.Errorf("缺少会话或用户信息")
    }
    // 查询是否已有点赞
    var existing models.MessageLike
    q := database.DB.Where("message_id = ?", messageID)
    if userID != nil && *userID != 0 {
        q = q.Where("user_id = ?", *userID)
    } else {
        q = q.Where("session_id = ?", sessionID)
    }
    if err := q.First(&existing).Error; err == nil && existing.ID != 0 {
        // 已点赞 -> 取消
        if err := database.DB.Delete(&existing).Error; err != nil {
            return false, 0, err
        }
    } else {
        // 未点赞 -> 新增
        like := models.MessageLike{ MessageID: messageID, SessionID: sessionID }
        if userID != nil && *userID != 0 { like.UserID = userID }
        if err := database.DB.Create(&like).Error; err != nil {
            return false, 0, err
        }
    }
    // 重新统计总数并同步
    var cnt int64
    if err := database.DB.Model(&models.MessageLike{}).Where("message_id = ?", messageID).Count(&cnt).Error; err != nil {
        return false, 0, err
    }
    if err := database.DB.Model(&models.Message{}).Where("id = ?", messageID).Update("like_count", cnt).Error; err != nil {
        return false, int(cnt), err
    }
    // 再查一次是否点赞状态
    var check models.MessageLike
    q2 := database.DB.Where("message_id = ?", messageID)
    if userID != nil && *userID != 0 { q2 = q2.Where("user_id = ?", *userID) } else { q2 = q2.Where("session_id = ?", sessionID) }
    liked := q2.First(&check).Error == nil && check.ID != 0
    return liked, int(cnt), nil
}
func GetMessagesGroupByDate() ([]struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}, error) {
	var results []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	// 移除 deleted_at 条件，因为该列不存在
	err := database.DB.Table("messages").
		Select("DATE(created_at) as date, COUNT(*) as count").
		// 移除这一行: Where("deleted_at IS NULL").
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&results).Error

	if err != nil {
		fmt.Printf("获取消息日历数据失败: %v\n", err)
		return nil, err
	}

	// 如果结果为空，返回空数组而不是nil
	if len(results) == 0 {
		return []struct {
			Date  string `json:"date"`
			Count int    `json:"count"`
		}{}, nil
	}

	return results, nil
}

// GetMessagePage 获取消息详情页
func GetMessagePage(id uint) (*models.Message, error) {
	message, err := repository.GetMessageByID(id, false)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, fmt.Errorf("消息不存在")
	}

	// 如果是私密消息，需要进行额外处理
	if message.Private {
		return nil, fmt.Errorf("无权访问")
	}

	return message, nil
}

func SearchMessages(keyword string, page, pageSize int, showPrivate bool) (dto.PageQueryResult, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 直接使用服务层实现
	query := database.DB.Model(&models.Message{}).
		Select("id, content, created_at, username, image_url, private, user_id").
		Where("content LIKE ?", "%"+keyword+"%")

	if !showPrivate {
		query = query.Where("private = ?", false)
	}

	var total int64
	var messages []models.Message

	err := query.Count(&total).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Order("created_at DESC").
		Find(&messages).Error

	if err != nil {
		return dto.PageQueryResult{}, err
	}

	// 确保返回的数据结构符合前端期望
	return dto.PageQueryResult{
		Total: total,
		Items: messages,
	}, nil
}
