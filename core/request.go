/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:48:07
 * :LastEditors  :loyd
 * :LastEditTime :2023-07-16 02:08:38
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"ip-controller/log"
	"ip-controller/pool"
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

/* 获取 IP  */
func RequestIP() {

	for {
		// 获取 IP 池中 IP 数量
		ip_pool_size := getPoolSize()

		// 保证 IP 池的数量不超过最大值
		setting_pool_size := Settings.Get("pool.size").(int64)
		if ip_pool_size >= setting_pool_size {
			time.Sleep(2 * time.Second)
			continue
		}

		// 请求代理网址
		r, err := requests.Get(Settings.Get("proxy.url").(string), nil)
		if err != nil {
			log.Log.Error("请求代理网址失败: ", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// 解析返回的数据
		ip := strings.TrimSpace(r.Text)
		log.Log.Debug("服务器返回: ", ip)

		// 验证 IP 是否合法
		if !isIPv4(ip) {
			time.Sleep(2 * time.Second)
			continue
		}

		// 计数
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
		log.Log.Debugf("获取到 IP: %s, 本小时内累积获取: %d, 今日累积获取: %d, 累积数量: %d", ip, total_hour, total_today, total)

		// 添加 IP 到 IP 池
		now := time.Now()
		pool.MemoryPool.Store(ip, now)

		// 休眠 1 秒
		time.Sleep(1 * time.Second)
	}

}
