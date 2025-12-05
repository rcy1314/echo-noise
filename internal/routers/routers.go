package routers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/controllers"
	"github.com/lin-snow/ech0/internal/middleware"
	"github.com/lin-snow/ech0/pkg"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 使用 pkg 中的 session 初始化
	pkg.InitSession(r)
	// 配置 CORS
	corsConfig := cors.DefaultConfig()
	// 动态按环境变量或放通所有来源（支持反代与跨域小组件）
	var allowList []string
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		for _, o := range strings.Split(origins, ",") {
			s := strings.TrimSpace(o)
			if s != "" {
				allowList = append(allowList, s)
			}
		}
	}
	if len(allowList) == 0 {
		corsConfig.AllowOriginFunc = func(origin string) bool { return true }
	} else {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			for _, o := range allowList {
				if origin == o {
					return true
				}
			}
			return false
		}
	}
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"X-Requested-With",
		"Accept",
		"Device-Type",
		"Authorization", // 新增授权头
		"Cache-Control",
		"Pragma",
		"Referer",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 86400

	// 不再使用 AllowOrigins 列表，统一使用 AllowOriginFunc 做灵活匹配

	r.Use(cors.New(corsConfig))

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
	vidDir := pickDir([]string{
		"./data/video",
		filepath.Join(wd, "data/video"),
		filepath.Join(exeDir, "data/video"),
		"/data/video",
		"/app/data/video",
	}, "./data/video")
	r.Static("/api/images", imgDir)
	r.Static("/video", vidDir)
	// 常用静态文件已在上方映射

	// API 路由组
	api := r.Group("/api")

	// 消息详情页路由（移到 API 组外）
	r.GET("/m/:id", controllers.GetMessagePage)

	// RSS 路由
	r.GET("/rss", controllers.GenerateRSS)                                               // 保持原有的 RSS 订阅链接
	api.POST("/rss/refresh", middleware.SessionAuthMiddleware(), controllers.RefreshRSS) // 添加刷新 RSS 的路由

	// PWA Manifest 路由
	r.GET("/manifest.webmanifest", controllers.GetWebManifest)

	// 公共路由
	api.GET("", controllers.GetStatus)
	api.GET("/frontend/config", controllers.GetFrontendConfig)
	api.GET("/settings", controllers.GetFrontendConfig)
	api.POST("/login", controllers.Login)
	api.POST("/register", controllers.Register)
	api.GET("/status", controllers.GetStatus)
	api.GET("/captcha", controllers.GetCaptcha)
	// api.GET("/config", controllers.GetFrontendConfig)
	api.GET("/messages", controllers.GetMessages)
	api.GET("/messages/:id", controllers.GetMessage)
	api.POST("/messages/page", controllers.GetMessagesByPage)
	api.GET("/messages/page", controllers.GetMessagesByPage)
	api.GET("/messages/calendar", controllers.GetMessagesCalendar) // 新增热力图专用路由
	api.GET("/messages/search", controllers.SearchMessages)        // 新增搜索消息路由
	api.GET("/version/check", controllers.CheckVersion)            // 添加版本检查路由
	api.GET("/version", controllers.GetVersion)                    // 当前运行版本（镜像标签/环境变量）
	// GitHub OAuth
	api.GET("/oauth/github/login", controllers.GithubLogin)
	r.GET("/oauth/github/callback", controllers.GithubCallback)
	api.POST("/password/forgot", controllers.PasswordForgot)

	// 添加标签和图像相关路由
	api.GET("/messages/tags/:tag", controllers.GetMessagesByTag)         // 获取指定标签的消息
	api.GET("/messages/tags", controllers.GetAllTags)                    // 获取所有标签列表
	api.GET("/messages/images", controllers.GetAllImages)                // 获取所有图片列表
	api.POST("/messages/:id/like", controllers.IncrementMessageLike)     // 点赞接口
	api.POST("/messages/:id/like/toggle", controllers.ToggleMessageLike) // 点赞切换
	api.GET("/guestbook/message", controllers.GetGuestbookMessageID)     // 获取留言板消息ID

	// 需要鉴权的路由
	authRoutes := api.Group("")
	authRoutes.Use(middleware.SessionAuthMiddleware())

	// 添加 token 认证的路由组
	tokenAuth := api.Group("/token")
	tokenAuth.Use(middleware.TokenAuthMiddleware()) // 使用 TokenAuthMiddleware
	{
		tokenAuth.POST("/messages", controllers.PostMessage)
		tokenAuth.PUT("/messages/:id", controllers.UpdateMessage)
		tokenAuth.PUT("/messages/:id/pin", controllers.UpdateMessagePinned)
		tokenAuth.DELETE("/messages/:id", controllers.DeleteMessage)
		tokenAuth.PUT("/settings", controllers.UpdateSetting)
	}
	// 需要鉴权的消息操作路由
	messages := authRoutes.Group("/messages")
	{
		messages.POST("", controllers.PostMessage)
		messages.PUT("/:id", controllers.UpdateMessage)
		messages.PUT("/:id/pin", controllers.UpdateMessagePinned)
		messages.DELETE("/:id", controllers.DeleteMessage)
	}

	// 评论系统（内置）公共路由
	api.GET("/messages/:id/comments", controllers.GetComments)
	api.POST("/messages/:id/comments", controllers.PostComment)
	// 管理员评论列表（提供公共路径，附加会话中间件以注入用户上下文）
	api.GET("/comments", middleware.SessionAuthMiddleware(), controllers.ListComments)
	// 评论删除（管理员）
	authRoutes.DELETE("/messages/:id/comments/:cid", controllers.DeleteComment)
	// 一次性回填评论 parent_id（管理员）
	authRoutes.POST("/comments/backfill", controllers.BackfillCommentParents)
	// 管理员评论列表管理（已在公共组注册路径，函数内部鉴权）
	// 添加推送配置路由
	notify := authRoutes.Group("/notify")
	{
		notify.POST("/test", controllers.TestNotify)        // 测试推送
		notify.POST("/send", controllers.SendNotify)        // 新增：实际推送路由
		notify.GET("/config", controllers.GetNotifyConfig)  // 获取配置
		notify.PUT("/config", controllers.SaveNotifyConfig) // 保存配置
	}

	email := authRoutes.Group("/email")
	{
		email.POST("/test", controllers.EmailTest)
	}

	// 数据库备份相关路由
	backup := authRoutes.Group("/backup")
	{
		backup.GET("/download", controllers.HandleBackupDownload)
		backup.POST("/restore", controllers.HandleBackupRestore)
		backup.POST("/storage/upload", controllers.HandleBackupUploadToURL)
		backup.POST("/storage/restore", controllers.HandleBackupRestoreFromURL)
		backup.POST("/storage/presign/upload", controllers.HandleBackupPresignUpload)
		backup.POST("/storage/presign/download", controllers.HandleBackupPresignDownload)
		backup.POST("/storage/sync-now", controllers.HandleBackupSyncNow)
	}

	// 图片上传路由
	authRoutes.POST("/images/upload", controllers.UploadImage) // 上传图片
	// 新增：视频上传路由（改为单数 video）
	authRoutes.POST("/video/upload", controllers.UploadVideo) // 上传视频

	// 附件管理路由
	attachments := authRoutes.Group("/attachments")
	{
		attachments.GET("/images", controllers.ListImageAttachments)
		attachments.GET("/images/", controllers.ListImageAttachments)
		attachments.GET("/video", controllers.ListVideoAttachments)
		attachments.GET("/video/", controllers.ListVideoAttachments)
		attachments.DELETE("/images/:name", middleware.AdminAuthMiddleware(), controllers.DeleteImageAttachment)
		attachments.DELETE("/images/:name/", middleware.AdminAuthMiddleware(), controllers.DeleteImageAttachment)
		attachments.DELETE("/video/:name", middleware.AdminAuthMiddleware(), controllers.DeleteVideoAttachment)
		attachments.DELETE("/video/:name/", middleware.AdminAuthMiddleware(), controllers.DeleteVideoAttachment)
	}

	// 用户相关路由
	user := authRoutes.Group("/user")
	{
		user.GET("", controllers.GetUserInfo)
		user.PUT("/change_password", controllers.ChangePassword)
		user.PUT("/update", controllers.UpdateUser)
		user.PUT("/admin", controllers.UpdateUserAdmin)
		user.DELETE("", controllers.DeleteUser)
		user.POST("/logout", controllers.Logout) // 添加退出登录路由
		user.POST("/reset_password", controllers.AdminResetPassword)
		// 添加 Token 相关路由
		user.GET("/token", controllers.GetUserToken)
		user.POST("/token/regenerate", controllers.RegenerateUserToken)
	}

	// 设置路由
	authRoutes.PUT("/settings", controllers.UpdateSetting)

	// 显式 /status 返回 SPA 入口，避免目录重定向影响
	r.GET("/status", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	// 先映射关键静态文件，避免被 SPA Fallback 覆盖
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.StaticFile("/favicon-32x32.png", "./public/favicon-32x32.png")
	r.StaticFile("/android-chrome-192x192.png", "./public/android-chrome-192x192.png")
	r.StaticFile("/favicon.svg", "./public/favicon.svg")
	// manifest 与 SW 使用控制器路由，避免重复注册

	// 使用静态中间件托管根目录，支持 SPA Fallback
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	// 显式映射 Nuxt 资源目录和常用静态
	r.Static("/_nuxt", "./public/_nuxt")
	r.Static("/assets", "./public/assets")

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/m/") ||
			strings.HasPrefix(path, "/messages/") ||
			path == "/" ||
			!strings.HasPrefix(path, "/api") {
			c.File("./public/index.html")
		} else {
			c.JSON(http.StatusNotFound, gin.H{"code": 0, "msg": "接口不存在"})
		}
	})

	return r
}

func pickDir(candidates []string, fallback string) string {
	for _, d := range candidates {
		if d == "" {
			continue
		}
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			return d
		}
	}
	return fallback
}
