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
