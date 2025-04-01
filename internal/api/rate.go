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
