/*
 * :Date         :2024-05-14 14:51:22
 * :LastEditTime :2024-05-14 15:31:45
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
	ips   map[string]*rate.Limiter
	mu    *sync.Mutex
	limit rate.Limit // 每秒可以生成的令牌数
	burst int        // 令牌桶的容量
}

var RateLimiter = &IPRateLimiter{
	ips:   make(map[string]*rate.Limiter),
	mu:    &sync.Mutex{},
	limit: rate.Limit(10), // 每秒5个请求
	burst: 20,             // 桶的大小
}

// 添加IP到限流器
func (r *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()
	limiter := rate.NewLimiter(r.limit, r.burst)
	r.ips[ip] = limiter
	return limiter
}

// 获取IP的限流器
func (r *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	r.mu.Lock()
	limiter, exists := r.ips[ip]
	r.mu.Unlock()
	if !exists {
		return r.AddIP(ip)
	}
	return limiter
}

// Gin中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.RemoteIP()
		limiter := RateLimiter.GetLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 424,
				"msg":  "访问太频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
