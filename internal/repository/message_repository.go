package repository

import (
	"errors"
	"strings"

	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
	"gorm.io/gorm"
)

// GetAllMessages 从数据库获取所有留言
func GetAllMessages(showPrivate bool) ([]models.Message, error) {
	var messages []models.Message

    // 是否将私密内容也查询出来（置顶优先）
    if showPrivate {
        if err := database.DB.Order("pinned DESC, created_at DESC").Find(&messages).Error; err != nil {
            return nil, err
        }
    } else {
        if err := database.DB.Where("private = ?", false).Order("pinned DESC, created_at DESC").Find(&messages).Error; err != nil {
            return nil, err
        }
    }

	return messages, nil
}

// GetMessageByID 根据 ID 获取留言
func GetMessageByID(id uint, showPrivate bool) (*models.Message, error) {
    var message models.Message
    result := database.DB.First(&message, id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, errors.New("消息不存在")
        }
        return nil, result.Error
    }

    if !showPrivate && message.Private {
        return nil, errors.New("无权访问该消息")
    }

    return &message, nil
}

// CreateMessage 保存一条留言
func CreateMessage(message *models.Message) error {
	// 防止 XSS 攻击，使用 bluemonday 清理器来清理 HTML 标签
	// p := bluemonday.UGCPolicy()                   // 创建一个新的 HTML 清理器
	// message.Content = p.Sanitize(message.Content) // 清理内容中的 HTML 标签

	message.Content = strings.TrimSpace(message.Content) // 去除内容前后的空格
	// message.Username = strings.TrimSpace(message.Username) // 去除用户名前后的空格

	if message.Content == "" && message.ImageURL == "" {
		return errors.New(models.CannotBeEmptyMessage) // 如果内容和图片链接都为空，则返回错误
	}

	result := database.DB.Create(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteMessage 根据 ID 删除留言
// 级联删除关联的内置评论
func DeleteMessage(id uint) error {
    tx := database.DB.Begin()
    // 先删除关联评论
    if err := tx.Where("message_id = ?", id).Delete(&models.Comment{}).Error; err != nil {
        tx.Rollback()
        return err
    }
    // 再删除消息本身
    var message models.Message
    result := tx.Delete(&message, id)
    if result.Error != nil {
        tx.Rollback()
        return result.Error
    }
    if result.RowsAffected == 0 {
        tx.Rollback()
        return gorm.ErrRecordNotFound
    }
    return tx.Commit().Error
}
// UpdateMessageContent 更新留言内容
func UpdateMessageContent(id uint, content string) error {
    // 开始事务
    tx := database.DB.Begin()
    
    // 先查询原消息
    var oldMessage models.Message
    if err := tx.First(&oldMessage, id).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 验证和处理新内容
    content = strings.TrimSpace(content)
    if content == "" && oldMessage.ImageURL == "" {
        tx.Rollback()
        return errors.New(models.CannotBeEmptyMessage)
    }
    
    // 创建新消息，保留原有信息
    newMessage := models.Message{
        ID:        oldMessage.ID, // 保持原有 ID
        UserID:    oldMessage.UserID,
        Username:  oldMessage.Username,
        Content:   content,
        ImageURL:  oldMessage.ImageURL,
        Private:   oldMessage.Private,
        CreatedAt: oldMessage.CreatedAt,
    }
    
    // 删除旧消息
    if err := tx.Unscoped().Delete(&oldMessage).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 创建新消息
    if err := tx.Create(&newMessage).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Error
}

// GetTotalMessages 获取消息总数
func GetTotalMessages() (int64, error) {
    var count int64
    if err := database.DB.Model(&models.Message{}).Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}
