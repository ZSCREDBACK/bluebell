package main

import (
	"context"
	"fmt"
	"goScaffold/dao/mysql"
	"goScaffold/dao/redis"
	"goScaffold/logger"
	"goScaffold/routes"
	"goScaffold/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Go Web项目通用脚手架模板-升级版

func main() {
	// 1.加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}
	fmt.Println("Init config success")

	// 2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		return
	}
	defer zap.L().Sync() // 将缓冲区的日志追加到文件中
	zap.L().Debug("Init logger success")

	// 3.初始化MySQL数据库连接
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("初始化数据库失败: %v\n", err)
		return
	}
	defer mysql.Close()
	zap.L().Debug("connect mysql success")

	// 4.初始化Redis数据库连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("初始化Redis失败: %v\n", err)
		return
	}
	defer redis.Close()
	zap.L().Debug("connect redis success")

	// 5.注册路由
	r := routes.Setup(settings.Conf)

	// 6.启动服务(添加优雅关机的功能,Ctrl+C或kill-2)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			settings.Conf.Host,
			settings.Conf.Port,
		),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Listen: ", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
