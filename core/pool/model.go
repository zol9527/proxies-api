/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-18 16:17:02
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-21 16:17:16
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import (
	"ip-controller/config"
	"ip-controller/core/log"
	"math/rand"
	"sync"
	"time"
)

var Pool *PoolStruct = &PoolStruct{
	Pool: make(map[string]*DomainStruct),
}

// 池
type PoolStruct struct {
	Lock sync.Mutex
	Pool map[string]*DomainStruct
}

// 域基础信息
type DomainStruct struct {
	Lock sync.Mutex
	Meta *DomainMetaStruct
}

// 域详细信息
type DomainMetaStruct struct {
	Last_used_time time.Time
	Memory_pool    map[string]*IPStruct
}

// IPStruct 信息
type IPStruct struct {
	Initial_time   time.Time
	Last_used_time time.Time
}

// 添加域
func (p *PoolStruct) Create(domain string) bool {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	// 查看 pool 中是否存在，不存在则创建
	if _, ok := p.Pool[domain]; !ok {

		// 初始化 domain
		now := time.Now()
		new_meta := &DomainMetaStruct{
			Last_used_time: now,
			Memory_pool:    make(map[string]*IPStruct),
		}

		// 创建
		p.Pool[domain] = &DomainStruct{
			Meta: new_meta,
		}

		return true
	} else {
		// 返回 nil 代表已存在
		return false
	}
}

// 删除域
func (p *PoolStruct) Delete(domain string) bool {

	// 加锁
	p.Lock.Lock()
	defer p.Lock.Unlock()

	if _, ok := p.Pool[domain]; ok {
		delete(p.Pool, domain)
		return true
	} else {
		return false
	}
}

// 添加 IP 到 *全部* 域
func (p *PoolStruct) AddIP(ip string) bool {

	// 加锁
	p.Lock.Lock()
	defer p.Lock.Unlock()

	// 遍历全部 domain, 添加 ip
	for _, domain := range p.Pool {

		// domain 操作加锁
		domain.Lock.Lock()

		// 添加进入 memory_pool
		now := time.Now()
		domain.Meta.Memory_pool[ip] = &IPStruct{
			Initial_time:   now,
			Last_used_time: now,
		}

		domain.Lock.Unlock()
	}

	return true
}

// 删除 IP
func (p *PoolStruct) DeleteIP(domain string, ip string) bool {

	// 判断域名存在
	if _, ok := p.Pool[domain]; ok {

		func() {

			// 加锁
			p.Pool[domain].Lock.Lock()
			defer p.Pool[domain].Lock.Unlock()

			// 查看 IP 是否已经在里面
			if _, ok := p.Pool[domain].Meta.Memory_pool[ip]; ok {
				delete(p.Pool[domain].Meta.Memory_pool, ip)
			} else {
				log.Debugln("Delete IP - IP not exists in the domain")
			}

		}()

		return true
	}

	log.Debugln("Delete IP - Domain not exists")
	return false
}

// 获取 IP
func (p *PoolStruct) GetIP(domain string) string {

	// 加锁
	p.Pool[domain].Lock.Lock()
	defer p.Pool[domain].Lock.Unlock()

	// 判断域名存在
	if _, ok := p.Pool[domain]; ok {

		// 判断 memory_pool 大小
		size := len(p.Pool[domain].Meta.Memory_pool)

		// 如果 memory_pool 为空，则返回空字符串
		if size == 0 {
			return ""
		}

		// 序列化 keys
		keys := make([]string, 0)
		for k := range p.Pool[domain].Meta.Memory_pool {
			keys = append(keys, k)
		}

		// 随机获取
		indexKey := rand.Intn(size)
		ip := keys[indexKey]

		// 判断 ip 存活时间
		ip_struct := p.Pool[domain].Meta.Memory_pool[ip]
		// 延迟时间
		delay := config.Settings.Get("pool.delay").(int64)

		// 比较上一次使用时间是否超过 last time
		if time.Now().Unix()-ip_struct.Last_used_time.Unix() > delay {
			return ip
		} else {
			return ""
		}
	} else {
		log.Debugln("Get IP - Domain not exists")
		return ""
	}
}
