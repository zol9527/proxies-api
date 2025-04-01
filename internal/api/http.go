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
