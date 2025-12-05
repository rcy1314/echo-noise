package models

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type UserStatus struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	Status   Status `json:"status"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Username  string    `gorm:"type:varchar(100)" json:"username,omitempty"`
	ImageURL  string    `gorm:"type:text" json:"image_url,omitempty"`
	Private   bool      `gorm:"default:false" json:"private"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Notify    bool      `gorm:"default:false" json:"notify"` // 新增推送通知字段
	Pinned    bool      `gorm:"default:false" json:"pinned"`
	LikeCount int       `gorm:"default:0" json:"like_count"`
}

// MessageLike 点赞记录（用于幂等与取消点赞）
type MessageLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageID uint      `gorm:"index;not null" json:"message_id"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"`
	SessionID string    `gorm:"type:varchar(191);index" json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageID uint      `gorm:"index;not null" json:"message_id"`
	Nick      string    `gorm:"type:varchar(100)" json:"nick"`
	Mail      string    `gorm:"type:varchar(191)" json:"mail"`
	Link      string    `gorm:"type:varchar(191)" json:"link"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	Username           string     `gorm:"type:varchar(191);not null;uniqueIndex" json:"username"`
	Password           string     `gorm:"type:varchar(191);not null" json:"password"`
	IsAdmin            bool       `json:"is_admin"`
	Token              string     `gorm:"type:varchar(191)" json:"token"`
	AvatarURL          string     `gorm:"type:varchar(191)" json:"avatar_url"`
	Description        string     `gorm:"type:varchar(191)" json:"description"`
	Email              string     `gorm:"type:varchar(191)" json:"email"`
	EmailVerified      bool       `json:"email_verified"`
	EmailPending       string     `gorm:"type:varchar(191)" json:"-"`
	EmailVerifyCode    string     `gorm:"type:varchar(20)" json:"-"`
	EmailVerifyExpires *time.Time `json:"-"`
}

// 生成 Token 的工具函数
func GenerateToken(length int) string {
	b := make([]byte, length/2)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type Status struct {
	SysAdminID    uint         `json:"sys_admin_id"`
	Username      string       `json:"username"`
	Users         []UserStatus `json:"users"`
	TotalMessages int          `json:"total_messages"`
}

type UserSession struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin"`
	LoginTime time.Time `json:"login_time"`
}

type Setting struct {
	gorm.Model
	AllowRegistration bool `gorm:"default:true"`
}

