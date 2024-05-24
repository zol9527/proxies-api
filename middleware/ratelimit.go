/*
 * :Date         :2024-05-14 14:51:22
 * :LastEditTime :2024-05-24 16:47:42
 * :Description  :
 *
 */

package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter 结构体
type IPRateLimiter struct {
	hosts map[string]*rate.Limiter
	mu    *sync.RWMutex
	limit rate.Limit // 每秒可以生成的令牌数
	burst int        // 令牌桶的容量
}

var RateLimiter = &IPRateLimiter{
	hosts: make(map[string]*rate.Limiter),
	mu:    &sync.RWMutex{},
	limit: rate.Limit(10), // 每秒10个请求
	burst: 20,             // 桶的大小
}

// 获取IP的限流器
func (r *IPRateLimiter) GetLimiter(host string) *rate.Limiter {
	r.mu.RLock()
	limiter, exists := r.hosts[host]
	r.mu.RUnlock()
	if !exists {
		r.mu.Lock()
		// 双重检查，防止并发问题
		limiter, exists = r.hosts[host]
		if !exists {
			limiter = rate.NewLimiter(r.limit, r.burst)
			r.hosts[host] = limiter
		}
		r.mu.Unlock()
	}
	return limiter
}

// Gin中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.URL.Query().Get("domain")
		if host == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "缺少参数",
			})
			c.Abort()
			return
		}
		limiter := RateLimiter.GetLimiter(host)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 424,
				"msg":  "访问太频繁，请稍后再试",
				"host": host,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
