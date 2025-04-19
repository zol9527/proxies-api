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
