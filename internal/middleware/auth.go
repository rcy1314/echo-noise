package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "github.com/lin-snow/ech0/internal/dto"
    "github.com/lin-snow/ech0/internal/models"
    "github.com/lin-snow/ech0/internal/database"
)

func SessionAuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        userID := session.Get("user_id")

        if userID == nil {
            // 定义公共路由
            publicPaths := []string{
                "/api/messages/page",
                "/api/messages/",
                "/api/messages",
                "/api/messages/search",
                "/api/messages/tags",
                "/api/messages/tags/",
                "/api/messages/images",
                "/api/messages/calendar",
                "/api/frontend/config",
                "/api/status",
                "/api/version/check",
            }

            // 检查是否是公共路由
            for _, path := range publicPaths {
                if strings.HasPrefix(ctx.Request.URL.Path, path) {
                    ctx.Set("user_id", uint(0))
                    ctx.Next()
                    return
                }
            }

            ctx.JSON(http.StatusUnauthorized, dto.Fail[any]("未登录或登录已过期"))
            ctx.Abort()
            return
        }

        // 将用户信息存储到上下文中
        ctx.Set("user_id", userID.(uint))
        ctx.Set("username", session.Get("username"))
        ctx.Set("is_admin", session.Get("is_admin"))
        ctx.Next()
    }
}

func AdminAuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        isAdmin := session.Get("is_admin")

        if isAdmin == nil || !isAdmin.(bool) {
            ctx.JSON(http.StatusForbidden, dto.Fail[any]("需要管理员权限"))
            ctx.Abort()
            return
        }

        ctx.Next()
    }
}

func TokenAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.JSON(http.StatusOK, dto.Fail[any]("未提供认证信息"))
            c.Abort()
            return
        }

        // 提取 token
        var token string
        if strings.HasPrefix(auth, "Bearer ") {
            token = strings.TrimPrefix(auth, "Bearer ")
        } else {
            token = auth
        }

        db, err := database.GetDB()
        if err != nil {
            c.JSON(http.StatusOK, dto.Fail[any]("系统错误"))
            c.Abort()
            return
        }

        // 查询用户
        var user models.User
        if err := db.Where("token = ?", token).First(&user).Error; err != nil {
            c.JSON(http.StatusOK, dto.Fail[any]("无效的token"))
            c.Abort()
            return
        }

        // 检查 token 是否为空
        if user.Token == "" {
            c.JSON(http.StatusOK, dto.Fail[any]("token已失效"))
            c.Abort()
            return
        }

        // 设置用户信息到上下文
        c.Set("user_id", user.ID)
        c.Set("username", user.Username)
        c.Set("is_admin", user.IsAdmin)
        c.Next()
    }
}
