/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-12 11:57:13
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-21 16:03:29
 * :Description  :清理过期 IP
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import (
	"ip-controller/config"
	"time"
)

func Cleaner() {

	timer := time.NewTicker(5 * time.Second)

	// 计数
	handled_count := 0
	has_handled := false

	// 定时运行
	for range timer.C {

		// 清理不用的域名
		for domain_name, domain := range Pool.Pool {

			// 上次使用时间到现在为止以及大于了 2 天
			if domain.Meta.Last_used_time.Add(2 * 24 * time.Hour).Before(time.Now()) {

				// 删除域名
				Pool.Delete(domain_name)
			}

			// 遍历域下 ip
			for ip, ip_meta := range domain.Meta.Memory_pool {

				// 最开始获取的时间到现在已经大于了配置
				if ip_meta.Initial_time.Add(time.Duration(config.Settings.Get("ip.exist_time").(int64)) * time.Second).Before(time.Now()) {

					// 删除 IP
					Pool.DeleteIP(domain_name, ip)
					if !has_handled {
						handled_count++
					}
				}
			}

			has_handled = true
		}

		// 删减
		current_ip_number, _ := sm.LoadOrStore("current_ip_number", 0)
		if current_ip_number.(int) > 0 {
			current_ip_number = current_ip_number.(int) - handled_count
			sm.Store("current_ip_number", current_ip_number)
		}

		// 重置
		handled_count = 0
		has_handled = false
	}
}
