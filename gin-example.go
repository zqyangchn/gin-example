package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-example/models"
	"gin-example/pkg/cache"
	"gin-example/pkg/database"
	"gin-example/pkg/logging"
	"gin-example/pkg/setting"
	"gin-example/routers"
)

func init() {
	// 配置初始化
	if err := setting.Setup(); err != nil {
		panic(err)
	}
	// 初始化日志
	logging.Setup()
	// 初始化缓存
	if err := cache.Setup(); err != nil {
		logging.Logger.Fatal("cache initialization failed", zap.Error(err))
	}
	// 初始化数据库
	if err := database.Setup(); err != nil {
		logging.Logger.Fatal("database initialization failed", zap.Error(err))
	}
	// 数据库表结构变更
	if err := models.Setup(); err != nil {
		logging.Logger.Fatal("models initialization failed", zap.Error(err))
	}

}

// @title API swagger
// @version 1.0
// @description smp ops system.
// @termsOfService http://github.com/zqyangchn

// @contact.name API Support
// @contact.url http://github.com/zqyangchn
// @contact.email zqyangchn@gmail.com
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	router := routers.NewRouter()
	srv := &http.Server{
		Addr:           ":" + setting.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logging.Logger.Info("web server starting ...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatal("web Server start Failed", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 15 秒的超时时间）
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt)
	<-osSignal

	// 启动服务器关闭流程
	logging.Logger.Info("shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// srv.Shutdown(ctx) 关闭服务器监听端口, 不再接受新的请求
	if err := srv.Shutdown(ctx); err != nil {
		logging.Logger.Fatal("Server Shutdown:", zap.Error(err))
	}
	logging.Logger.Info("server shutdown completed !")
}
