/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:48:07
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-21 15:56:07
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import (
	"ip-controller/config"
	"ip-controller/core/log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/wangluozhe/requests"
)

var (
	total       int = 0
	total_hour  int = 0
	total_today int = 0
	date            = time.Now().Day()
	hour            = time.Now().Hour()
)

// 当前总共 IP 数量
var sm sync.Map

// 传输新 IP
var IPPipeline = make(chan string, 200)

/* 获取 IP */
func AcquireIP() {

	for {

		// 判断是否已达最大 IP 数量
		pool_size := config.Settings.Get("pool.size").(int64)
		current_ip_number, _ := sm.LoadOrStore("current_ip_number", 0)
		if current_ip_number.(int) >= int(pool_size) {
			time.Sleep(2 * time.Second)
			log.Debug("临时 IP 池已满, 休眠 2 秒后继续")
			continue
		}

		// 请求代理网址
		next_link := config.Settings.Get("proxy.url").(string)
		r, err := requests.Get(next_link, nil)
		if err != nil {
			log.Debug("请求代理网址失败: ", err, " 休眠 2 秒后继续")
			time.Sleep(2 * time.Second)
			continue
		}

		// 解析返回的数据, 并验证 IP 是否是 IPv4 格式
		ip := strings.TrimSpace(r.Text)
		log.Debug("服务器返回: ", ip)
		ip_v4 := strings.Split(ip, ":")[0]
		result := net.ParseIP(ip_v4)
		if result == nil {
			log.Debug("返回的数据不是 IP v4 格式, 休眠 2 秒后继续")
			time.Sleep(2 * time.Second)
			continue
		}

		// IP 数量计算
		current_number, _ := sm.LoadOrStore("current_ip_number", 0)
		current_number = current_number.(int) + 1
		sm.Store("current_ip_number", current_number)
		total += 1
		if time.Now().Day() != date {
			date = time.Now().Day()
			total_today = 0
		}
		total_today += 1
		if time.Now().Hour() != hour {
			hour = time.Now().Hour()
			total_hour = 0
		}
		total_hour += 1
		log.Debugf("获取到 IP: %s \n 本小时内累积获取: %d \n 今日累积获取: %d \n 累积数量: %d", ip, total_hour, total_today, total)

		// 添加 IP 进入管道
		IPPipeline <- ip

		// 休眠 1 秒
		time.Sleep(2 * time.Second)
	}

}
