This file is a merged representation of the entire codebase, combined into a single document by Repomix.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Files are sorted by Git change count (files with more changes are at the bottom)

## Additional Info

# Directory Structure
```
cmd/
  main.go
config/
  proxy_pool_config.toml
internal/
  api/
    fetch.go
    http.go
    rate_test.go
    rate.go
  config/
    config_test.go
    config.go
  log/
    log.go
  models/
    domaininfo_test.go
    domaininfo.go
    ipdetails.go
    pool.go
  pool/
    init.go
.gitignore
go.mod
README.md
repomix.config.json
```

# Files

## File: cmd/main.go
````go
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
````

## File: config/proxy_pool_config.toml
````toml
[server]
port = 8000

# 代理池配置
[pool]
size = 100

# 自定义代理源
[proxy_source]
urls = [
    "http://example.com/proxy1"
]
request_frequency = 10
local_path = [
    "/var/proxy_pool/proxy.txt",
    "/var/proxy_pool/proxy2.txt",
    "/var/proxy_pool/proxy3.txt",
]
read_frequency = 10

# 单独对每个 ip 设置, delay 也可以使用 URL 参数动态设置
[ip]
lifetime_seconds = 900
delay = [0, 5]
````

## File: internal/api/http.go
````go
// Author       :loyd
// Date         :2025-03-30 18:19:08
// LastEditors  :loyd
// LastEditTime :2025-03-30 18:21:14
// Description  :通过 URL 参数传递数组的 HTTP 服务
//

package api

import (
	"strconv"
	"strings"

	"proxies-api/internal/models"
	"proxies-api/internal/pool"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Domain string `json:"domain"`
	Delay  []int  `json:"delay,omitempty"`
}

/**
 * Function: API_Get
 *
 * Purpose: 通过 URL 参数获取代理 IP 信息，包括处理数组形式的 delay 参数
 *
 * Parameters:
 *   - c (*gin.Context): Gin 框架的上下文对象
 *
 * Returns:
 */
func API_Get(c *gin.Context) {
	// 从 URL 获取域名参数
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(400, gin.H{"error": "❌ 缺少域名参数"})
		return
	}

	// 获取 delay 参数数组，支持两种方式：
	// 1. 逗号分隔：?delay=1,5
	// 2. 重复参数：?delay=1&delay=5
	var delay_array []int

	// 方式1：通过逗号分隔的方式获取数组
	delay_csv := c.Query("delay")
	if delay_csv != "" {
		delay_strs := strings.Split(delay_csv, ",")
		for _, str := range delay_strs {
			if val, err := strconv.Atoi(str); err == nil {
				delay_array = append(delay_array, val)
			}
		}
	} else {
		// 方式2：通过重复参数的方式获取数组
		delay_strs := c.QueryArray("delay")
		for _, str := range delay_strs {
			if val, err := strconv.Atoi(str); err == nil {
				delay_array = append(delay_array, val)
			}
		}
	}

	// 处理查询逻辑
	client_pool := pool.AccessPool()
	domain_info, exists := client_pool.GetDomainInfo(domain)
	if !exists {
		// 添加新域名
		client_pool.AddDomain(domain)
		domain_info, _ = client_pool.GetDomainInfo(domain)

		// 如果有传递 delay 参数，设置延迟
		if len(delay_array) > 0 {
			domain_info.SetDelay(delay_array)
		}

		// 复制 origin ips
		origins := client_pool.GetOriginIPs()
		for _, ip := range origins {
			ip_info := models.NewIpDetails(ip.GetIP(), ip.GetPort())
			domain_info.Add(ip_info)
		}
	}

	// 返回一个 IP
	fetched_ip, err := domain_info.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": "❌ 获取 IP 失败"})
		return
	}

	// 返回结果
	c.JSON(200, gin.H{
		"ip":    fetched_ip.GetIP(),
		"port":  fetched_ip.GetPort(),
		"delay": delay_array, // 返回解析的 delay 数组，方便调试
	})
}

func API_Status(c *gin.Context) {

	status := GetProxyPoolStatus()
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "✅ 获取状态成功",
		"status": status,
	})
}

/**
 * Function: RegisterUrlRoutes
 *
 * Purpose: 注册使用 URL 参数的 HTTP 路由
 *
 * Parameters:
 *   - router (*gin.Engine): Gin 引擎实例
 *
 * Returns:
 */
func RegisterUrlRoutes(router *gin.Engine) {
	router.GET("/proxy", API_Get)
	router.GET("/status", API_Status)
}
````

## File: internal/api/rate_test.go
````go
// Author       :loyd
// Date         :2025-03-30 11:24:58
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:28:50
// Description  :测试基于IP的请求限流中间件

package api

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

// 测试常量
const (
	TEST_HOST        = "example.com"
	TEST_RATE_LIMIT  = 5
	TEST_BURST_SIZE  = 10
	TEST_MINI_EXPIRY = 1 // 测试用的较短过期时间（分钟）
)

// Function: TestNewIPRateLimiter
//
// Purpose: 测试限流器初始化功能
//
// Parameters:
//   - t (*testing.T): 测试对象
//
// Returns: 无
func TestNewIPRateLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(TEST_RATE_LIMIT), TEST_BURST_SIZE)

	assert.NotNil(t, limiter, "📋 限流器不应为空")
	assert.NotNil(t, limiter.hosts, "📋 hosts映射不应为空")
	assert.NotNil(t, limiter.mu, "📋 互斥锁不应为空")
	assert.Equal(t, rate.Limit(TEST_RATE_LIMIT), limiter.rate_limit, "📋 限流率应匹配")
	assert.Equal(t, TEST_BURST_SIZE, limiter.burst_size, "📋 突发大小应匹配")
}

// Function: TestGetLimiter
//
// Purpose: 测试获取主机限流器功能
//
// Parameters:
//   - t (*testing.T): 测试对象
//
// Returns: 无
func TestGetLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(TEST_RATE_LIMIT), TEST_BURST_SIZE)

	// 首次获取会创建新的限流器
	host_limiter := limiter.GetLimiter(TEST_HOST)
	assert.NotNil(t, host_limiter, "📋 获取的限流器不应为空")

	// 再次获取应返回相同的限流器
	host_limiter2 := limiter.GetLimiter(TEST_HOST)
	assert.Equal(t, host_limiter, host_limiter2, "📋 应返回相同的限流器实例")

	// 检查hosts映射是否包含此主机
	limiter.mu.RLock()
	_, exists := limiter.hosts[TEST_HOST]
	limiter.mu.RUnlock()
	assert.True(t, exists, "📋 hosts映射应包含请求的主机")
}

// Function: TestConcurrentGetLimiter
//
// Purpose: 测试并发获取限流器场景
//
// Parameters:
//   - t (*testing.T): 测试对象
//
// Returns: 无
func TestConcurrentGetLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(TEST_RATE_LIMIT), TEST_BURST_SIZE)
	const num_goroutines = 100

	var wg sync.WaitGroup
	wg.Add(num_goroutines)

	for i := 0; i < num_goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = limiter.GetLimiter(TEST_HOST)
		}()
	}

	wg.Wait()

	// 验证无论多少goroutine，只创建了一个限流器
	limiter.mu.RLock()
	assert.Equal(t, 1, len(limiter.hosts), "📋 应只创建一个限流器实例")
	limiter.mu.RUnlock()
}

// Function: TestRateLimiting
//
// Purpose: 测试限流器的实际限流功能
//
// Parameters:
//   - t (*testing.T): 测试对象
//
// Returns: 无
func TestRateLimiting(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(2), 1) // 每秒2个请求，突发1个
	host_limiter := limiter.GetLimiter(TEST_HOST)

	// 第一个请求应通过
	assert.True(t, host_limiter.Allow(), "📋 第一个请求应通过")

	// 连续的请求应被限流
	assert.False(t, host_limiter.Allow(), "📋 超出突发容量的请求应被限流")

	// 等待补充令牌
	time.Sleep(500 * time.Millisecond)
	assert.True(t, host_limiter.Allow(), "📋 等待后新的请求应通过")
}

// Function: TestCleanupExpiredLimiters
//
// Purpose: 测试清理过期限流器功能
//
// Parameters:
//   - t (*testing.T): 测试对象
//
// Returns: 无
func TestCleanupExpiredLimiters(t *testing.T) {
	// 使用较短的过期时间
	const quick_expiry = 100 * time.Millisecond

	// 创建一个自定义的IPRateLimiter用于测试
	limiter := &IPRateLimiter{
		hosts:      make(map[string]*hostLimiter),
		mu:         &sync.RWMutex{},
		rate_limit: rate.Limit(TEST_RATE_LIMIT),
		burst_size: TEST_BURST_SIZE,
	}

	// 添加一个过期的限流器
	limiter.mu.Lock()
	limiter.hosts[TEST_HOST] = &hostLimiter{
		limiter:     rate.NewLimiter(rate.Limit(TEST_RATE_LIMIT), TEST_BURST_SIZE),
		last_access: time.Now().Add(-quick_expiry * 2), // 设置为已过期
	}
	limiter.mu.Unlock()

	// 手动触发清理
	limiter.mu.Lock()
	now := time.Now()
	for host, info := range limiter.hosts {
		if now.Sub(info.last_access) > quick_expiry {
			delete(limiter.hosts, host)
		}
	}
	limiter.mu.Unlock()

	// 验证限流器已被清理
	limiter.mu.RLock()
	count := len(limiter.hosts)
	limiter.mu.RUnlock()
	assert.Equal(t, 0, count, "📋 过期的限流器应被清理")
}

// Function: TestRateLimitMiddleware
//
// Purpose: 测试限流中间件功能
//
// Parameters:
//   - t (*testing.T): 测试对象
//
// Returns: 无
func TestRateLimitMiddleware(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 测试缺少domain参数的情况
	t.Run("MissingDomain", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request, _ = http.NewRequest("GET", "/test", nil)

		middleware := RateLimitMiddleware()
		middleware(c)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "📋 缺少domain参数应返回400")
	})

	// 测试正常请求
	t.Run("ValidRequest", func(t *testing.T) {
		// 使用高限制，确保不被限流
		SetRateLimit(100, 100)

		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request, _ = http.NewRequest("GET", "/test?domain="+TEST_HOST, nil)

		middleware := RateLimitMiddleware()
		middleware(c)

		assert.Equal(t, http.StatusOK, recorder.Code, "📋 有效请求应正常通过")
	})

	// 测试被限流的请求
	t.Run("RateLimited", func(t *testing.T) {
		// 使用低限制，确保被限流
		SetRateLimit(1, 1)

		// 第一个请求消耗掉唯一的令牌
		limiter := RateLimiter.GetLimiter(TEST_HOST)
		limiter.Allow()

		// 下一个请求应被限流
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request, _ = http.NewRequest("GET", "/test?domain="+TEST_HOST, nil)

		middleware := RateLimitMiddleware()
		middleware(c)

		assert.Equal(t, http.StatusTooManyRequests, recorder.Code, "📋 超出频率的请求应被限流")
	})

	// 恢复默认值
	SetRateLimit(DEFAULT_RATE_LIMIT, DEFAULT_BURST_SIZE)
}

