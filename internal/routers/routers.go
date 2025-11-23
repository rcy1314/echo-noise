package routers

import (
	"net/http"
	"os"
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
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"X-Requested-With",
		"Accept",
		"Device-Type",
		"Authorization", // 新增授权头
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 86400

	allowed := []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		// 支持以逗号分隔的多个来源
		var list []string
		for _, o := range strings.Split(origins, ",") {
			s := strings.TrimSpace(o)
			if s != "" {
				list = append(list, s)
			}
		}
		if len(list) > 0 {
			allowed = list
		}
	}

	// 在开发模式下，限定来源为本地前端；生产下建议通过环境变量配置
	if config.Config.Server.Mode == "debug" {
		corsConfig.AllowOrigins = allowed
	} else {
		corsConfig.AllowOrigins = allowed
	}

	r.Use(cors.New(corsConfig))

	// 映射静态文件目录
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.Static("/api/images", "./data/images")
	r.Static("/video", "./data/video")

	// API 路由组
	api := r.Group("/api")

	// 消息详情页路由（移到 API 组外）
	r.GET("/m/:id", controllers.GetMessagePage)

	// RSS 路由
	r.GET("/rss", controllers.GenerateRSS)                                               // 保持原有的 RSS 订阅链接
	api.POST("/rss/refresh", middleware.SessionAuthMiddleware(), controllers.RefreshRSS) // 添加刷新 RSS 的路由

	// 公共路由
	api.GET("/frontend/config", controllers.GetFrontendConfig)
	api.POST("/login", controllers.Login)
	api.POST("/register", controllers.Register)
	api.GET("/status", controllers.GetStatus)
	// api.GET("/config", controllers.GetFrontendConfig)
	api.GET("/messages", controllers.GetMessages)
	api.GET("/messages/:id", controllers.GetMessage)
	api.POST("/messages/page", controllers.GetMessagesByPage)
	api.GET("/messages/page", controllers.GetMessagesByPage)
	api.GET("/messages/calendar", controllers.GetMessagesCalendar) // 新增热力图专用路由
	api.GET("/messages/search", controllers.SearchMessages)        // 新增搜索消息路由
	api.GET("/version/check", controllers.CheckVersion)            // 添加版本检查路由

	// 添加标签和图像相关路由
	api.GET("/messages/tags/:tag", controllers.GetMessagesByTag) // 获取指定标签的消息
	api.GET("/messages/tags", controllers.GetAllTags)            // 获取所有标签列表
	api.GET("/messages/images", controllers.GetAllImages)        // 获取所有图片列表

	// 需要鉴权的路由
	authRoutes := api.Group("")
	authRoutes.Use(middleware.SessionAuthMiddleware())

	// 添加 token 认证的路由组
	tokenAuth := api.Group("/token")
	tokenAuth.Use(middleware.TokenAuthMiddleware()) // 使用 TokenAuthMiddleware
	{
		tokenAuth.POST("/messages", controllers.PostMessage)
	}
	// 需要鉴权的消息操作路由
	messages := authRoutes.Group("/messages")
	{
		messages.POST("", controllers.PostMessage)
		messages.PUT("/:id", controllers.UpdateMessage)
		messages.DELETE("/:id", controllers.DeleteMessage)
	}
	// 添加推送配置路由
	notify := authRoutes.Group("/notify")
	{
		notify.POST("/test", controllers.TestNotify)        // 测试推送
		notify.POST("/send", controllers.SendNotify)        // 新增：实际推送路由
		notify.GET("/config", controllers.GetNotifyConfig)  // 获取配置
		notify.PUT("/config", controllers.SaveNotifyConfig) // 保存配置
	}

	// 数据库备份相关路由
	backup := authRoutes.Group("/backup")
	{
		backup.GET("/download", controllers.HandleBackupDownload)
		backup.POST("/restore", controllers.HandleBackupRestore)
	}

	// 图片上传路由
	authRoutes.POST("/images/upload", controllers.UploadImage) // 上传图片
	// 新增：视频上传路由
	authRoutes.POST("/videos/upload", controllers.UploadVideo) // 上传视频

	// 用户相关路由
	user := authRoutes.Group("/user")
	{
		user.GET("", controllers.GetUserInfo)
		user.PUT("/change_password", controllers.ChangePassword)
		user.PUT("/update", controllers.UpdateUser)
		user.PUT("/admin", controllers.UpdateUserAdmin)
		user.POST("/logout", controllers.Logout) // 添加退出登录路由
		// 添加 Token 相关路由
		user.GET("/token", controllers.GetUserToken)
		user.POST("/token/regenerate", controllers.RegenerateUserToken)
	}

	// 设置路由
	authRoutes.PUT("/settings", controllers.UpdateSetting)

	// 404 处理
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
