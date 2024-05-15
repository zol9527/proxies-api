/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-12 11:57:13
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-15 22:10:25
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

	// 获取设置
	existTime := config.Settings.Get("ip.exist_time").(int64)

	// 定时运行
	for range timer.C {

		// 清理不用的域名: 上次使用时间过大
		needClearDomains := make([]string, 0)
		for domain_name, domain := range Pool.Domain {
			domain.Lock.RLock()
			// 上次使用时间到现在为止以及大于了 12 小时
			if domain.LastUsedTime.Add(12 * time.Hour).Before(time.Now()) {
				// 需要删除的域名
				needClearDomains = append(needClearDomains, domain_name)
			}
			domain.Lock.RUnlock()
		}

		// 删除域名
		for _, domainName := range needClearDomains {
			Pool.Delete(domainName, "")
		}

		// 清理过期的 ip
		needClearIPs := make([]string, 0)
		for _, ip_meta := range BasePool {
			// 最开始获取的时间到现在已经大于了配置
			if ip_meta.InitialTime.Add(time.Duration(existTime) * time.Second).Before(time.Now()) {
				// 删除 IP
				needClearIPs = append(needClearIPs, ip_meta.Address)
				handled_count++
			}
		}

		// 删除 IP
		for _, ip := range needClearIPs {
			Pool.Delete("", ip)
		}

		// 删减
		log.Debugf("清理了 %d 个 IP", handled_count)

		// 重置
		handled_count = 0
	}
}
