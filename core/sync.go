/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-06-02 16:38:33
 * :LastEditors  :loyd
 * :LastEditTime :2023-07-14 17:56:24
 * :Description  :同步 IP 到 Domain Pool 中
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"ip-controller/log"
	"ip-controller/pool"
	"sync"
	"time"

	"github.com/pelletier/go-toml"
)

/* 动态读取配置文件 */
func SyncSettings() {
	timer := time.NewTicker(time.Second * 10)
	defer timer.Stop()

	for range timer.C {
		Settings, _ = toml.LoadFile("project.toml")
		log.Log.Debug("配置文件已更新")
	}
}

// func inDomainPool(slice []IPInDomain, ip string) bool {
// 	for _, v := range slice {
// 		if v.IP == ip {
// 			return true
// 		}
// 	}
// 	return false
// }

/* 同步新的 ip 到二级代理池 */
func SyncIP() {
	timer := time.NewTicker(time.Second * 2)
	defer timer.Stop()

	// 定时启动更新
	for range timer.C {

		// 获取一些变量
		max_life_time := Settings.Get("ip.exist_time").(int64)
		now := time.Now()

		// 删除一级代理中过期的 ip
		count_delete := 0
		ready_to_delete := make([]string, 0)
		pool.MemoryPool.Range(func(key, value any) bool {
			if time.Since(value.(time.Time)).Seconds() > float64(max_life_time) {
				ready_to_delete = append(ready_to_delete, key.(string))
			}
			return true
		})
		for _, ip := range ready_to_delete {
			pool.MemoryPool.Delete(ip)
			count_delete++
		}
		log.Log.Debug("删除过期 IP: ", count_delete, " 个")

		// 删除超过1天没有使用的域名
		to_delete := make([]string, 0)
		pool.MemorySecondaryPool.Range(func(domain_name, _ any) bool {
			domain_map, _ := pool.MemorySecondaryPool.Load(domain_name)
			last_use, _ := domain_map.(*sync.Map).Load("last_use")
			if last_use.(time.Time).Second() > int(max_life_time) {
				to_delete = append(to_delete, domain_name.(string))
			}
			return true
		})
		for _, domain_name := range to_delete {
			pool.MemorySecondaryPool.Delete(domain_name)
			log.Log.Debug("删除过期域名: ", domain_name)
		}

		// 添加新 ip 到二级代理中

		// 删除二级代理中已过期的 ip
	}

	// 	// 删除过期的 ip

	// // 同步 IP 进入 Domain Pool
	// for name, slice := range DomainPool {
	// 	tmp_new_ips := make([]IPInDomain, 0)
	// 	Pool.Range(func(p_ip, p_time interface{}) bool {
	// 		// 如果 IP 池中的 IP 不在 Domain Pool 中，则添加到临时变量中
	// 		if !inDomainPool(slice, p_ip.(string)) {
	// 			tmp_new_ips = append(tmp_new_ips, IPInDomain{p_ip.(string), p_time.(time.Time), time.Now()})
	// 		}
	// 		return true
	// 	})
	// 	old_ips := DomainPool[name]
	// 	new_ips := append(old_ips, tmp_new_ips...)
	// 	DomainPool[name] = new_ips
	// }

	// // 删除 Domain Pool 中过期的 IP
	// for name, slice := range DomainPool {
	// 	// 获取需要删除的 IP
	// 	ready_to_delete := make([]IPInDomain, 0)
	// 	for _, ip := range slice {
	// 		if time.Since(ip.AcquireTime).Seconds() > float64(max_life_time) {
	// 			ready_to_delete = append(ready_to_delete, ip)
	// 		}
	// 	}
	// 	// 遍历需要删除的 ip
	// 	for _, ip := range ready_to_delete {
	// 		for i, v := range slice {
	// 			if v.IP == ip.IP {
	// 				slice = append(slice[:i], slice[i+1:]...)
	// 			}
	// 		}
	// 	}
	// 	DomainPool[name] = slice
	// }
	// cost := time.Since(now)
	// log.Debug("完成同步, 同步耗时: ", cost)
	// Locker.Unlock()
}
