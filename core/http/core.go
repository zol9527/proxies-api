/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-21 16:28:26
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package http

import (
	"fmt"
	"ip-controller/config"
	"net/http"
	"strconv"
	"ip-controller/core/pool"
	"time"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Domain string `form:"domain"`
	Delay  int64  `form:"delay, default=12"`
}

/* 获取IP */
func core(c *gin.Context) {

	// 获取客户端 Query "domain"
	domain_query := c.Query("domain")
	if domain_query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "domain查询参数不能为空",
		})
		return
	}

	// 查看 MemoryPool 中是否存在这个域名, 不存在则添加
	_, ok := pool.Pool.Pool[domain_query]
	if !ok {
		pool.Pool.Create(domain_query)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "新建域名IP池，稍后会主动向其添加可用IP，请再次请求",
		})
		return
	}

	// 检查 MemoryPool 中的这个域名是否有 IP, 有则获取一个IP，并返回且
	ip := pool.Pool.GetIP(domain_query)
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "当前域名IP池为空",
		})
		return
	}

	// 获取上次使用时间与设置中的使用延迟
	last_use := ip.LastuseTime
	delay := config.Settings.Get("pool.delay").(int64)

	// 保证在每个 domain 中只有一个用户在操作
	pt_domain.Lock.Lock()
	defer pt_domain.Lock.Unlock()

	// 如果上次使用时间与当前时间的差值大于延迟时间, 则返回这个 ip
	if time.Since(last_use).Seconds() > float64(delay) {
		// 更新这个 ip 的使用时间
		ip.LastuseTime = time.Now()
		// 返回 ip
		c.JSON(http.StatusOK, gin.H{
			"http":  ip,
			"https": ip,
		})
		return
	} else {
		// 如果上次使用时间与当前时间的差值小于延迟时间, 则返回错误
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "当前 IP 未到可用时间",
		})
		return
	}
}

func Server() {

	// 获取 port
	port_str := Settings.Get("server.port").(int64)
	port := strconv.FormatInt(port_str, 10)

	r := gin.New()

	// 路由
	r.GET("/get", core)

	// 启动
	r.Run(fmt.Sprintf(":%s", port))
}
