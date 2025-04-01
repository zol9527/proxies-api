package models

import (
	"fmt"
	"net/http"
	"net/url"
	"proxies-api/internal/log"
	"time"

	"github.com/duke-git/lancet/v2/netutil"
)

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

	request := &netutil.HttpRequest{
		Method: "GET",
		RawURL: "https://www.baidu.com",
		Headers: http.Header{
			"User-Agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"},
			"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
			"Accept-Language": []string{"zh-CN,zh;q=0.9,en;q=0.8"},
			"Cache-Control":   []string{"no-cache"},
		},
	}

	clientCfg := netutil.HttpClientConfig{
		Timeout: 5 * time.Second,
		Proxy: &url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s:%d", i.ip, i.port),
		},
	}

	client := netutil.NewHttpClientWithConfig(&clientCfg)
	resp, err := client.SendRequest(request)
	if err != nil {
		log.Error("❌ IP %s:%d 连接失败: %v", i.ip, i.port, err)
		return false
	}

	if resp.StatusCode != 200 {
		log.Error("❌ IP %s:%d 返回状态码: %d", i.ip, i.port, resp.StatusCode)
		return false
	}

	defer resp.Body.Close()
	return true
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
