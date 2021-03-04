package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// swagger
	_ "gin-example/docs"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware

	"gin-example/middleware/session-auth"
	"gin-example/middleware/zaplogger"
	"gin-example/pkg/logging"
	"gin-example/pkg/setting"
	"gin-example/routers/api"
	"gin-example/routers/api/v1"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(zaplogger.Ginzap(logging.Logger))
	r.Use(zaplogger.RecoveryWithZap(logging.Logger, true))

	gin.DebugPrintRouteFunc = logging.GinDebugPrintRouteZapLoggerFunc

	/*
		################  swagger  ################
		go get -u github.com/swaggo/swag/cmd/swag
		go get -u github.com/swaggo/gin-swagger
		go get -u github.com/swaggo/files
		go get -u github.com/alecthomas/template

		swag init -g http-service.go

		"github.com/swaggo/files"       // swagger embed files
		"github.com/swaggo/gin-swagger" // gin-swagger middleware
		_ "i-morefun.net/morefun/go-programming-tour-book/http-service/docs"
	*/
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// jwt auth
	r.POST("/auth", api.GetAuth)
	r.GET("/reverse/solution/jwt", api.ReverseSolutionJWT)

	r.GET("/stream", api.Stream)

	// 上传文件
	r.POST("/upload/file", api.UploadFile)
	r.StaticFS("/static", http.Dir(setting.AppSetting.UploadSavePath))

	sr := r.Group("/", sessionauth.EnableCookieSession())
	{
		// 新建用户
		sr.POST("/register", api.Register)
		// 登陆
		sr.POST("/login", api.Login)
		// 新建用户
		sr.GET("/logout", api.Logout)

		authorized := sr.Group("/", sessionauth.AuthSessionMiddle())
		{
			// 错误码汇总
			authorized.GET("/error/message", api.GetErrorMessages)

			// apiv1
			apiv1 := authorized.Group("/api/v1")
			{
				//获取用户列表
				apiv1.GET("/users", api.GetUsers)

				//获取标签列表
				apiv1.GET("/tags", v1.GetTags)
				//新建标签
				apiv1.POST("/tags", v1.AddTag)
				//更新指定标签
				apiv1.PUT("/tags/:id", v1.EditTag)
				//删除指定标签
				apiv1.DELETE("/tags/:id", v1.DeleteTag)
			}
		}
	}

	return r
}
