/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-06-02 16:38:33
 * :LastEditors  :loyd
 * :LastEditTime :2023-07-15 16:19:37
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

func cleanPrimaryProxyPool(max_life_time int64) {

	// 定义变量接收需要删除的 ip，与计数
	count_delete := 0
	to_delete := make([]string, 0)

	// 遍历一级代理池
	pool.MemoryPool.Range(func(key, value any) bool {

		// 转义变量
		ip := key.(string)
		acquire_time := value.(time.Time)

		// 判断是否过期
		if time.Since(acquire_time).Seconds() > float64(max_life_time) {
			to_delete = append(to_delete, ip)
		}
		return true
	})
	for _, ip := range to_delete {
		pool.MemoryPool.Delete(ip)
		count_delete++
	}
	log.Log.Debug("删除过期 IP: ", count_delete, " 个")
}

func cleanNotUsedSecondaryProxyPool(max_hour int) {

	// 定义变量接收需要删除的 ip
	to_delete := make([]string, 0)
	count_delete := 0

	// 遍历二级代理
	pool.MemorySecondaryPool.Range(func(secondary_pool_name, secondary_pool any) bool {

		// 二级代理池
		pt_secondary_pool := secondary_pool.(*sync.Map)

		// 获取二级代理上次的使用时间
		last_use, ok := pt_secondary_pool.Load("last_use")
		if !ok {
			log.Log.Debug("获取二级代理上次使用时间失败")
			return true
		}

		// 判断是否过期
		if time.Since(last_use.(time.Time)).Seconds() > float64(max_hour*60*60) {
			to_delete = append(to_delete, secondary_pool_name.(string))
		}
		return true
	})

	// 遍历需要删除的二级代理名
	for _, secondary_pool_name := range to_delete {
		pool.MemorySecondaryPool.Delete(secondary_pool_name)
		count_delete++
	}
	log.Log.Debug("删除过期域名: ", count_delete, " 个")
}

func addNewIPToSecondaryProxyPool() {

	// 遍历一级代理 ip 池
	pool.MemoryPool.Range(func(key, value any) bool {

		// 获取 ip 和初始时间
		ip := key.(string)
		acquire_time := value.(time.Time)

		// 遍历每个二级代理
		pool.MemorySecondaryPool.Range(func(secondary_pool_name, secondary_pool any) bool {

			// 二级代理池
			pt_secondary_pool := secondary_pool.(*sync.Map)

			// 获取二级代理池的全部 ip
			ip_pool, ok := pt_secondary_pool.Load("ip_pool")
			// 如果获取时发生错误，处理下一个
			if !ok {
				log.Log.Error("向二级代理添加ip时, 获取二级代理池失败")
				return true
			}

			// 使用一级代理的 ip 来判断是否在二级代理池中
			pt_ip_pool := ip_pool.(*sync.Map)
			_, ok = pt_ip_pool.Load(ip)
			// ip 已经存在
			if ok {
				return true
			}

			// 添加这个 ip
			info := map[string]interface{}{
				"acquire_time": acquire_time,
				"last_use":     time.Now(),
			}
			pt_ip_pool.Store(ip, info)

			// 继续处理下一个二级代理的 ip
			return true
		})

		// 进入下一个一级代理的 ip
		return true
	})
}

func cleanSecondaryProxyPool(max_life_time int64) {

	// 遍历二级代理池
	pool.MemorySecondaryPool.Range(func(secondary_pool_name, secondary_pool any) bool {

		// 二级代理池
		pt_secondary_pool := secondary_pool.(*sync.Map)

		// 获取二级代理池的全部 ip
		ip_pool, ok := pt_secondary_pool.Load("ip_pool")
		// 如果获取时发生错误，处理下一个
		if !ok {
			log.Log.Error("清理二级代理时, 获取二级代理池失败")
			return true
		}

		// 获取二级代理池全部 ip 的指针
		pt_ip_pool := ip_pool.(*sync.Map)

		// 遍历池中的全部 ip
		to_delete := make([]string, 0)
		pt_ip_pool.Range(func(ip, info any) bool {
			// 获取 ip 的信息
			ip_info := info.(map[string]interface{})
			// 判断是否过期
			if time.Since(ip_info["acquire_time"].(time.Time)).Seconds() > float64(max_life_time) {
				to_delete = append(to_delete, ip.(string))
			}
			return true
		})
		// 遍历需要删除 ip, 最后在 pointer 中删除
		for _, ip := range to_delete {
			pt_ip_pool.Delete(ip)
		}

		// 继续处理下一个二级代理的
		return true
	})
}

/* 同步新的 ip 到二级代理池 */
func SyncIP() {
	timer := time.NewTicker(time.Second * 2)
	defer timer.Stop()

	// 定时启动更新
	for range timer.C {

		// 获取一些变量
		max_life_time := Settings.Get("ip.exist_time").(int64)

		// 删除一级代理池过期的 ip
		cleanPrimaryProxyPool(max_life_time)

		// 删除超过一定时间的二级代理
		cleanNotUsedSecondaryProxyPool(12)

		// 添加新 ip 到二级代理中
		addNewIPToSecondaryProxyPool()

		// 删除二级代理中已过期的 ip
		cleanSecondaryProxyPool(max_life_time)
	}

}