// Function: TestSetRateLimit
//
// Purpose: 测试动态设置限流参数功能
//
// Parameters:
//   - t (*testing.T): 测试对象
//
// Returns: 无
func TestSetRateLimit(t *testing.T) {
	// 保存原始限流器以便恢复
	original_limiter := RateLimiter

	// 测试设置新的限流参数
	new_rate := 15
	new_burst := 25
	SetRateLimit(new_rate, new_burst)

	assert.Equal(t, rate.Limit(new_rate), RateLimiter.rate_limit, "📋 应更新限流率")
	assert.Equal(t, new_burst, RateLimiter.burst_size, "📋 应更新突发大小")

	// 恢复原始限流器
	RateLimiter = original_limiter
}
````

## File: internal/api/rate.go
````go
// Author       :loyd
// Date         :2025-03-30 11:24:58
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:28:50
// Description  :基于IP的请求限流中间件
//

package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 限流器配置常量
const (
	DEFAULT_RATE_LIMIT = 10 // 每秒请求数限制
	DEFAULT_BURST_SIZE = 20 // 默认令牌桶容量
	CLEANUP_INTERVAL   = 10 // 清理间隔（分钟）
	LIMITER_EXPIRY     = 60 // 限流器过期时间（分钟）
)

// IPRateLimiter 结构体定义
type IPRateLimiter struct {
	hosts      map[string]*hostLimiter // 主机对应的限流器映射
	mu         *sync.RWMutex           // 读写锁，保护并发访问
	rate_limit rate.Limit              // 每秒可以生成的令牌数
	burst_size int                     // 令牌桶的容量
}

// 主机限流器，包含过期时间
type hostLimiter struct {
	limiter     *rate.Limiter // 限流器实例
	last_access time.Time     // 最后访问时间
}

// NewIPRateLimiter 创建一个新的IP限流器
//
// 功能: 初始化并返回一个新的IPRateLimiter实例
//
// 参数:
//   - rate_limit (rate.Limit): 每秒允许的请求数
//   - burst_size (int): 令牌桶大小
//
// 返回值: *IPRateLimiter - 初始化后的限流器实例
func NewIPRateLimiter(rate_limit rate.Limit, burst_size int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		hosts:      make(map[string]*hostLimiter),
		mu:         &sync.RWMutex{},
		rate_limit: rate_limit,
		burst_size: burst_size,
	}

	// 启动定期清理过期限流器的协程
	go limiter.cleanupExpiredLimiters()

	return limiter
}

// 全局默认限流器实例
var RateLimiter = NewIPRateLimiter(DEFAULT_RATE_LIMIT, DEFAULT_BURST_SIZE)

// GetLimiter 获取特定主机的限流器
//
// 功能: 返回指定主机的限流器，若不存在则创建新的
//
// 参数:
//   - host (string): 主机标识符，通常是域名或IP地址
//
// 返回值: *rate.Limiter - 对应主机的限流器实例
func (r *IPRateLimiter) GetLimiter(host string) *rate.Limiter {
	r.mu.RLock()
	limiterInfo, exists := r.hosts[host]
	r.mu.RUnlock()

	if !exists {
		r.mu.Lock()
		// 双重检查，防止并发问题
		limiterInfo, exists = r.hosts[host]
		if !exists {
			limiterInfo = &hostLimiter{
				limiter:     rate.NewLimiter(r.rate_limit, r.burst_size),
				last_access: time.Now(),
			}
			r.hosts[host] = limiterInfo
		}
		r.mu.Unlock()
	} else {
		// 更新最后访问时间
		r.mu.Lock()
		limiterInfo.last_access = time.Now()
		r.mu.Unlock()
	}

	return limiterInfo.limiter
}

// cleanupExpiredLimiters 清理过期的限流器
//
// 功能: 定期检查并移除长时间未使用的限流器，防止内存泄漏
//
// 参数: 无
//
// 返回值: 无
func (r *IPRateLimiter) cleanupExpiredLimiters() {
	ticker := time.NewTicker(CLEANUP_INTERVAL * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		for host, info := range r.hosts {
			// 删除超过过期时间未使用的限流器
			if now.Sub(info.last_access) > LIMITER_EXPIRY*time.Minute {
				delete(r.hosts, host)
			}
		}
		r.mu.Unlock()
	}
}

/**
 * Function: RateLimitMiddleware
 *
 * Purpose: 为Gin框架提供基于域名参数的请求限流功能，仅对特定路由生效
 *
 * Parameters:
 *
 * Returns: gin.HandlerFunc - Gin中间件处理函数
 */
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求路径
		path := c.Request.URL.Path

		// 仅对 /proxy 路由进行域名检查和限流
		if path == "/proxy" {
			host := c.Request.URL.Query().Get("domain")
			if host == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "❌ 缺少必要的domain参数",
				})
				c.Abort()
				return
			}

			limiter := RateLimiter.GetLimiter(host)
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code": 429, // 标准的HTTP状态码
					"msg":  "⚠️ 请求频率超限，请稍后再试",
					"host": host,
				})
				c.Abort()
				return
			}
		}

		// 对所有路由继续处理请求
		c.Next()
	}
}

// SetRateLimit 设置全局限流参数
//
// 功能: 动态调整全局默认限流器的参数
//
// 参数:
//   - new_limit (int): 新的每秒请求数限制
//   - new_burst (int): 新的令牌桶容量
//
// 返回值: 无
func SetRateLimit(new_limit int, new_burst int) {
	RateLimiter = NewIPRateLimiter(rate.Limit(new_limit), new_burst)
}
````

## File: internal/config/config_test.go
````go
// Author       :loyd
// Date         :2025-03-30 12:00:00
// LastEditors  :loyd
// LastEditTime :2025-03-30 12:00:00
// Description  :配置模块测试

package config

import (
	"os"
	"path/filepath"
	"testing"
)

/**
 * Function: TestLoadConfig_ValidFile
 *
 * Purpose: 测试加载有效配置文件的功能
 *
 * Parameters:
 *   - t (*testing.T): 测试对象
 *
 * Returns: 无
 */
