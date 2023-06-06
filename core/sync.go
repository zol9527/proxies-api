/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-06-02 16:38:33
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-06 11:49:45
 * :Description  :同步 IP 到 Domain Pool 中
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"time"

	"github.com/pelletier/go-toml"
)

/* 动态读取配置文件 */
func SyncSettings() {
	timer := time.NewTicker(time.Second * 10)
	defer timer.Stop()

	for range timer.C {
		Settings, _ = toml.LoadFile("project.toml")
		log.Debug("配置文件已更新")
	}
}

func inDomainPool(slice []DomainIP, ip string) bool {
	for _, v := range slice {
		if v.IP == ip {
			return true
		}
	}
	return false
}

func SyncIP() {
	timer := time.NewTicker(time.Second * 2)
	defer timer.Stop()

	for range timer.C {

		Locker.Lock()
		log.Debug("开始同步 IP")

		// 获取一些变量
		max_life_time := Settings.Get("ip.exist_time").(int64)
		now := time.Now()

		// 更新 ip
	Loop:
		for {
			select {
			case ip := <-SyncPipe:
				Pool.Store(ip.IP, ip.AcquireTime)
			default:
				break Loop
			}
		}

		// 删除过期的 ip
		count_delete := 0
		ready_to_delete := make([]NewIP, 0)
		Pool.Range(func(ip, t interface{}) bool {
			if time.Since(t.(time.Time)).Seconds() > float64(max_life_time) {
				ready_to_delete = append(ready_to_delete, NewIP{ip.(string), t.(time.Time)})
			}
			return true
		})
		for _, ip := range ready_to_delete {
			Pool.Delete(ip.IP)
			count_delete++
		}
		log.Debug("删除过期 IP: ", count_delete, " 个")

		// 同步 IP 进入 Domain Pool
		for name, slice := range DomainPool {
			tmp_new_ips := make([]DomainIP, 0)
			Pool.Range(func(p_ip, p_time interface{}) bool {
				// 如果 IP 池中的 IP 不在 Domain Pool 中，则添加到临时变量中
				if !inDomainPool(slice, p_ip.(string)) {
					tmp_new_ips = append(tmp_new_ips, DomainIP{p_ip.(string), p_time.(time.Time), time.Now()})
				}
				return true
			})
			old_ips := DomainPool[name]
			new_ips := append(old_ips, tmp_new_ips...)
			DomainPool[name] = new_ips
		}

		// 删除 Domain Pool 中过期的 IP
		for name, slice := range DomainPool {
			// 获取需要删除的 IP
			ready_to_delete := make([]DomainIP, 0)
			for _, ip := range slice {
				if time.Since(ip.AcquireTime).Seconds() > float64(max_life_time) {
					ready_to_delete = append(ready_to_delete, ip)
				}
			}
			// 遍历需要删除的 ip
			for _, ip := range ready_to_delete {
				for i, v := range slice {
					if v.IP == ip.IP {
						slice = append(slice[:i], slice[i+1:]...)
					}
				}
			}
			DomainPool[name] = slice
		}
		cost := time.Since(now)
		log.Debug("完成同步, 同步耗时: ", cost)
		Locker.Unlock()
	}
}
