/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-06-02 16:38:33
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-03 13:48:43
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
	for range timer.C {
		Settings, _ = toml.LoadFile("project.toml")
		log.Debug("配置文件已更新")
	}
}

/* 更新 IP 池 */
func SyncPoolIP() {
	// 定时删除过期的 IP
	timer := time.NewTicker(time.Second * 10)
	for range timer.C {

		Locker.Lock()
		// 获取配置文件中每个 IP 的最大存活时间
		exist_time := Settings.Get("ip.exist_time").(int64)

		// 准备一个变量接收即将删除的 IP
		ready_to_delete := make([]interface{}, 0)

		// 遍历 IP 池, 找出过期的 IP
		Pool.Range(func(k, v interface{}) bool {
			ip_start_time := v.(time.Time)
			// 如果 IP 存活时间超过最大存活时间，则删除
			if time.Since(ip_start_time).Seconds() > float64(exist_time) {
				ready_to_delete = append(ready_to_delete, k)
			}
			return true
		})

		// 计数
		count := 0

		// 遍历变量，删除过期 IP
		for _, v := range ready_to_delete {
			Pool.Delete(v)
			count++
		}

		// 打印日志
		log.Debugf("同步 IP - 删除过期 IP %d 个", count)
		Locker.Unlock()
	}
}

/* 同步到 Domain 池 */
func SyncDomainIP() {

	timer := time.NewTicker(time.Second * 5)
	for range timer.C {

		Locker.Lock()
		// 定义一些需要的变量
		max_life_time := Settings.Get("ip.exist_time").(int64)
		now := time.Now()

		// 遍历 Domain Pool
		DomainPool.Range(func(d_name, d_pool interface{}) bool {

			// 准备遍历 domain pool
			d_pool_map := to_map(d_pool.([]DomainIP))
			// 同步新的 IP 准备添加到 Domain Pool 中
			// 1. 准备一个临时变量接收 IP
			tmp_new_ips := make([]DomainIP, 0)
			// 2. 遍历 Pool 中的 IP
			Pool.Range(func(p_ip, p_time interface{}) bool {
				// 如果 IP 池中的 IP 不在 Domain Pool 中，则添加到临时变量中
				if _, ok := d_pool_map[p_ip.(string)]; !ok {
					tmp_new_ips = append(tmp_new_ips, DomainIP{p_ip.(string), p_time.(time.Time), time.Now()})
				}
				return true
			})
			// 更新 IP
			// 读取
			old_ips, _ := DomainPool.Load(d_name.(string))
			new_ips := append(old_ips.([]DomainIP), tmp_new_ips...)
			// 设置
			DomainPool.Store(d_name.(string), new_ips)

			// 删除过期的 IP
			log.Debug("正在删除 Domain 中过期 IP")
			tmp_pool := make([]DomainIP, 0)
			for _, ip := range d_pool.([]DomainIP) {
				// 如果 IP 存活时间超过最大存活时间，跳过不进行 append 操作
				if now.Sub(ip.Start).Seconds() > float64(max_life_time) {
					continue
				}
				tmp_pool = append(tmp_pool, ip)
			}
			// 设置
			DomainPool.Store(d_name.(string), tmp_pool)

			return true
		})

		log.Debug("同步 Domain 总耗时: ", time.Since(now))
		Locker.Unlock()
	}
}

func SyncInOne() {
	go SyncSettings()
	go SyncPoolIP()
	go SyncDomainIP()
}