func TestLoadConfig_ValidFile(t *testing.T) {
	// 创建临时目录
	temp_dir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("❌ 无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(temp_dir)

	// 创建有效的配置文件
	config_content := `
[server]
port = 8080

[pool]
size = 100

[proxy]
url = "http://example.com/proxy"

[ip]
lifetime_seconds = 3600
delay = [10, 20, 30]
`
	config_path := filepath.Join(temp_dir, "config.toml")
	err = os.WriteFile(config_path, []byte(config_content), 0644)
	if err != nil {
		t.Fatalf("❌ 无法写入测试配置文件: %v", err)
	}

	// 测试加载配置
	config, err := LoadConfig(config_path)
	if err != nil {
		t.Fatalf("❌ 加载有效配置文件应成功，但返回错误: %v", err)
	}

	// 验证配置内容
	if config.Server.Port != 8080 {
		t.Errorf("❌ 服务器端口配置错误，期望 8080，实际为 %d", config.Server.Port)
	}
	if config.Pool.Size != 100 {
		t.Errorf("❌ 代理池大小配置错误，期望 100，实际为 %d", config.Pool.Size)
	}
	if config.ProxySource.Url != "http://example.com/proxy" {
		t.Errorf("❌ 代理URL配置错误，期望 'http://example.com/proxy'，实际为 '%s'",
			config.ProxySource.Url)
	}
	if config.Ip.Lifetime_seconds != 3600 {
		t.Errorf("❌ IP生命周期配置错误，期望 3600，实际为 %d",
			config.Ip.Lifetime_seconds)
	}
	if len(config.Ip.Delay) != 3 || config.Ip.Delay[0] != 10 ||
		config.Ip.Delay[1] != 20 || config.Ip.Delay[2] != 30 {
		t.Errorf("❌ 延迟配置错误，期望 [10,20,30]，实际为 %v", config.Ip.Delay)
	}

	// 验证全局配置已设置
	global_config := GetConfig()
	if global_config != config {
		t.Error("❌ 全局配置未正确设置")
	}
}

/**
 * Function: TestLoadConfig_FileNotFound
 *
 * Purpose: 测试加载不存在的配置文件时的错误处理
 *
 * Parameters:
 *   - t (*testing.T): 测试对象
 *
 * Returns: 无
 */
func TestLoadConfig_FileNotFound(t *testing.T) {
	// 测试加载不存在的文件
	_, err := LoadConfig("/path/to/nonexistent/config.toml")
	if err == nil {
		t.Error("❌ 加载不存在的配置文件应返回错误，但没有")
	}
}

/**
 * Function: TestLoadConfig_InvalidTOML
 *
 * Purpose: 测试加载格式无效的TOML配置文件时的错误处理
 *
 * Parameters:
 *   - t (*testing.T): 测试对象
 *
 * Returns: 无
 */
func TestLoadConfig_InvalidTOML(t *testing.T) {
	// 创建临时目录
	temp_dir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("❌ 无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(temp_dir)

	// 创建无效的TOML配置文件
	invalid_content := `
[server]
port = "not an integer" # 应该是整数，不是字符串
`
	invalid_path := filepath.Join(temp_dir, "invalid.toml")
	err = os.WriteFile(invalid_path, []byte(invalid_content), 0644)
	if err != nil {
		t.Fatalf("❌ 无法写入测试配置文件: %v", err)
	}

	// 测试加载无效配置
	_, err = LoadConfig(invalid_path)
	if err == nil {
		t.Error("❌ 加载格式无效的配置文件应返回错误，但没有")
	}
}
````

## File: internal/config/config.go
````go
// Author       :loyd
// Date         :2025-03-29 01:46:00
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:37:27
// Description  :读取配置文件
//

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

// 定义应用程序配置结构
type Config struct {
	Server      ServerConfig      `toml:"server"`
	Pool        PoolConfig        `toml:"pool"`
	ProxySource ProxySourceConfig `toml:"proxy_source"`
	Ip          IpConfig          `toml:"ip"`
}

// 服务器相关配置
type ServerConfig struct {
	Port int `toml:"port"`
}

// 代理池相关配置
type PoolConfig struct {
	Size int `toml:"size"`
}

// 代理源相关配置
type ProxySourceConfig struct {
	Url              []string `toml:"urls"`
	RequestFrequency int      `toml:"request_frequency"`
	LocalPath        []string `toml:"local_path"`
	ReadFrequency    int      `toml:"read_frequency"`
}

// IP相关配置
type IpConfig struct {
	Lifetime_seconds int   `toml:"lifetime_seconds"`
	Delay            []int `toml:"delay"`
}

// 全局配置实例
var APP_CONFIG *Config

/**
 * Function: LoadConfig
 *
 * Purpose: 从指定路径读取TOML配置文件并解析到Config结构体
 *
 * Parameters:
 *   - config_path (string): 配置文件的路径
 *
 * Returns: (*Config, error) - 返回配置对象指针和可能的错误
 */
func LoadConfig(config_path string) (*Config, error) {
	// 获取绝对路径
	abs_path, err := filepath.Abs(config_path)
	if err != nil {
		return nil, fmt.Errorf("❌ 无法获取配置文件绝对路径: %w", err)
	}

	// 读取配置文件
	fmt.Printf("📖 正在读取配置文件: %s\n", abs_path)
	config_data, err := os.ReadFile(abs_path)
	if err != nil {
		return nil, fmt.Errorf("❌ 读取配置文件失败: %w", err)
	}

	// 解析TOML配置到结构体
	config := &Config{}
	err = toml.Unmarshal(config_data, config)
	if err != nil {
		return nil, fmt.Errorf("❌ 解析TOML配置失败: %w", err)
	}

	// 设置全局配置
	APP_CONFIG = config

	fmt.Println("✅ 配置加载成功!")
	return config, nil
}

/**
 * Function: GetConfig
 *
 * Purpose: 获取全局配置实例
 *
 * Parameters: 无
 *
 * Returns: *Config - 返回全局配置对象指针
 */
func GetConfig() *Config {
	return APP_CONFIG
}
````

## File: internal/log/log.go
````go
// Author       :loyd
// Date         :2025-03-29 02:00:56
// LastEditors  :loyd
// LastEditTime :2025-03-30 19:14:09
// Description  :日志工具封装
//

package log

import (
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 常量定义
const (
	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"
)

// LogConfig 日志配置结构体
type LogConfig struct {
	LogLevel   string // 日志级别，可选值: debug, info, warn, error
	LogFile    string // 日志文件路径，为空则只输出到控制台
	MaxSize    int    // 单个日志文件最大大小(MB)，超过后会轮转
	MaxBackups int    // 最大备份日志文件数量
	MaxAge     int    // 日志文件保留天数
	Compress   bool   // 是否压缩旧日志文件
}

// 定义私有的全局logger实例
var log = logrus.New()

// 初始化日志
func init() {
	// 默认配置 - 在SetupLogger被显式调用前使用
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.SetReportCaller(true)
}

// SetupLogger 配置日志系统
func SetupLogger(config LogConfig) {
	// 创建日志目录
	if config.LogFile != "" {
		log_dir := path.Dir(config.LogFile)
		if _, err := os.Stat(log_dir); os.IsNotExist(err) {
			err := os.MkdirAll(log_dir, 0755)
			if err != nil {
				log.Fatalf("❌ 创建日志目录失败: %v", err)
			}
		}

		// 配置日志轮转
		lumberjack_logger := &lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    config.MaxSize, // 单位MB
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge, // 天数
			Compress:   config.Compress,
		}

		// 同时输出到控制台和文件
		multi_writer := io.MultiWriter(os.Stdout, lumberjack_logger)
		log.SetOutput(multi_writer)
	}

	// 设置日志格式
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", filename + ":" + strconv.Itoa(f.Line)
		},
	})

	// 设置日志级别
	switch config.LogLevel {
	case LOG_LEVEL_DEBUG:
		log.SetLevel(logrus.DebugLevel)
	case LOG_LEVEL_INFO:
		log.SetLevel(logrus.InfoLevel)
	case LOG_LEVEL_WARN:
		log.SetLevel(logrus.WarnLevel)
	case LOG_LEVEL_ERROR:
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	Info("🚀 日志系统初始化成功")
}

// 对外暴露的日志方法，直接调用

// Debug 输出调试日志
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info 输出信息日志
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn 输出警告日志
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error 输出错误日志
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatal 输出致命错误日志
func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// WithFields 创建带字段的日志条目
func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}
````

## File: internal/models/domaininfo_test.go
````go
package models

import (
	"fmt"
	"testing"
	"time"
)

// Function: TestGetIPCount
//
// Purpose: 测试DomainInfo结构体的GetIPCount方法
//
// Parameters:
//   - t (*testing.T): 测试上下文对象
//
// Returns:
func TestGetIPCount(t *testing.T) {
	// 测试用例1: 空IP池
	t.Run("EmptyIPPool", func(t *testing.T) {
		// 创建一个空的DomainInfo实例
		domainInfo := NewDomainInfo("test.com")

		// 验证IP数量为0
		if count := domainInfo.GetIPCount(); count != 0 {
			t.Errorf("❌ 期望IP数量为0，但得到 %d", count)
		}
	})

	// 测试用例2: 含有多个IP的池
	t.Run("MultipleIPs", func(t *testing.T) {
		// 创建DomainInfo实例
		domainInfo := NewDomainInfo("example.com")

		// 添加3个IP
		ip1 := &IpDetails{
			ip:            "192.168.1.1",
			last_use_time: time.Now(),
		}
		ip2 := &IpDetails{
			ip:            "192.168.1.2",
			last_use_time: time.Now(),
		}
		ip3 := &IpDetails{
			ip:            "192.168.1.3",
			last_use_time: time.Now(),
		}

		domainInfo.Add(ip1)
		domainInfo.Add(ip2)
		domainInfo.Add(ip3)

		// 验证IP数量为3
		if count := domainInfo.GetIPCount(); count != 3 {
			t.Errorf("❌ 期望IP数量为3，但得到 %d", count)
		}

		// 删除一个IP后再次验证
		domainInfo.Delete(ip2)
		if count := domainInfo.GetIPCount(); count != 2 {
			t.Errorf("❌ 删除一个IP后期望数量为2，但得到 %d", count)
		}
	})

	// 测试用例3: 测试IP重复添加
	t.Run("DuplicateIPs", func(t *testing.T) {
		domainInfo := NewDomainInfo("duplicate.com")

		ip := &IpDetails{
			ip:            "192.168.1.1",
			last_use_time: time.Now(),
		}

		// 多次添加相同IP
		domainInfo.Add(ip)
		domainInfo.Add(ip)

		// OrderedMap会覆盖相同键的值，所以数量应该为1
		if count := domainInfo.GetIPCount(); count != 1 {
			t.Errorf("❌ 添加重复IP后期望数量为1，但得到 %d", count)
		}
	})
}

// Function: TestGet
//
// Purpose: 测试DomainInfo结构体的Get方法
//
// Parameters:
//   - t (*testing.T): 测试上下文对象
//
// Returns:
func TestGet(t *testing.T) {
	// 测试用例1: 空IP池
	t.Run("EmptyIPPool", func(t *testing.T) {
		// 创建一个空的DomainInfo实例
		domain_info := NewDomainInfo("empty.com")

		// 尝试获取IP
		ip, err := domain_info.Get()

		// 验证返回错误且IP为nil
		if ip != nil {
			t.Errorf("❌ 期望IP为nil，但得到 %v", ip)
		}
		if err == nil {
			t.Error("❌ 期望返回错误，但得到nil")
		}
		expected_err_msg := fmt.Sprintf("当前域名%s%s", "empty.com", ERROR_NO_IP_AVAILABLE)
		if err.Error() != expected_err_msg {
			t.Errorf("❌ 期望错误信息为 '%s'，但得到 '%s'", expected_err_msg, err.Error())
		}
	})

	// 测试用例2: 单个IP
	t.Run("SingleIP", func(t *testing.T) {
		// 创建DomainInfo实例
		domain_info := NewDomainInfo("single.com")

		// 添加一个IP
		old_time := time.Now().Add(-1 * time.Hour) // 1小时前
		ip_details := &IpDetails{
			ip:            "192.168.1.1",
			last_use_time: old_time,
		}
		domain_info.Add(ip_details)

		// 记录当前域名的最后使用时间
		old_domain_use_time := domain_info.GetLastUseTime()

		// 获取IP
		result_ip, err := domain_info.Get()

		// 验证获取的IP
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if result_ip == nil {
			t.Error("❌ 期望返回IP，但得到nil")
		}
		if result_ip.ip != "192.168.1.1" {
			t.Errorf("❌ 期望IP为 '192.168.1.1'，但得到 '%s'", result_ip.ip)
		}

		// 验证last_use_time已更新
		if !result_ip.last_use_time.After(old_time) {
			t.Errorf("❌ IP的最后使用时间应该更新")
		}

		// 验证域名的最后使用时间也已更新
		if !domain_info.GetLastUseTime().After(old_domain_use_time) {
			t.Error("❌ 域名的最后使用时间应该更新")
		}
	})

	// 测试用例3: 多个IP并验证轮换
	t.Run("MultipleIPs", func(t *testing.T) {
		// 创建DomainInfo实例
		domain_info := NewDomainInfo("multiple.com")

		// 添加三个IP
		ip1 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now()}
		ip2 := &IpDetails{ip: "192.168.1.2", last_use_time: time.Now()}
		ip3 := &IpDetails{ip: "192.168.1.3", last_use_time: time.Now()}

		domain_info.Add(ip1) // 此时顺序: [ip1]
		domain_info.Add(ip2) // 此时顺序: [ip1, ip2]
		domain_info.Add(ip3) // 此时顺序: [ip1, ip2, ip3]

		// 获取第一个IP
		first_ip, err := domain_info.Get()
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if first_ip.ip != "192.168.1.1" {
			t.Errorf("❌ 期望首个IP为 '192.168.1.1'，但得到 '%s'", first_ip.ip)
		}

		// 获取第二个IP
		second_ip, err := domain_info.Get()
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if second_ip.ip != "192.168.1.2" {
			t.Errorf("❌ 期望第二个IP为 '192.168.1.2'，但得到 '%s'", second_ip.ip)
		}

		// 获取第三个IP
		third_ip, err := domain_info.Get()
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if third_ip.ip != "192.168.1.3" {
			t.Errorf("❌ 期望第三个IP为 '192.168.1.3'，但得到 '%s'", third_ip.ip)
		}

		// 再次获取，应该回到第一个IP（轮换）
		fourth_ip, err := domain_info.Get()
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if fourth_ip.ip != "192.168.1.1" {
			t.Errorf("❌ 期望第四个IP为 '192.168.1.1'，但得到 '%s'", fourth_ip.ip)
		}
	})
}

