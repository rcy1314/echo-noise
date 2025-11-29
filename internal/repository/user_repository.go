package repository

import (
	"sync"
	"time"

	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
)

type userCacheItem struct {
	user      *models.User
	timestamp time.Time
}

var (
	userCacheMap = make(map[uint]userCacheItem)
	cacheMutex   sync.RWMutex
	cacheExpiry  = 30 * time.Minute
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	err := database.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := database.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(id uint) (*models.User, error) {
	cacheMutex.RLock()
	if item, exists := userCacheMap[id]; exists {
		if time.Since(item.timestamp) < cacheExpiry {
			cacheMutex.RUnlock()
			return item.user, nil
		}
	}
	cacheMutex.RUnlock()

	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	cacheMutex.Lock()
	userCacheMap[id] = userCacheItem{
		user:      &user,
		timestamp: time.Now(),
	}
	cacheMutex.Unlock()

	return &user, nil
}

func GetSysAdmin() (*models.User, error) {
	var user models.User
	err := database.DB.Where("is_admin = ?", true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func clearUserCache(id uint) {
	cacheMutex.Lock()
	delete(userCacheMap, id)
	cacheMutex.Unlock()
}

func UpdateUser(user *models.User) error {
	if user == nil {
		return nil
	}
	err := database.DB.Save(user).Error
	if err != nil {
		return err
	}
	clearUserCache(user.ID)
	return nil
}
func GetSettingByID(id uint) (*models.Setting, error) {
    var setting models.Setting
    result := database.DB.First(&setting, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &setting, nil
}

func UpdateSetting(setting *models.Setting, updates map[string]interface{}) error {
    result := database.DB.Model(setting).Updates(updates)
    return result.Error
}
func UpdateUserField(userID uint, field string, value interface{}) error {
    return database.DB.Model(&models.User{}).Where("id = ?", userID).Update(field, value).Error
}
func UpdateUserToken(userID uint, token string) error {
    err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("token", token).Error
    if err != nil {
        return err
    }
    clearUserCache(userID)  // 清除用户缓存
    return nil
}
// DeleteUser 删除用户
func DeleteUser(id uint) error {
    err := database.DB.Delete(&models.User{}, id).Error
    if err != nil {
        return err
    }
    clearUserCache(id)
    return nil
}

// CountAdmins 统计管理员数量
func CountAdmins() (int64, error) {
    var count int64
    err := database.DB.Model(&models.User{}).Where("is_admin = ?", true).Count(&count).Error
    if err != nil {
        return 0, err
    }
    return count, nil
}

// BatchCreateUsers 批量创建用户
func BatchCreateUsers(users []*models.User) error {
    return database.DB.Create(&users).Error
}

// BatchUpdateUsers 批量更新用户
func BatchUpdateUsers(users []*models.User) error {
    for _, user := range users {
        if err := UpdateUser(user); err != nil {
            return err
        }
    }
    return nil
}

// ClearExpiredCache 清理过期缓存
func ClearExpiredCache() {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()
    
    now := time.Now()
    for id, item := range userCacheMap {
        if now.Sub(item.timestamp) > cacheExpiry {
            delete(userCacheMap, id)
        }
    }
}
