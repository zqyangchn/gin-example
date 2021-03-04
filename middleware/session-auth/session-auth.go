package sessionauth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-example/models"
	"gin-example/pkg/app"
	"gin-example/pkg/errcode"
	"gin-example/pkg/gin-sessions"
)

// 使用 Cookie 保存 session
func EnableCookieSession() gin.HandlerFunc {
	store, err := ginsessions.NewRedisStore([]byte("secret"))
	if err != nil {
		panic(err)
	}
	store.SetRedisKeyPrefix("session:")
	store.SetMaxAge(86400 * 2)

	return ginsessions.Sessions("smp", store)
}

// session中间件
func AuthSessionMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := ginsessions.GetSession(c)

		sessionValue := session.Get("userId")
		if sessionValue == nil {
			appG := app.Gin{Context: c}
			appG.Response(http.StatusUnauthorized,
				errcode.CookieSessionError.WithDetails("Session is null, 未登录"),
				struct{}{})
			c.Abort()
			return
		}

		// 设置简单的变量
		c.Set("userId", sessionValue.(uint))

		c.Next()
		return
	}
}

// 注册和登陆时都需要保存sessions信息
func SaveAuthSession(c *gin.Context, id uint) error {
	session := ginsessions.GetSession(c)
	session.Set("userId", id)
	return session.Save()
}

// 退出时清除session
func ClearAuthSession(c *gin.Context) error {
	session := ginsessions.GetSession(c)
	session.Clear()
	return session.Save()
}

func HasSession(c *gin.Context) bool {
	session := ginsessions.GetSession(c)
	if sessionValue := session.Get("userId"); sessionValue == nil {
		return false
	}
	return true
}

func GetSessionUserId(c *gin.Context) uint {
	session := ginsessions.GetSession(c)
	sessionValue := session.Get("userId")
	if sessionValue == nil {
		return 0
	}
	return sessionValue.(uint)
}

func GetUserSession(c *gin.Context) map[string]interface{} {
	hasSession := HasSession(c)
	userName := ""
	if hasSession {
		userId := GetSessionUserId(c)
		user, _ := models.UserDetail(userId)
		userName = user.Name
	}
	data := make(map[string]interface{})
	data["hasSession"] = hasSession
	data["userName"] = userName
	return data
}