// Function: TestAdd
//
// Purpose: 测试DomainInfo结构体的Add方法
//
// Parameters:
//   - t (*testing.T): 测试上下文对象
//
// Returns:
func TestAdd(t *testing.T) {
	// 测试用例1: 向空的DomainInfo添加IP
	t.Run("AddToEmpty", func(t *testing.T) {
		// 创建一个空的DomainInfo实例
		domainInfo := NewDomainInfo("test.com")

		// 初始IP数量应为0
		if count := domainInfo.GetIPCount(); count != 0 {
			t.Errorf("❌ 初始IP数量期望为0，但得到 %d", count)
		}

		// 添加一个IP
		ip := &IpDetails{
			ip:            "192.168.1.1",
			last_use_time: time.Now(),
		}
		domainInfo.Add(ip)

		// 验证IP已被添加
		if count := domainInfo.GetIPCount(); count != 1 {
			t.Errorf("❌ 添加IP后期望数量为1，但得到 %d", count)
		}

		// 验证可以获取到添加的IP
		allIPs := domainInfo.GetAllIPs()
		if len(allIPs) != 1 || allIPs[0].ip != "192.168.1.1" {
			t.Errorf("❌ 无法找到添加的IP '192.168.1.1'")
		}
	})

	// 测试用例2: 添加多个不同的IP
	t.Run("AddMultiple", func(t *testing.T) {
		// 创建DomainInfo实例
		domainInfo := NewDomainInfo("multiple.com")

		// 添加三个不同的IP
		ip1 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now()}
		ip2 := &IpDetails{ip: "192.168.1.2", last_use_time: time.Now()}
		ip3 := &IpDetails{ip: "192.168.1.3", last_use_time: time.Now()}

		domainInfo.Add(ip1)
		domainInfo.Add(ip2)
		domainInfo.Add(ip3)

		// 验证IP数量
		if count := domainInfo.GetIPCount(); count != 3 {
			t.Errorf("❌ 添加三个IP后期望数量为3，但得到 %d", count)
		}

		// 获取所有IP并验证是否都被正确添加
		allIPs := domainInfo.GetAllIPs()
		if len(allIPs) != 3 {
			t.Errorf("❌ GetAllIPs期望返回3个IP，但返回了 %d 个", len(allIPs))
		}

		// 创建一个map来检查所有IP是否存在
		ipMap := make(map[string]bool)
		for _, ip := range allIPs {
			ipMap[ip.ip] = true
		}

		// 验证所有添加的IP都在集合中
		if !ipMap["192.168.1.1"] || !ipMap["192.168.1.2"] || !ipMap["192.168.1.3"] {
			t.Errorf("❌ 无法找到全部添加的IP")
		}
	})

	// 测试用例3: 重复添加相同IP
	t.Run("AddDuplicate", func(t *testing.T) {
		// 创建DomainInfo实例
		domainInfo := NewDomainInfo("duplicate.com")

		// 添加一个IP
		ip1 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now()}
		domainInfo.Add(ip1)

		// 验证IP数量为1
		if count := domainInfo.GetIPCount(); count != 1 {
			t.Errorf("❌ 添加一个IP后期望数量为1，但得到 %d", count)
		}

		// 创建新的IpDetails但使用相同的IP地址
		ip2 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now().Add(1 * time.Hour)}
		domainInfo.Add(ip2)

		// 验证IP数量仍为1（因为使用了相同的键）
		if count := domainInfo.GetIPCount(); count != 1 {
			t.Errorf("❌ 添加重复IP后期望数量仍为1，但得到 %d", count)
		}

		// 验证存储的是最新添加的IP对象
		result, _ := domainInfo.Get()
		if result != ip2 {
			t.Errorf("❌ 期望存储的是最新添加的IP对象，但不是")
		}
	})

	// 测试用例4: 测试添加和删除后的IP位置顺序
	t.Run("AddDeleteOrder", func(t *testing.T) {
		// 创建DomainInfo实例
		domainInfo := NewDomainInfo("order.com")

		// 添加三个不同的IP
		ip1 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now()}
		ip2 := &IpDetails{ip: "192.168.1.2", last_use_time: time.Now()}
		ip3 := &IpDetails{ip: "192.168.1.3", last_use_time: time.Now()}

		domainInfo.Add(ip1)
		domainInfo.Add(ip2)
		domainInfo.Add(ip3)

		// 验证IP数量为3
		if count := domainInfo.GetIPCount(); count != 3 {
			t.Errorf("❌ 添加三个IP后期望数量为3，但得到 %d", count)
		}

		// 获取所有IP验证添加顺序
		allIPs := domainInfo.GetAllIPs()
		if len(allIPs) != 3 {
			t.Errorf("❌ 期望IP数量为3，但得到 %d", len(allIPs))
		}

		// 确认IP顺序为 ip1, ip2, ip3
		expectedOrder := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
		for i, ip := range allIPs {
			if ip.ip != expectedOrder[i] {
				t.Errorf("❌ 位置 %d 期望IP为 %s，但得到 %s", i, expectedOrder[i], ip.ip)
			}
		}

		// 删除中间的IP (ip2)
		domainInfo.Delete(ip2)

		// 验证IP数量为2
		if count := domainInfo.GetIPCount(); count != 2 {
			t.Errorf("❌ 删除一个IP后期望数量为2，但得到 %d", count)
		}

		// 获取所有IP验证删除后的顺序
		allIPs = domainInfo.GetAllIPs()

		// 确认IP顺序为 ip1, ip3
		expectedOrder = []string{"192.168.1.1", "192.168.1.3"}
		for i, ip := range allIPs {
			if ip.ip != expectedOrder[i] {
				t.Errorf("❌ 删除中间IP后，位置 %d 期望IP为 %s，但得到 %s", i, expectedOrder[i], ip.ip)
			}
		}

		// 添加新的IP
		ip4 := &IpDetails{ip: "192.168.1.4", last_use_time: time.Now()}
		domainInfo.Add(ip4)

		// 验证IP数量为3
		if count := domainInfo.GetIPCount(); count != 3 {
			t.Errorf("❌ 添加新IP后期望数量为3，但得到 %d", count)
		}

		// 获取所有IP验证添加后的顺序
		allIPs = domainInfo.GetAllIPs()

		// 确认IP顺序为 ip1, ip3, ip4
		expectedOrder = []string{"192.168.1.1", "192.168.1.3", "192.168.1.4"}
		for i, ip := range allIPs {
			if ip.ip != expectedOrder[i] {
				t.Errorf("❌ 添加新IP后，位置 %d 期望IP为 %s，但得到 %s", i, expectedOrder[i], ip.ip)
			}
		}

		// 验证轮换获取IP的顺序
		first, err := domainInfo.Get()
		if err != nil || first.ip != "192.168.1.1" {
			t.Errorf("❌ 期望获取第一个IP为 192.168.1.1，但得到 %v", first)
		}

		second, err := domainInfo.Get()
		if err != nil || second.ip != "192.168.1.3" {
			t.Errorf("❌ 期望获取第二个IP为 192.168.1.3，但得到 %v", second)
		}

		third, err := domainInfo.Get()
		if err != nil || third.ip != "192.168.1.4" {
			t.Errorf("❌ 期望获取第三个IP为 192.168.1.4，但得到 %v", third)
		}
	})

}
````

## File: internal/models/domaininfo.go
````go
package models

import (
	"fmt"
	"proxies-api/internal/log"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
)

// 常量定义
const (
	ERROR_NO_IP_AVAILABLE = "没有可用的IP"
)

// DomainInfo 域名信息结构体
type DomainInfo struct {
	lock           sync.RWMutex                            // 读写锁，用于保护并发访问
	domain         string                                  // 域名
	last_use_time  time.Time                               // 上次使用时间
	delay          []int                                   // 延迟时间
	ip_details_map *maputil.OrderedMap[string, *IpDetails] // IP池，存储与域名相关的IP详情
}

// Function: NewDomainInfo
//
// Purpose: 创建并初始化一个新的域名信息对象
//
// Parameters:
//   - domain (string): 域名字符串
//
// Returns: *DomainInfo - 初始化完成的DomainInfo实例指针
func NewDomainInfo(domain string) *DomainInfo {
	return &DomainInfo{
		domain:         domain,
		last_use_time:  time.Now(),
		ip_details_map: maputil.NewOrderedMap[string, *IpDetails](),
		lock:           sync.RWMutex{},
	}
}

// Function: SetDelay
//
// Purpose: 设置域名的延迟时间
//
// Parameters:
//   - delay (int): 要设置的延迟时间（毫秒）
//
// Returns:
func (d *DomainInfo) SetDelay(delay []int) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 超过 2 个延迟时间则不设置
	if len(delay) > 2 {
		log.Error("❌ 域名 %s 设置延迟失败，延迟时间超过限制", d.domain)
		return false
	}
	// 设置延迟时间
	d.delay = delay
	log.Debug("⏱️ 域名 %s 设置延迟时间: %v", d.domain, delay)
	return true
}

// Function: GetDomain
//
// Purpose: 获取DomainInfo对象存储的域名
//
// Parameters:
//
// Returns: string - 域名字符串
func (d *DomainInfo) GetDomain() string {
	return d.domain
}

// Function: GetLastUseTime
//
// Purpose: 获取域名最后一次被使用的时间
//
// Parameters:
//
// Returns: time.Time - 域名最后使用的时间点
func (d *DomainInfo) GetLastUseTime() time.Time {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.last_use_time
}

// Function: GetIPCount
//
// Purpose: 获取域名下可用IP的数量
//
// Parameters:
//
// Returns: int - 可用IP的总数量
func (d *DomainInfo) GetIPCount() int {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.ip_details_map.Len()
}

