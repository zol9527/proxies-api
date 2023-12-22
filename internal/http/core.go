/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-22 15:29:21
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package http

import (
	"fmt"
	"ip-controller/config"
	"ip-controller/internal/pool"
	"net/http"
	"strconv"
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
	exist := pool.DomainExist(domain_query)
	if !exist {
		pool.CreateDomain(domain_query)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "新建域名IP池，稍后会主动向其添加可用IP，请再次请求",
		})
		return
	}

	// 检查 MemoryPool 中的这个域名是否有 IP, 有则获取一个IP，并返回且
	ip := pool.GetIP(domain_query)
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "当前域名IP池为空",
		})
		return
	}

	// 确定使用延迟
	delay := config.Settings.Get("pool.delay").(int64)
	pool := pool.GetPool()
	if time.Since(pool.Pool[domain_query].Meta.Last_used_time).Seconds() > float64(delay) {

		// 返回
		c.JSON(http.StatusOK, gin.H{
			"http":  ip,
			"https": ip,
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "当前 IP 未到可用时间",
		})
		return
	}
}

func Server() {

	go func() {
		// 获取 port
		port_str := config.Settings.Get("server.port").(int64)
		port := strconv.FormatInt(port_str, 10)

		r := gin.New()

		// 路由
		r.GET("/get", core)

		// 启动
		r.Run(fmt.Sprintf(":%s", port))

	}()
}