type SiteConfig struct {
	gorm.Model
	SiteTitle        string `gorm:"type:varchar(100)"`
	SubtitleText     string `gorm:"type:varchar(191)"`
	AvatarURL        string `gorm:"type:varchar(191)"`
	Username         string `gorm:"type:varchar(50)"`
	Description      string `gorm:"type:varchar(191)"`
	Backgrounds      string `gorm:"type:text"`
	CardFooterTitle  string `gorm:"type:varchar(100)"`
	CardFooterLink   string `gorm:"type:varchar(100)"`
	PageFooterHTML   string `gorm:"type:text"`
	RSSTitle         string `gorm:"type:varchar(100)"`
	RSSDescription   string `gorm:"type:varchar(191)"`
	RSSAuthorName    string `gorm:"type:varchar(50)"`
	RSSFaviconURL    string `gorm:"type:varchar(191)"`
	WalineServerURL  string `gorm:"type:varchar(191)"`
	EnableGithubCard bool   `gorm:"default:false"`
	// PWA 配置
	PwaEnabled     bool   `gorm:"default:true"`
	PwaTitle       string `gorm:"type:varchar(100)"`
	PwaDescription string `gorm:"type:varchar(191)"`
	PwaIconURL     string `gorm:"type:varchar(191)"`
	// 主题默认模式: dark 或 light
	ContentThemeDefault string `gorm:"type:varchar(10)"`
	AnnouncementText    string `gorm:"type:varchar(191)"`
	AnnouncementEnabled bool   `gorm:"default:true"`
	Version             int    `json:"version"`
	SmtpEnabled         bool   `gorm:"default:false" json:"smtpEnabled"`
	SmtpDriver          string `gorm:"type:varchar(50)" json:"smtpDriver"`
	SmtpHost            string `gorm:"type:varchar(191)" json:"smtpHost"`
	SmtpPort            int    `json:"smtpPort"`
	SmtpUser            string `gorm:"type:varchar(191)" json:"smtpUser"`
	SmtpPass            string `gorm:"type:varchar(191)" json:"smtpPass"`
	SmtpFrom            string `gorm:"type:varchar(191)" json:"smtpFrom"`
	SmtpEncryption      string `gorm:"type:varchar(20)" json:"smtpEncryption"`
	SmtpTLS             bool   `gorm:"default:false" json:"smtpTLS"`
	// GitHub OAuth
	GithubOAuthEnabled bool   `gorm:"default:false"`
	GithubClientId     string `gorm:"type:varchar(191)"`
	GithubClientSecret string `gorm:"type:varchar(191)"`
	GithubCallbackURL  string `gorm:"type:varchar(191)"`
	// 云存储（S3/R2）备份设置
	StorageEnabled       bool   `gorm:"default:false"`
	StorageProvider      string `gorm:"type:varchar(20)"` // s3 或 r2
	StorageEndpoint      string `gorm:"type:varchar(191)"`
	StorageRegion        string `gorm:"type:varchar(100)"`
	StorageBucket        string `gorm:"type:varchar(191)"`
	StorageAccessKey     string `gorm:"type:varchar(191)"`
	StorageSecretKey     string `gorm:"type:varchar(191)"`
	StorageUsePathStyle  bool   `gorm:"default:true"`
	StoragePublicBaseURL string `gorm:"type:varchar(191)"`
	// 云存储自动同步
	StorageAutoSyncEnabled    bool       `gorm:"default:false"`
	StorageSyncMode           string     `gorm:"type:varchar(20)"` // instant 或 scheduled
	StorageSyncIntervalMinute int        `gorm:"default:15"`
	StorageLastSyncTime       *time.Time `json:"storageLastSyncTime"`
	// 音乐播放器配置（NeteaseMiniPlayer）
	MusicEnabled          bool   `gorm:"default:false"`
	MusicPlaylistId       string `gorm:"type:varchar(50)"`
	MusicSongId           string `gorm:"type:varchar(50)"`
	MusicPosition         string `gorm:"type:varchar(30)"` // static/top-left/top-right/bottom-left/bottom-right
	MusicTheme            string `gorm:"type:varchar(20)"` // auto/light/dark
	MusicLyric            bool   `gorm:"default:true"`
	MusicAutoplay         bool   `gorm:"default:false"`
	MusicDefaultMinimized bool   `gorm:"default:true"`
	MusicEmbed            bool   `gorm:"default:false"`
	MusicCssCdnURL        string `gorm:"type:varchar(255)"`
	MusicJsCdnURL         string `gorm:"type:varchar(255)"`
	// 评论系统配置
	CommentEnabled                bool   `gorm:"default:true"`
	CommentSystem                 string `gorm:"type:varchar(20)"` // builtin/waline/none/other
	CommentEmailEnabled           bool   `gorm:"default:false"`
	CommentLoginRequired          bool   `gorm:"default:true"`
	CommentEmailReplyName         string `gorm:"type:varchar(100)"`
	CommentEmailAdminPrefix       string `gorm:"type:varchar(50)"`
	CommentEmailReplyPrefix       string `gorm:"type:varchar(50)"`
	CommentEmailReplyTemplate     string `gorm:"type:text"`
	CommentEmailAdminTemplate     string `gorm:"type:text"`
	CommentEmailSiteURL           string `gorm:"type:varchar(191)"`
	CommentEmailReplyTemplateHTML string `gorm:"type:text"`
	CommentEmailAdminTemplateHTML string `gorm:"type:text"`
	// 扩展组件开关
	CalendarEnabled bool `gorm:"default:true"`
	TimeEnabled     bool `gorm:"default:true"`
	HitokotoEnabled bool `gorm:"default:true"`
	// 系统欢迎组件（左栏头像卡片专用，脱离用户资料）
	WelcomeAvatarURL   string `gorm:"type:varchar(255)"`
	WelcomeName        string `gorm:"type:varchar(100)"`
	WelcomeDescription string `gorm:"type:varchar(255)"`
	WelcomeUseAdmin    bool   `gorm:"default:true"`
	// 广告位配置（左侧）
	LeftAdEnabled     bool   `gorm:"default:true"`
	LeftAds           string `gorm:"type:text"`
	LeftAdsIntervalMs int    `gorm:"default:4000"`
}

func (s *SiteConfig) GetBackgroundsList() []string {
	if s.Backgrounds == "" {
		return []string{}
	}

	var backgrounds []string
	if err := json.Unmarshal([]byte(s.Backgrounds), &backgrounds); err != nil {
		// 如果解析失败，返回空数组
		return []string{}
	}
	return backgrounds
}
func UpdateMessage(id string, content string) error {
	// 先查询消息是否存在
	var message Message
	result := DB.First(&message, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("消息不存在")
		}
		return result.Error
	}

	// 更新消息内容
	result = DB.Model(&message).Updates(map[string]interface{}{
		"content": content,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新失败")
	}

	return nil
}
