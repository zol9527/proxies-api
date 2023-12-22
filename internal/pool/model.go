/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-18 16:17:02
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-22 12:00:53
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import (
	"ip-controller/config"
	"ip-controller/internal/log"
	"math/rand"
	"sync"
	"time"
)

var pool *Pool = &Pool{
	Pool: make(map[string]*Domain),
}

// 池
type Pool struct {
	Lock sync.Mutex
	Pool map[string]*Domain
}

// 域基础信息
type Domain struct {
	Lock sync.Mutex
	Meta *DomainMeta
}

// 域详细信息
type DomainMeta struct {
	Last_used_time time.Time
	Memory_pool    map[string]*IP
}

// IP 信息
type IP struct {
	Initial_time   time.Time
	Last_used_time time.Time
}

// 获取 pool 指针对象
func GetPool() *Pool {
	return pool
}

func DomainExist(domain string) bool {
	if _, ok := pool.Pool[domain]; ok {
		return true
	} else {
		return false
	}
}

// 添加域
func CreateDomain(domain string) bool {

	// 加锁
	pool.Lock.Lock()
	pool.Lock.Unlock()

	// 查看 pool 中是否存在这个域
	if _, ok := pool.Pool[domain]; !ok {

		// 初始化 domain
		now := time.Now()
		new_meta := &DomainMeta{
			Last_used_time: now,
			Memory_pool:    make(map[string]*IP),
		}

		// 创建
		pool.Pool[domain] = &Domain{
			Meta: new_meta,
		}

		return true
	} else {
		// 返回 nil 代表已存在
		return false
	}
}

// 删除域
func DeleteDomain(domain string) bool {

	// 加锁
	pool.Lock.Lock()
	defer pool.Lock.Unlock()

	if _, ok := pool.Pool[domain]; ok {
		delete(pool.Pool, domain)
		return true
	} else {
		return false
	}
}

// 添加 IP 到 *全部* 域
func AddIP(ip string) bool {

	// 加锁
	pool.Lock.Lock()
	defer pool.Lock.Unlock()

	// 遍历全部 domain, 添加 ip
	for _, domain := range pool.Pool {

		// domain 操作加锁
		domain.Lock.Lock()

		// 添加进入 memory_pool
		now := time.Now()
		domain.Meta.Memory_pool[ip] = &IP{
			Initial_time:   now,
			Last_used_time: now,
		}

		domain.Lock.Unlock()
	}

	return true
}

// 删除 IP
func DeleteIP(domain string, ip string) bool {

	// 判断域名存在
	if _, ok := pool.Pool[domain]; ok {

		// 加锁
		pool.Pool[domain].Lock.Lock()

		// 查看 IP 是否已经在里面
		if _, ok := pool.Pool[domain].Meta.Memory_pool[ip]; ok {
			delete(pool.Pool[domain].Meta.Memory_pool, ip)
		} else {
			log.Debugln("Delete IP - IP not exists in the domain")
		}
		pool.Pool[domain].Lock.Unlock()
		return true
	}

	log.Debugln("Delete IP - Domain not exists")
	return false
}

// 获取 IP
func GetIP(domain string) string {

	// 加锁
	pool.Pool[domain].Lock.Lock()
	defer pool.Pool[domain].Lock.Unlock()

	// 判断域名存在
	if _, ok := pool.Pool[domain]; ok {

		// 判断 memory_pool 大小
		size := len(pool.Pool[domain].Meta.Memory_pool)

		// 如果 memory_pool 为空，则返回空字符串
		if size == 0 {
			log.Debugln("Get IP - 无可用 ip")
			return ""
		}

		// 序列化 keys
		keys := make([]string, 0)
		for k := range pool.Pool[domain].Meta.Memory_pool {
			keys = append(keys, k)
		}

		// 随机获取
		indexKey := rand.Intn(size)
		ip := keys[indexKey]

		// 判断 ip 存活时间
		ip_struct := pool.Pool[domain].Meta.Memory_pool[ip]
		// 延迟时间
		delay := config.Settings.Get("pool.delay").(int64)

		// 比较上一次使用时间是否超过 last time
		if time.Now().Unix()-ip_struct.Last_used_time.Unix() > delay {

			// 更新 last use time
			pool.Pool[domain].Meta.Memory_pool[ip].Last_used_time = time.Now()
			return ip
		} else {
			return ""
		}
	} else {
		log.Debugln("Get IP - Domain not exists")
		return ""
	}
}