// Function: Get
//
// Purpose: 从域名的IP池中获取一个IP并将其移至队列末尾
//
// Parameters:
//
// Returns:
//   - *IpDetails - 获取到的IP详情对象，如果没有可用IP则为nil
//   - error - 可能的错误，如无可用IP时返回错误信息
func (d *DomainInfo) Get() (*IpDetails, error) {
	// 使用读锁进行并发保护
	d.lock.RLock()
	defer d.lock.RUnlock()

	// 检查IP池是否为空
	if d.ip_details_map.Len() == 0 {
		log.Error("❌ 域名 %s 没有可用IP", d.domain)
		return nil, fmt.Errorf("当前域名%s%s", d.domain, ERROR_NO_IP_AVAILABLE)
	}

	// 获取第一个元素
	elem, ok := d.ip_details_map.Front()
	if !ok {
		log.Error("❌ 域名 %s 获取IP失败", d.domain)
		return nil, fmt.Errorf("当前域名%s%s", d.domain, ERROR_NO_IP_AVAILABLE)
	}

	// 移动到末尾（删除并重新添加）
	ip_details := elem.Value
	d.ip_details_map.Delete(elem.Key)
	d.ip_details_map.Set(ip_details.ip, ip_details)

	// 更新使用时间
	d.last_use_time = time.Now()
	ip_details.last_use_time = time.Now()

	log.Debug("🔄 域名 %s 获取IP: %s 成功", d.domain, ip_details.ip)
	return ip_details, nil
}

// Function: Add
//
// Purpose: 向域名的IP池中添加新的IP
//
// Parameters:
//   - ip (*IpDetails): 要添加的IP详情对象指针
//
// Returns:
func (d *DomainInfo) Add(ip *IpDetails) {
	// 使用写锁进行并发保护
	d.lock.Lock()
	defer d.lock.Unlock()

	// 添加IP到有序映射中
	d.ip_details_map.Set(ip.ip, ip)

	log.Debug("🚀 域名 %s 添加IP %s，当前可用IP数量: %d",
		d.domain,
		ip.ip,
		d.ip_details_map.Len())
}

// Function: Delete
//
// Purpose: 从域名的IP池中删除指定IP
//
// Parameters:
//   - ip (*IpDetails): 要删除的IP详情对象指针
//
// Returns:
func (d *DomainInfo) Delete(ip *IpDetails) {
	// 使用写锁进行并发保护
	d.lock.Lock()
	defer d.lock.Unlock()

	// 检查IP是否存在并删除
	if _, exists := d.ip_details_map.Get(ip.ip); exists {
		d.ip_details_map.Delete(ip.ip)
		log.Debug("✅ 域名 %s 删除IP %s，当前可用IP数量: %d",
			d.domain,
			ip.ip,
			d.ip_details_map.Len())
	} else {
		log.Debug("❌ 域名 %s 删除IP %s 失败，IP不存在",
			d.domain,
			ip.ip)
	}
}

// Function: DeleteByIP
//
// Purpose: 通过IP地址字符串从域名的IP池中删除指定IP
//
// Parameters:
//   - ipAddress (string): 要删除的IP地址字符串
//
// Returns: bool - 删除操作是否成功，成功返回true，失败返回false
func (d *DomainInfo) DeleteByIP(ipAddress string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	if _, exists := d.ip_details_map.Get(ipAddress); exists {
		d.ip_details_map.Delete(ipAddress)
		log.Debug("✅ 域名 %s 删除IP %s，当前可用IP数量: %d",
			d.domain,
			ipAddress,
			d.ip_details_map.Len())
		return true
	}

	log.Debug("❌ 域名 %s 删除IP %s 失败，IP不存在", d.domain, ipAddress)
	return false
}

// Function: GetAllIPs
//
// Purpose: 获取域名下所有的IP详情列表
//
// Parameters:
//
// Returns: []*IpDetails - IP详情对象指针的切片，包含所有IP信息
func (d *DomainInfo) GetAllIPs() []*IpDetails {
	d.lock.RLock()
	defer d.lock.RUnlock()

	result := make([]*IpDetails, 0, d.ip_details_map.Len())
	for elem := range d.ip_details_map.Iter() {
		result = append(result, elem.Value)
	}

	return result
}
````

## File: internal/models/pool.go
````go
package models

import (
	"proxies-api/internal/log"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
)

// Pool 结构体定义了代理IP池的核心数据结构
type Pool struct {
	size           int64                                       // 池的大小
	ip_ttl         int64                                       // 每个IP的生命周期 seconds
	origin_ip_map  *maputil.ConcurrentMap[string, *IpDetails]  // 原始IP映射
	domains        *maputil.ConcurrentMap[string, *DomainInfo] // 域名信息映射
	clean_ticker   *time.Ticker                                // 清理IP的定时器
	ticker_stopped bool                                        // 标记定时器是否已停止
	ticker_mutex   sync.Mutex                                  // 保护定时器操作的互斥锁
	stop_chan      chan struct{}                               // 用于停止定时器的通道
}

/**
 * Function: GetOriginIPs
 *
 * Purpose: 获取池中所有原始IP及其详情的映射
 *
 * Parameters:
 *
 * Returns: map[string]*IpDetails - IP地址到IP详情对象的映射，包含池中所有原始IP信息
 */
func (p *Pool) GetOriginIPs() map[string]*IpDetails {
	ips := make(map[string]*IpDetails)
	p.origin_ip_map.Range(func(key string, value *IpDetails) bool {
		ips[key] = value
		return true
	})
	return ips
}

/**
 * Function: CleanTicker
 *
 * Purpose: 启动一个定时器定期清理池中过期的IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: 无
 */
func (p *Pool) CleanTicker() {
	p.ticker_mutex.Lock()
	defer p.ticker_mutex.Unlock()

	// 如果已经有定时器在运行，则返回
	if p.clean_ticker != nil {
		log.Debug("⚠️ 清理定时器已经在运行中")
		return
	}

	// 创建定时器和停止通道
	p.clean_ticker = time.NewTicker(10 * time.Second)
	p.stop_chan = make(chan struct{})
	p.ticker_stopped = false

	// 在后台启动清理任务
	go func() {
		for {
			select {
			case <-p.clean_ticker.C:
				// 执行清理操作
				p.cleanExpiredIPs()
			case <-p.stop_chan:
				// 收到停止信号，结束goroutine
				log.Debug("🛑 IP清理定时器接收到停止信号")
				return
			}
		}
	}()

	log.Debug("🕒 IP清理定时器已启动，每10秒清理一次过期IP")
}

/**
 * Function: cleanExpiredIPs
 *
 * Purpose: 清理所有域名中过期的IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: 无
 */
func (p *Pool) cleanExpiredIPs() {
	// 遍历所有域名
	p.domains.Range(func(key string, value *DomainInfo) bool {
		// 检查域名的IP池是否为空
		if value.GetIPCount() == 0 {
			log.Debug("❌ 域名 %s 的IP池为空，跳过清理", key)
			return true
		}

		// 清理过期IP
		value.ip_details_map.Range(func(ip string, ipDetails *IpDetails) bool {
			if time.Since(ipDetails.last_use_time).Seconds() > float64(p.ip_ttl) {
				value.Delete(ipDetails)
				log.Debug("✅ 清理过期IP %s 成功", ip)
			}
			return true
		})

		return true
	})
}

/**
 * Function: StopCleanTicker
 *
 * Purpose: 停止IP清理定时器
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功停止定时器
 */
func (p *Pool) StopCleanTicker() bool {
	p.ticker_mutex.Lock()
	defer p.ticker_mutex.Unlock()

	// 检查定时器是否已存在且未停止
	if p.clean_ticker != nil && !p.ticker_stopped {
		// 停止定时器
		p.clean_ticker.Stop()
		// 发送停止信号给goroutine
		close(p.stop_chan)
		// 标记为已停止
		p.ticker_stopped = true
		log.Debug("✅ IP清理定时器已成功停止")
		return true
	}

	log.Debug("⚠️ IP清理定时器未运行或已停止")
	return false
}

/**
 * Function: AddIPToPool
 *
 * Purpose: 向池中添加新的IP和端口，并分配到所有域名
 *
 * Parameters:
 *   - ip (string): 要添加的IP地址
 *   - port (int): 要添加的端口号
 *
 * Returns: 无
 */
func (p *Pool) AddIPToPool(ip string, port int) {
	// 添加原始IP到池中
	ipDetails := NewIpDetails(ip, port)
	p.origin_ip_map.Set(ip, ipDetails)
	log.Debug("✅ 添加原始IP %s:%d 成功", ip, port)

	// 添加到全部 domains
	p.domains.Range(func(key string, value *DomainInfo) bool {
		value.Add(ipDetails)
		return true
	})
}

/**
 * Function: DeleteIPFromPool
 *
 * Purpose: 从池中删除指定IP，并从所有域名中移除该IP
 *
 * Parameters:
 *   - ip (string): 要删除的IP地址
 *
 * Returns:
 */
func (p *Pool) DeleteIPFromPool(ip string) {
	// 删除原始IP从池中
	if ipDetails, ok := p.origin_ip_map.Get(ip); ok {
		p.origin_ip_map.Delete(ip)
		log.Debug("✅ 删除原始IP %s 成功", ip)

		// 从所有域名中删除
		p.domains.Range(func(key string, value *DomainInfo) bool {
			value.Delete(ipDetails)
			return true
		})
	} else {
		log.Debug("❌ 原始IP %s 不存在，无法删除", ip)
	}
}

/**
 * Function: GetSize
 *
 * Purpose: 获取IP池的大小
 *
 * Parameters:
 *
 * Returns: int64 - IP池的大小，表示可以容纳的最大域名数量
 */
func (p *Pool) GetSize() int64 {
	return p.size
}

/**
 * Function: GetDomains
 *
 * Purpose: 获取IP池中所有域名及其信息的映射
 *
 * Parameters:
 *
 * Returns: map[string]*DomainInfo - 域名到域名信息对象的映射
 */
func (p *Pool) GetDomains() map[string]*DomainInfo {
	domains := make(map[string]*DomainInfo)
	p.domains.Range(func(key string, value *DomainInfo) bool {
		domains[key] = value
		return true
	})
	return domains
}

/**
 * Function: GetDomainInfo
 *
 * Purpose: 获取指定域名的详细信息
 *
 * Parameters:
 *   - domain (string): 要查询的域名
 *
 * Returns:
 *   - *DomainInfo - 域名信息对象，如不存在则为nil
 *   - bool - 表示域名是否存在于池中
 */
func (p *Pool) GetDomainInfo(domain string) (*DomainInfo, bool) {
	domainInfo, ok := p.domains.Get(domain)
	if !ok {
		log.Debug("❌ 域名 %s 不存在", domain)
		return nil, false
	}
	log.Debug("✅ 获取域名 %s 的信息成功", domain)
	return domainInfo, true
}

/**
 * Function: AddDomain
 *
 * Purpose: 向IP池中添加新的域名
 *
 * Parameters:
 *   - domain (string): 要添加的域名
 *
 * Returns:
 */
