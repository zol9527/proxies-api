/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-15 22:26:47
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package http

import (
	"fmt"
	"ip-controller/config"
	"ip-controller/internal/pool"
	"ip-controller/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Domain string `form:"domain"`
	Delay  int64  `form:"delay, default=12"`
}

/* 获取IP */
func core(c *gin.Context) {

	// 获取变量
	MemoryPool := pool.GetMemoryPool()

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
	exist := MemoryPool.Exist(domain_query, "")
	if !exist {
		MemoryPool.Create(domain_query)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "新建域名IP池，稍后会主动向其添加可用IP，请再次请求",
		})
		return
	}

	// 检查 MemoryPool 中的这个域名是否有 IP, 有则获取一个IP，并返回且
	ip := MemoryPool.Get(domain_query)
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "当前域名IP无可用, 或者 ip 没有到可用时间, 可再次请求",
		})
		return
	}

	// 返回 IP
	c.JSON(http.StatusOK, gin.H{
		"http":  ip,
		"https": ip,
	})
}

func Server() {

	go func() {
		// 获取 port
		port_str := config.Settings.Get("server.port").(int64)
		port := strconv.FormatInt(port_str, 10)

		r := gin.New()
		r.Use(middleware.RateLimitMiddleware())

		// 路由
		r.GET("/get", core)

		// 启动
		r.Run(fmt.Sprintf(":%s", port))

	}()
}
