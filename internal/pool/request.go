/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:48:07
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-24 11:21:23
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

		// 获取设置
		maxPoolSize := config.Settings.Get("pool.size").(int64)
		proxyUrl := config.Settings.Get("proxy.url").(string)

		// 判断 ip 池是否已经满了
		if len(BaseIPList) >= int(maxPoolSize) {
			log.Debug("IP 池已满, 休眠 2 秒后继续")
			time.Sleep(2 * time.Second)
			continue
		}

		// 请求代理网址
		r, err := requests.Get(proxyUrl, nil)
		if err != nil {
			if r != nil {
				closeErr := r.Body.Close()
				if closeErr != nil {
					log.Debug("关闭 Body 失败: ", closeErr)
					continue
				}
				log.Debug("请求代理网址失败: ", r.Text, " 休眠 2 秒后继续")
				continue
			}
			log.Debug("请求代理网址失败: ", err, " 休眠 2 秒后继续")
			time.Sleep(2 * time.Second)
			continue
		}

		// 存储变量准备关闭 Body
		respText := r.Text
		closeErr := r.Body.Close()
		if closeErr != nil {
			log.Debug("关闭 Body 失败: ", closeErr)
		}

		// 服务器实际返回
		println(respText)

		// 验证 IP 是否是 IPv4 格式
		ip := strings.TrimSpace(respText)
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
		log.Debugf("获取到 IP: %s \n 本小时内累积获取: %d \n 今日累积获取: %d \n 累积数量: %d, 代理池中数量: %d", ip, total_hour, total_today, total, len(BaseIPList))

		// 添加 ip 进入全部域名中
		Pool.Add("", ip)
	}
}