func (p *Pool) AddDomain(domain string) {
	// 检查域名是否已存在
	ok := p.domains.Has(domain)
	if ok {
		log.Debug("❌ 域名 %s 已存在，无法重复添加", domain)
		return
	}

	// 创建新的域名信息对象并添加到池中
	domainInfo := NewDomainInfo(domain)
	p.domains.Set(domain, domainInfo)
	log.Debug("✅ 添加域名 %s 成功", domain)
}

/**
 * Function: DeleteDomain
 *
 * Purpose: 从IP池中删除指定域名及其所有IP信息
 *
 * Parameters:
 *   - domain (string): 要删除的域名
 *
 * Returns:
 */
func (p *Pool) DeleteDomain(domain string) {
	// 检查域名是否存在
	ok := p.domains.Has(domain)
	if !ok {
		log.Debug("❌ 域名 %s 不存在，无法删除", domain)
		return
	}

	// 删除域名信息对象
	p.domains.Delete(domain)
	log.Debug("✅ 删除域名 %s 成功", domain)
}

/**
 * Function: NewPool
 *
 * Purpose: 创建并初始化一个新的池对象
 *
 * Parameters:
 *   - size (int64): 池的大小
 *   - lifetime_per_ip (int64): 每个IP的生命周期(秒)
 *   - urlResource (string): 资源地址字符串
 *
 * Returns: *Pool - 初始化完成的Pool实例指针
 */
func NewPool(size int64, lifetime_per_ip int64) *Pool {
	return &Pool{
		size:          size,
		ip_ttl:        lifetime_per_ip,
		domains:       maputil.NewConcurrentMap[string, *DomainInfo](int(size)),
		origin_ip_map: maputil.NewConcurrentMap[string, *IpDetails](100),
	}
}
````

## File: internal/pool/init.go
````go
// Author       :loyd
// Date         :2025-03-29 13:35:01
// LastEditors  :loyd
// LastEditTime :2025-03-30 11:59:44
// Description  :内存池启动
//

package pool

import (
	"fmt"
	"os"
	"proxies-api/internal/config"
	"proxies-api/internal/models"
	"sync"
)

// 全局内存池实例
var pool_instance *models.Pool

// 默认配置文件路径列表
var DEFAULT_CONFIG_PATH = []string{
	"../config/proxy_pool_config.toml",
	"../../config/proxy_pool_config.toml",
	"./config/proxy_pool_config.toml",
}
var once sync.Once

// Function: InitializeConfigAndPool
//
// Purpose: 初始化配置文件并创建内存池实例，确保只初始化一次
//
// Parameters:
//   - 无
//
// Returns: void - 无返回值
func InitializeConfigAndPool() *models.Pool {
	// 确保只初始化一次
	once.Do(func() {
		// 配置文件
		found := false
		for _, path := range DEFAULT_CONFIG_PATH {
			if _, err := os.Stat(path); err == nil {
				config.LoadConfig(path)
				found = true
				fmt.Println("🔍 找到配置文件:", path)
				break
			}
		}

		// 如果没有找到配置文件，报错并退出
		if !found {
			fmt.Println("❌ 错误: 未找到配置文件，请确保以下路径至少有一个有效的配置文件")
			for _, path := range DEFAULT_CONFIG_PATH {
				fmt.Println("  - " + path)
			}
			os.Exit(1)
		}

		// 初始化内存池
		pool_instance = models.NewPool(
			int64(config.APP_CONFIG.Pool.Size),
			int64(config.APP_CONFIG.Ip.Lifetime_seconds))

		fmt.Println("✅ 内存池初始化成功，容量:", config.APP_CONFIG.Pool.Size)

		// 启动定时器清理过期IP
		pool_instance.CleanTicker()
		fmt.Println("⏰ IP清理定时器已启动")
	})

	return pool_instance
}

// Function: AccessPool
//
// Purpose: 获取内存池实例，如果未初始化则先初始化
//
// Parameters:
//   - 无
//
// Returns: *models.Pool - 返回内存池实例指针
func AccessPool() *models.Pool {
	InitializeConfigAndPool()
	return pool_instance
}
````

## File: .gitignore
````
# ---------- 通用配置 ----------

# 日志

# 测试文件
**/test/

# 临时文件
**/tmp/
**/*.tmp

# 二进制文件
**/*.o
**/*.pyc
**/*.swp
**/*.class
**/*.exe
**/*.dll
**/*.so
**/*.dylib
**/*.a

# 压缩文件
**/*.rar
**/*.tgz
**/*.zip
**/*.7z
**/*.gz
**/*.bz
**/*.xz
**/*.tar
**/*.tar.gz

# 输出目录
**/dist/
**/build/
**/out/
**/target/

# 特殊文件
**/.DS_Store
**/.idea/
**/Thumb.db

# vscode 相关配置
**/.vscode/

# 数据目录
**/data*/
**/archive*/

# ---------- Python 相关 ----------
# 环境
**/.venv/

# python 缓存
**/__pycache__/
**/.pytest_cache/

# ---------- 前端相关 ----------
# node_modules
**/node_modules/

# ---------- 自定义 ----------
**/.vscode/
````

## File: repomix.config.json
````json
{
  "output": {
    "filePath": "repomix-output.md",
    "style": "markdown",
    "parsableStyle": false,
    "fileSummary": true,
    "directoryStructure": true,
    "removeComments": false,
    "removeEmptyLines": false,
    "compress": false,
    "topFilesLength": 5,
    "showLineNumbers": false,
    "copyToClipboard": false,
    "git": {
      "sortByChanges": true,
      "sortByChangesMaxCommits": 100
    }
  },
  "include": [],
  "ignore": {
    "useGitignore": true,
    "useDefaultPatterns": true,
    "customPatterns": []
  },
  "security": {
    "enableSecurityCheck": true
  },
  "tokenCount": {
    "encoding": "o200k_base"
  }
}
````

## File: internal/api/fetch.go
````go
// Author       :loyd
// Date         :2025-03-30 18:58:57
// LastEditors  :loyd
// LastEditTime :2025-03-30 20:04:22
// Description  :从互联网中或者本地文件夹读取
//

package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"proxies-api/internal/config"
	"proxies-api/internal/pool"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/strutil"
)

// IPFetchService 定义IP获取服务结构体
type IPFetchService struct {
	file_ticker   *time.Ticker  // 文件读取定时器
	url_ticker    *time.Ticker  // URL请求定时器
	stop_chan     chan struct{} // 停止信号通道
	ticker_mutex  sync.Mutex    // 互斥锁保护定时器操作
	is_running    bool          // 服务是否在运行
	paths         []string      // 本地文件路径
	urls          []string      // 远程URL列表
	read_interval time.Duration // 文件读取间隔
	req_interval  time.Duration // URL请求间隔
}

// 全局IP获取服务实例
var ip_fetch_service *IPFetchService
var service_once sync.Once

// 重试相关常量
const (
	DEFAULT_MAX_RETRIES    = 3
	DEFAULT_RETRY_INTERVAL = 1 * time.Second
	DEFAULT_MAX_INTERVAL   = 30 * time.Second
	DEFAULT_TIMEOUT        = 30 * time.Second
)

// RetryConfig 重试配置结构体
type RetryConfig struct {
	MaxRetries    int           // 最大重试次数
	RetryInterval time.Duration // 重试间隔
	MaxInterval   time.Duration // 最大重试间隔
	Timeout       time.Duration // 请求超时时间
}

// 默认重试配置
var defaultRetryConfig = RetryConfig{
	MaxRetries:    DEFAULT_MAX_RETRIES,
	RetryInterval: DEFAULT_RETRY_INTERVAL,
	MaxInterval:   DEFAULT_MAX_INTERVAL,
	Timeout:       DEFAULT_TIMEOUT,
}

// 默认时间间隔（秒）
const DEFAULT_READ_INTERVAL = 60
const DEFAULT_REQUEST_INTERVAL = 300

/**
 * Function: retrieveIPFetchService
 *
 * Purpose: 获取单例IPFetchService实例
 *
 * Parameters:
 *   - 无
 *
 * Returns: *IPFetchService - IP获取服务实例
 */
func retrieveIPFetchService() *IPFetchService {
	service_once.Do(func() {
		// 从配置读取参数
		req_frequency := config.GetConfig().ProxySource.RequestFrequency
		read_frequency := config.GetConfig().ProxySource.ReadFrequency
		paths := config.GetConfig().ProxySource.LocalPath
		urls := config.GetConfig().ProxySource.Url

		// 验证时间间隔，确保大于0
		if req_frequency <= 0 {
			fmt.Printf("⚠️ 请求频率配置无效: %d，使用默认值: %d秒\n",
				req_frequency, DEFAULT_REQUEST_INTERVAL)
			req_frequency = DEFAULT_REQUEST_INTERVAL
		}

		if read_frequency <= 0 {
			fmt.Printf("⚠️ 读取频率配置无效: %d，使用默认值: %d秒\n",
				read_frequency, DEFAULT_READ_INTERVAL)
			read_frequency = DEFAULT_READ_INTERVAL
		}

		// 创建IP获取服务实例
		ip_fetch_service = &IPFetchService{
			stop_chan:     make(chan struct{}),
			is_running:    false,
			paths:         paths,
			urls:          urls,
			read_interval: time.Duration(read_frequency) * time.Second,
			req_interval:  time.Duration(req_frequency) * time.Second,
		}
	})

	return ip_fetch_service
}

/**
 * Function: RefreshIPRegistry
 *
 * Purpose: 启动IP获取服务，从文件和URL中周期性获取代理IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: 无
 */
func RefreshIPRegistry() {
	service := retrieveIPFetchService()
	service.Start()
}

/**
 * Function: StopIPFetchService
 *
 * Purpose: 停止IP获取服务
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功停止服务
 */
func StopIPFetchService() bool {
	service := retrieveIPFetchService()
	return service.Stop()
}

/**
 * Function: Start
 *
 * Purpose: 启动IP获取服务，定期从文件和URL获取IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功启动服务
 */
