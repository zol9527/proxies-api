// Author       :loyd
// Date         :2025-03-30 12:41:41
// LastEditors  :loyd
// LastEditTime :2025-03-31 00:05:46
// Description  :程序入口点
//

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"proxies-api/internal/api"
	"proxies-api/internal/config"
	"proxies-api/internal/log"
	"proxies-api/internal/pool"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Function: main
//
// Purpose: 程序入口函数，初始化并启动代理IP池系统
//
// Parameters:
//   - 无
//
// Returns: 无
func main() {
	fmt.Println("🚀 代理IP池管理系统启动中...")

	// 初始化配置和IP池
	pool_instance := pool.InitializeConfigAndPool()
	if pool_instance == nil {
		fmt.Println("❌ 初始化IP池失败，程序退出")
		os.Exit(1)
	}

	// 设置日志系统
	setupLogSystem()

	// 初始化并启动HTTP服务
	router, server := setupHTTPServer()

	// 启动IP获取服务
	go startIPFetchService()

	// 监听系统信号以优雅关闭
	setupGracefulShutdown(router, server)

}

// Function: setupLogSystem
//
// Purpose: 配置并初始化日志系统
//
// Parameters:
//   - 无
//
// Returns: 无
func setupLogSystem() {
	// 创建日志目录
	log_dir := "./logs"
	os.MkdirAll(log_dir, 0755)

	// 配置日志系统
	log_config := log.LogConfig{
		LogLevel:   "debug",
		LogFile:    log_dir + "/proxy_pool.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}

	log.SetupLogger(log_config)
	log.Info("📝 日志系统初始化完成")
}

// Function: setupHTTPServer
//
// Purpose: 设置并启动HTTP服务器
//
// Parameters:
//   - 无
//
// Returns: (*gin.Engine, *http.Server) - Gin引擎实例和HTTP服务器实例
func setupHTTPServer() (*gin.Engine, *http.Server) {
	// 设置gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由器
	router := gin.Default()

	// 注册中间件
	router.Use(api.RateLimitMiddleware())

	// 注册路由
	api.RegisterUrlRoutes(router)

	// 获取配置的端口
	port := config.APP_CONFIG.Server.Port
	if port == 0 {
		port = 8000 // 默认端口
	}

	// 创建HTTP服务器
	addr := ":" + strconv.Itoa(port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 启动HTTP服务
	go func() {
		log.Info("🌐 HTTP服务器启动于 %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("❌ HTTP服务器启动失败: %v", err)
		}
	}()

	return router, server
}

// Function: startIPFetchService
//
// Purpose: 启动IP获取服务
//
// Parameters:
//   - 无
//
// Returns: 无
func startIPFetchService() {
	log.Info("🔍 IP获取服务启动")

	// 延迟一秒等待其他服务初始化完成
	time.Sleep(1 * time.Second)

	// 开始更新IP库
	api.RefreshIPRegistry()
}

// Function: setupGracefulShutdown
//
// Purpose: 设置优雅关闭机制
//
// Parameters:
//   - router (*gin.Engine): Gin引擎实例
//   - server (*http.Server): HTTP服务器实例
//
// Returns: 无
func setupGracefulShutdown(router *gin.Engine, server *http.Server) {
	// 创建系统信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-quit

	log.Info("🛑 接收到关闭信号，系统正在优雅关闭...")

	// 创建一个超时上下文用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭所有现有连接，但等待请求完成
	if err := server.Shutdown(ctx); err != nil {
		log.Error("❌ HTTP服务器关闭出错: %v", err)
	} else {
		log.Info("✅ HTTP服务器已优雅关闭")
	}

	// 停止IP获取服务
	if api.StopIPFetchService() {
		log.Info("✅ IP获取服务已停止")
	}

	// 关闭IP池和其他资源
	pool_instance := pool.AccessPool()
	if pool_instance != nil {
		pool_instance.StopCleanTicker()
		log.Info("⏹️ IP清理定时器已停止")
	}

	log.Info("👋 系统已安全关闭")
}
