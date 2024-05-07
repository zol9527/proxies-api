/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:48:07
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-07 20:36:02
 * :Description  :请求代理ip 并添加到管道
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import (
	"ip-controller/config"
	"ip-controller/internal/log"
	"net"
	"strings"
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

/* 获取 IP */
func AcquireIP() {

	for {

		// 判断是否已达最大 IP 数量
		setting_pool_size := config.Settings.Get("pool.size").(int64)

		if len(BasePool) > int(setting_pool_size) {
			time.Sleep(2 * time.Second)
			log.Debug("IP 池已满, 休眠 2 秒后继续")
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
		log.Debugf("获取到 IP: %s \n 本小时内累积获取: %d \n 今日累积获取: %d \n 累积数量: %d, 代理池中数量: %d", ip, total_hour, total_today, total, len(BasePool))

		// 添加 ip 进入全部域名中
		AddIP(ip)
	}

}