func (s *IPFetchService) Start() bool {
	s.ticker_mutex.Lock()
	defer s.ticker_mutex.Unlock()

	// 如果已经在运行，则返回
	if s.is_running {
		fmt.Println("⚠️ IP获取服务已经在运行中")
		return false
	}

	// 再次检查时间间隔是否大于0
	if s.read_interval <= 0 {
		fmt.Println("⚠️ 读取间隔必须大于0，使用默认间隔60秒")
		s.read_interval = DEFAULT_READ_INTERVAL * time.Second
	}

	if s.req_interval <= 0 {
		fmt.Println("⚠️ 请求间隔必须大于0，使用默认间隔300秒")
		s.req_interval = DEFAULT_REQUEST_INTERVAL * time.Second
	}

	// 确保IP池已初始化
	pool.AccessPool()
	fmt.Println("🚀 启动IP获取服务...")

	// 立即执行一次搜索
	searchFromFiles(s.paths)
	searchFromURLs(s.urls)

	// 创建定时器
	s.file_ticker = time.NewTicker(s.read_interval)
	s.url_ticker = time.NewTicker(s.req_interval)
	s.stop_chan = make(chan struct{})
	s.is_running = true

	// 启动后台协程执行定时任务
	go func() {
		for {
			select {
			case <-s.file_ticker.C:
				fmt.Println("🕒 执行定时文件读取...")
				needUpdate := NeedUpdateIPRegistry()
				if needUpdate {
					searchFromFiles(s.paths)
				}
			case <-s.url_ticker.C:
				fmt.Println("🌐 执行定时URL请求...")
				needUpdate := NeedUpdateIPRegistry()
				if needUpdate {
					searchFromURLs(s.urls)
				}
			case <-s.stop_chan:
				fmt.Println("🛑 IP获取服务接收到停止信号")
				return
			}
		}
	}()

	fmt.Println("✅ IP获取服务已启动")
	return true
}

/**
 * Function: Stop
 *
 * Purpose: 停止IP获取服务
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功停止服务
 */
func (s *IPFetchService) Stop() bool {
	s.ticker_mutex.Lock()
	defer s.ticker_mutex.Unlock()

	// 检查服务是否正在运行
	if !s.is_running {
		fmt.Println("⚠️ IP获取服务未运行")
		return false
	}

	// 停止定时器
	if s.file_ticker != nil {
		s.file_ticker.Stop()
		s.file_ticker = nil
	}

	if s.url_ticker != nil {
		s.url_ticker.Stop()
		s.url_ticker = nil
	}

	// 发送停止信号
	close(s.stop_chan)

	// 更新状态
	s.is_running = false
	fmt.Println("✅ IP获取服务已停止")

	return true
}

/**
 * Function: NeedUpdateIPRegistry
 *
 * Purpose: 检查代理池中的IP数量是否低于阈值，判断是否需要更新IP库
 *
 * Parameters: 无
 *
 * Returns: bool - 如果IP数量低于阈值返回true表示需要更新，否则返回false
 */
func NeedUpdateIPRegistry() bool {

	// 检查池大小与可用IP
	available_ips := pool.AccessPool().GetOriginIPs()
	available_count := len(available_ips)

	// 获取配置中设置的最小IP数量要求
	MIN_IP_COUNT := int(pool.AccessPool().GetSize())

	// 如果当前可用IP数量低于阈值，则需要更新
	if available_count < MIN_IP_COUNT {
		fmt.Printf("⚠️ 当前可用IP数量(%d)低于设定阈值(%d)，需要更新IP库\n",
			available_count, MIN_IP_COUNT)
		return true
	}

	fmt.Printf("✅ 当前可用IP数量充足: %d/%d\n", available_count, MIN_IP_COUNT)
	return false
}

/**
 * Function: searchFromFiles
 *
 * Purpose: 从文件中读取并解析IP地址
 *
 * Parameters:
 *   - paths ([]string): 文件路径列表
 *
 * Returns: 无
 */
func searchFromFiles(paths []string) {
	// 读取配置文件
	file_paths := findValidFilePaths(paths)
	if len(file_paths) == 0 {
		return
	}

	// 如果找到了文件，打印找到的文件路径
	fmt.Println("✅ 找到配置文件:", file_paths)

	// 从文件中加载
	ip_regex := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::(\d{1,5}))?`

	for _, file_path := range file_paths {
		// 从文件读取内容
		file_content, err := os.ReadFile(file_path)
		if err != nil {
			fmt.Println("❌ 读取文件失败:", err)
			continue
		}

		// 解析并处理IP
		parseAndProcessIPs(string(file_content), ip_regex)
	}
}

/**
 * Function: ContainsIgnoreCase
 *
 * Purpose: 检查字符串是否包含子字符串，忽略大小写
 *
 * Parameters:
 *   - source (string): 源字符串
 *   - substr (string): 要检查的子字符串
 *
 * Returns: bool - 如果源字符串包含子字符串（忽略大小写）则返回true，否则返回false
 */
func ContainsIgnoreCase(source string, substr string) bool {
	// 转换为小写后进行比较
	return strings.Contains(
		strings.ToLower(source),
		strings.ToLower(substr),
	)
}

/**
 * Function: searchFromURLs
 *
 * Purpose: 从URL中获取并解析IP地址，支持网页内容解析和错误重试
 *
 * Parameters:
 *   - urls ([]string): URL列表
 *
 * Returns: 无
 */
func searchFromURLs(urls []string) {
	if len(urls) == 0 {
		fmt.Println("⚠️ 没有配置URL来源")
		return
	}

	ip_regex := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::(\d{1,5}))?`

	// 创建自定义HTTP客户端，添加超时设置
	client := &http.Client{
		Timeout: defaultRetryConfig.Timeout,
	}

	for _, url := range urls {
		fmt.Println("🌐 正在从URL获取IP:", url)

		// 创建请求对象
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("❌ 创建请求对象失败:", err)
			continue
		}

		// 添加用户代理头部，降低被拒绝的可能性
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		// 使用重试机制发送请求
		var resp *http.Response
		var retryErr error
		for retry := 0; retry <= defaultRetryConfig.MaxRetries; retry++ {
			if retry > 0 {
				// 计算退避时间
				backoff := defaultRetryConfig.RetryInterval * time.Duration(1<<uint(retry-1))
				if backoff > defaultRetryConfig.MaxInterval {
					backoff = defaultRetryConfig.MaxInterval
				}
				fmt.Printf("🔄 第 %d 次重试，等待 %v 后重试...\n", retry, backoff)
				time.Sleep(backoff)
			}

			// 发送HTTP请求
			resp, retryErr = client.Do(req)
			if retryErr == nil && resp.StatusCode == http.StatusOK {
				break
			}

			if retryErr != nil {
				fmt.Printf("❌ 第 %d 次请求失败: %v\n", retry+1, retryErr)
			} else {
				fmt.Printf("❌ 第 %d 次请求返回非200状态码: %d %s\n", retry+1, resp.StatusCode, resp.Status)
				resp.Body.Close()
			}

			if retry == defaultRetryConfig.MaxRetries {
				fmt.Printf("❌ 达到最大重试次数 %d，放弃请求\n", defaultRetryConfig.MaxRetries)
				continue
			}
		}

		if retryErr != nil {
			fmt.Printf("❌ 所有重试都失败，放弃请求: %v\n", retryErr)
			continue
		}

		defer resp.Body.Close()

		// 读取响应内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("❌ 读取响应失败:", err)
			continue
		}

		content_type := resp.Header.Get("Content-Type")
		fmt.Printf("📄 内容类型: %s\n", content_type)

		// 根据内容类型进行不同处理
		if ContainsIgnoreCase(content_type, "html") {
			fmt.Println("🔍 检测到HTML内容，进行特殊解析...")
			// 解析HTML内容中的IP
			parseHTMLContent(string(body), ip_regex)
		} else {
			// 解析普通文本内容
			parseAndProcessIPs(string(body), ip_regex)
		}
	}
}

/**
 * Function: parseHTMLContent
 *
 * Purpose: 解析HTML内容中的IP地址
 *
 * Parameters:
 *   - html_content (string): HTML内容
 *   - ip_regex (string): 用于匹配IP的正则表达式
 *
 * Returns: 无
 */
func parseHTMLContent(html_content string, ip_regex string) {
	// 移除HTML标签，只保留文本内容
	text_content := strutil.ReplaceWithMap(html_content, map[string]string{
		`<[^>]*>`: " ",
	})

	// 移除多余空白
	text_content = strutil.ReplaceWithMap(text_content, map[string]string{
		`\s+`: " ",
	})

	fmt.Println("🧹 已清理HTML标签，正在解析IP...")

	// 解析清理后的文本中的IP
	parseAndProcessIPs(text_content, ip_regex)
}

/**
 * Function: findValidFilePaths
 *
 * Purpose: 检验并收集所有有效的文件路径
 *
 * Parameters:
 *   - paths ([]string): 待验证的文件路径列表
 *
 * Returns: []string - 所有存在且有效的文件路径列表，如果没有找到则返回空数组
 */
func findValidFilePaths(paths []string) []string {
	var validated_paths []string

	// 遍历配置中的本地路径
	for _, path := range paths {
		// 检查文件是否存在
		if _, err := os.Stat(path); err == nil {
			// 📁 文件存在，添加到有效路径列表
			validated_paths = append(validated_paths, path)
		}
	}

	return validated_paths
}

/**
 * Function: parseAndProcessIPs
 *
 * Purpose: 使用正则表达式解析文本中的IP地址并添加到代理池
 *
 * Parameters:
 *   - content (string): 要解析的文本内容
 *   - regex (string): 用于匹配IP的正则表达式
 *
 * Returns: 无
 */
func parseAndProcessIPs(content string, regex string) {
	// 确保IP池已初始化
	poolInstance := pool.AccessPool()

	// 匹配所有IP和端口
	matches := strutil.RegexMatchAllGroups(regex, content)

	// 输出匹配结果并添加到池中
	for _, match := range matches {
		if len(match) > 2 && match[2] != "" {
			// 有端口的情况
			ip := match[1]
			port_str := match[2]

			// 转换端口为整数
			port, err := strconv.Atoi(port_str)
			if err != nil {
				fmt.Printf("❌ 端口格式错误: %s:%s\n", ip, port_str)
				continue
			}

			fmt.Printf("✅ 找到IP和端口: %s:%d 并添加到代理池\n", ip, port)
			poolInstance.AddIPToPool(ip, port)
		} else if len(match) > 1 {
			// 没有端口的情况，使用默认端口80
			ip := match[1]
			fmt.Printf("✅ 找到IP: %s (使用默认端口:80) 并添加到代理池\n", ip)
			poolInstance.AddIPToPool(ip, 80)
		}
	}
}

/**
 * Function: GetProxyPoolStatus
 *
 * Purpose: 获取代理池的当前状态信息
 *
 * Parameters: 无
 *
 * Returns: map[string]interface{} - 包含池状态的映射，如池大小、原始IP数量等
 */
func GetProxyPoolStatus() map[string]interface{} {
	poolInstance := pool.AccessPool()

	// 获取原始IP数量
	originIPs := poolInstance.GetOriginIPs()
	originIPCount := len(originIPs)

	// 获取域名数量
	domains := poolInstance.GetDomains()
	domainCount := len(domains)

	// 构建状态信息
	status := map[string]interface{}{
		"pool_size":       poolInstance.GetSize(),
		"origin_ip_count": originIPCount,
		"domain_count":    domainCount,
		"domains":         domains,
	}

	return status
}
````

