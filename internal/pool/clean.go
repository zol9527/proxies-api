/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-12 11:57:13
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-07 20:39:24
 * :Description  :清理过期 IP
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import (
	"ip-controller/config"
	"ip-controller/internal/log"
	"time"
)

func Cleaner() {

	timer := time.NewTicker(5 * time.Second)

	// 计数
	handled_count := 0

	// 定时运行
	for range timer.C {

		// 清理不用的域名
	ClearLoop:
		for domain_name, domain := range Pool.Pool {

			// 上次使用时间到现在为止以及大于了 2 天
			if domain.LastUsedTime.Add(2 * 24 * time.Hour).Before(time.Now()) {

				// 删除域名
				DeleteDomain(domain_name)

				// 继续
				continue ClearLoop
			}
		}

		// 清楚过期的 ip
		for ip, ip_meta := range BasePool {
			// 最开始获取的时间到现在已经大于了配置
			if ip_meta.InitialTime.Add(time.Duration(config.Settings.Get("ip.exist_time").(int64)) * time.Second).Before(time.Now()) {
				// 删除 IP
				DeleteIP(ip)
				handled_count++
			}
		}

		// 删减
		log.Debugf("清理了 %d 个 IP", handled_count)

		// 重置
		handled_count = 0
	}
}
