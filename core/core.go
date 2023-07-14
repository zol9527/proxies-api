/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :loyd
 * :LastEditTime :2023-07-15 01:59:16
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"fmt"
	"ip-controller/pool"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
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

	// 查看是否存在这个二级代理
	_, ok := pool.MemorySecondaryPool.Load(domain_query)
	// 不存在则添加
	if !ok {
		var tmp_secondary_pool sync.Map
		tmp_secondary_pool.Store("domain", domain_query)
		tmp_secondary_pool.Store("last_use", time.Now())
		tmp_secondary_pool.Store("ip_pool", &sync.Map{})
		pool.MemorySecondaryPool.Store(domain_query, &tmp_secondary_pool)
	}

	// 检查二级池中是否存在 ip
	secondary_pool, _ := pool.MemorySecondaryPool.Load(domain_query)
	pt_secondary_pool := secondary_pool.(*sync.Map)

	// 获取池
	ip_pool, _ := pt_secondary_pool.Load("ip_pool")
	pt_ip_pool := ip_pool.(*sync.Map)

	// 检查这个池大小
	size := 0
	pt_ip_pool.Range(func(key, value interface{}) bool {
		size++
		return true
	})

	// 如果没有 ip 则拒绝该请求
	if size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "当前 IP 池为空",
		})
		return
	}

	// 将所有 ip 放入一个切片中
	var ip_slice []string
	pt_ip_pool.Range(func(key, value interface{}) bool {
		ip_slice = append(ip_slice, key.(string))
		return true
	})

	// 随机获取一个 ip
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(size)
	ip := ip_slice[index]

	// 获取这个 ip 的信息
	ip_info, _ := pt_ip_pool.Load(ip)

	// 获取上次使用时间与设置中的使用延迟
	last_use := ip_info.(map[string]interface{})["last_use"].(time.Time)
	delay := Settings.Get("pool.delay").(int64)

	// 如果上次使用时间与当前时间的差值大于延迟时间, 则返回这个 ip
	if time.Since(last_use).Seconds() > float64(delay) {
		// 更新这个 ip 的使用时间
		ip_info.(map[string]interface{})["last_use"] = time.Now()
		pt_ip_pool.Store(ip, ip_info)
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