## File: internal/models/ipdetails.go
````go
package models

import (
	"fmt"
	"net/http"
	"net/url"
	"proxies-api/internal/log"
	"time"

	"github.com/duke-git/lancet/v2/netutil"
)

// 重试相关常量
const (
	IP_CHECK_MAX_RETRIES    = 2
	IP_CHECK_RETRY_INTERVAL = 1 * time.Second
	IP_CHECK_MAX_INTERVAL   = 5 * time.Second
	IP_CHECK_TIMEOUT        = 5 * time.Second
	IP_CHECK_SUCCESS_CODES  = 200
)

// IPCheckConfig IP检查配置结构体
type IPCheckConfig struct {
	MaxRetries    int           // 最大重试次数
	RetryInterval time.Duration // 重试间隔
	MaxInterval   time.Duration // 最大重试间隔
	Timeout       time.Duration // 请求超时时间
	SuccessCodes  int           // 成功的状态码
	TestURL       string        // 测试URL
}

// 默认IP检查配置
var defaultIPCheckConfig = IPCheckConfig{
	MaxRetries:    IP_CHECK_MAX_RETRIES,
	RetryInterval: IP_CHECK_RETRY_INTERVAL,
	MaxInterval:   IP_CHECK_MAX_INTERVAL,
	Timeout:       IP_CHECK_TIMEOUT,
	SuccessCodes:  IP_CHECK_SUCCESS_CODES,
	TestURL:       "https://www.baidu.com",
}

// IpDetails IP详情结构体
type IpDetails struct {
	ip            string    // IP 地址
	port          int       // 端口号
	acquire_time  time.Time // 获取时间
	last_use_time time.Time // 上次使用时间
}

// Function: GetIP
//
// Purpose: 获取IP详情对象中存储的IP地址
//
// Parameters:
//
// Returns: string - IP地址字符串
func (i *IpDetails) GetIP() string {
	return i.ip
}

func (i *IpDetails) CheckConnection() bool {
	// 创建请求配置
	request := &netutil.HttpRequest{
		Method: "GET",
		RawURL: defaultIPCheckConfig.TestURL,
		Headers: http.Header{
			"User-Agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"},
			"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
			"Accept-Language": []string{"zh-CN,zh;q=0.9,en;q=0.8"},
			"Cache-Control":   []string{"no-cache"},
		},
	}

	// 创建客户端配置
	clientCfg := netutil.HttpClientConfig{
		Timeout: defaultIPCheckConfig.Timeout,
		Proxy: &url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s:%d", i.ip, i.port),
		},
	}

	// 创建HTTP客户端
	client := netutil.NewHttpClientWithConfig(&clientCfg)

	// 重试逻辑
	var lastErr error
	for retry := 0; retry <= defaultIPCheckConfig.MaxRetries; retry++ {
		if retry > 0 {
			// 计算退避时间
			backoff := defaultIPCheckConfig.RetryInterval * time.Duration(1<<uint(retry-1))
			if backoff > defaultIPCheckConfig.MaxInterval {
				backoff = defaultIPCheckConfig.MaxInterval
			}
			log.Debug("🔄 IP %s:%d 第 %d 次重试，等待 %v 后重试...", i.ip, i.port, retry, backoff)
			time.Sleep(backoff)
		}

		// 发送请求
		resp, err := client.SendRequest(request)
		if err != nil {
			lastErr = err
			log.Debug("❌ IP %s:%d 第 %d 次连接失败: %v", i.ip, i.port, retry+1, err)
			continue
		}

		// 检查状态码
		if resp.StatusCode != defaultIPCheckConfig.SuccessCodes {
			lastErr = fmt.Errorf("状态码: %d", resp.StatusCode)
			log.Debug("❌ IP %s:%d 第 %d 次返回状态码: %d", i.ip, i.port, retry+1, resp.StatusCode)
			resp.Body.Close()
			continue
		}

		// 成功，关闭响应体并返回
		resp.Body.Close()
		log.Debug("✅ IP %s:%d 连接检查成功", i.ip, i.port)
		return true
	}

	// 所有重试都失败
	log.Error("❌ IP %s:%d 所有重试都失败，最后错误: %v", i.ip, i.port, lastErr)
	return false
}

// Function: GetPort
//
// Purpose: 获取IP详情对象中存储的端口号
//
// Parameters:
//
// Returns: int - 端口号整数值
func (i *IpDetails) GetPort() int {
	return i.port
}

// Function: GetAcquireTime
//
// Purpose: 获取IP被获取的时间点
//
// Parameters:
//
// Returns: time.Time - IP首次获取的时间
func (i *IpDetails) GetAcquireTime() time.Time {
	return i.acquire_time
}

// Function: GetLastUseTime
//
// Purpose: 获取IP最后一次被使用的时间点
//
// Parameters:
//
// Returns: time.Time - IP最后一次使用的时间
func (i *IpDetails) GetLastUseTime() time.Time {
	return i.last_use_time
}

// Function: NewIpDetails
//
// Purpose: 创建并初始化一个新的IP详情对象
//
// Parameters:
//   - ip (string): IP地址字符串
//   - port (int): IP的端口号
//
// Returns: *IpDetails - 初始化完成的IpDetails实例指针
func NewIpDetails(ip string, port int) *IpDetails {
	now := time.Now()
	return &IpDetails{
		ip:            ip,
		port:          port,
		acquire_time:  now,
		last_use_time: now,
	}
}
````

## File: go.mod
````
module proxies-api

go 1.24.1

require (
	github.com/pelletier/go-toml v1.9.5
	github.com/sirupsen/logrus v1.9.2
)

require (
	github.com/duke-git/lancet/v2 v2.3.5
	github.com/stretchr/testify v1.8.3
	golang.org/x/time v0.11.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20221208152030-732eee02a75a // indirect
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
````

## File: README.md
````markdown
# 代理IP池管理系统

## :star: 概述

本项目是由 Golang 编写的IP控制池，同一个 IP 可以交由不同 `domain` 处理。系统提供高性能的代理IP管理，支持自动获取、清理和分配IP，并针对不同域名提供独立的IP轮询机制和延迟设置。

## :rocket: 特性

- 🔄 IP自动轮询分配与负载均衡
- 🕒 可配置的IP生命周期管理
- 🌐 支持多域名独立管理
- 📊 基于域名的请求频率限制
- 📥 同时支持从URL和本地文件获取IP
- 🧹 自动清理过期IP
- 🛡️ 内置健康检查机制

## :hammer_and_wrench: 安装

```bash
# 克隆仓库
git clone https://github.com/zol9527/proxies-api.git
cd proxies-api

# 安装依赖
go mod download

# 编译项目
go build -o proxy-pool

# 运行服务
./proxy-pool
```

## :fire: 使用

1. 保证代理服务商获取的 `ip` 格式: `<xxx.xxx.xxx.xxx:ip_v4>:<xxx:port>`

2. 在本地或服务器上运行本程序

3. 使用程序请求

    ```python
    # 示例：请求 http://localhost:8000/proxy?domain=baidu.com
    import requests
    
    response = requests.get("http://localhost:8000/proxy", params={
        "domain": "baidu.com",
        "delay": "2,5"  # 可选参数，设置延迟范围
    })
    
    if response.status_code == 200:
        data = response.json()
        proxy = f"{data['ip']}:{data['port']}"
        print(f"获取到代理: {proxy}")
    else:
        print(f"错误: {response.json()['msg']}")
    ```

## :gear: API 文档

### 获取代理 IP

**请求**:
- 路径: `/proxy`
- 方法: `GET`
- 参数:
  - `domain` (必填): 请求域名
  - `delay` (可选): 延迟设置，可以是逗号分隔的数组 (如 `2,5`) 或多个参数 (如 `delay=2&delay=5`)

**成功响应** (200):
```json
{
  "ip": "192.168.1.1",
  "port": 8080,
  "delay": [2, 5]
}
```

**错误响应** (400/429/500):
```json
{
  "code": 400,
  "msg": "❌ 缺少必要的domain参数"
}
```

### 获取系统状态

**请求**:
- 路径: `/status`
- 方法: `GET`
- 参数: 无

**成功响应** (200):
```json
{
  "code": 200,
  "msg": "✅ 获取状态成功",
  "status": {
    "domains_count": 3,
    "origins_count": 50,
    "domains": [
      {
        "name": "example.com",
        "ips_count": 50
      },
      {
        "name": "test.com",
        "ips_count": 50
      }
    ]
  }
}
```

### 限流机制

系统内置请求限流机制，防止单个域名请求过多，默认限制为每秒10个请求，突发上限为20个请求。
注意：限流机制仅应用于 `/proxy` 端点，`/status` 端点不受限流影响。

## :robot: 配置

在 `config/proxy_pool_config.toml` 中设置：

```toml
[server]
port = 8000  # HTTP服务端口

# 代理池配置
[pool]
size = 100  # IP池大小

# 代理源配置
[proxy_source]
urls = [  # 代理源URL列表
    "http://example.com/proxy1",
    "http://example.com/proxy2"
]
request_freq = 10  # URL请求频率(秒)
local_path = [  # 本地代理文件路径
    "/var/proxy_pool/proxy.txt",
    "/var/proxy_pool/proxy2.txt"
]
read_freq = 10  # 文件读取频率(秒)

# IP配置
[ip]
lifetime_seconds = 900  # IP生存时间(秒)
delay = [0, 5]  # 默认延迟范围(秒)
```

配置项说明:
- ✅ 每个域名的IP使用间隔
- ✅ 每个IP的存活时间
- ✅ IP池大小
- ✅ 不同域名单独设置 IP 的延迟时间 (通过API参数)

## :wrench: 开发

项目结构:
```
proxies-api/
├── cmd/            # 命令行入口
│   └── main.go     # 主程序入口
├── config/         # 配置文件
├── internal/       # 内部实现
│   ├── api/        # HTTP接口
│   ├── config/     # 配置解析
│   ├── log/        # 日志工具
│   ├── models/     # 数据结构
│   └── pool/       # IP池管理
└── README.md       # 文档
```

### 单元测试

```bash
go test ./...
```

## :warning: 注意事项

- 确保配置文件中的路径有效
- 对于高并发场景，建议适当调整限流参数
- 系统会定期清理过期IP，无需手动操作

## :page_facing_up: 许可证

本项目采用 MIT 许可证
````
