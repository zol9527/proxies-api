/*
 * :Date         :2024-05-06 16:15:09
 * :LastEditTime :2024-05-07 20:35:38
 * :Description  :
 *
 */

package pool

import (
	"ip-controller/config"
	"sync"
	"time"
)

// 维护一个临时的 ip 池作为基准
var BasePool map[string]*IP
var Pool *MemoryPool

func init() {
	// 初始化临时 ip 池，以及 pool
	Pool = NewMemoryPool()
	BasePool = make(map[string]*IP, 0)
}

type MemoryPool struct {
	Lock sync.Mutex
	Pool map[string]*Domain
}

type Domain struct {
	LastUsedTime time.Time
	IP           map[string]*IP
}

type IP struct {
	InitialTime time.Time
	DomainUsage map[string]time.Time
}

// 创建内存池
func NewMemoryPool() *MemoryPool {
	return &MemoryPool{
		Pool: make(map[string]*Domain),
	}
}

// 检查当前域名是否存在
func DomainExist(domain string) bool {

	// 创建先加锁判断是否存在
	Pool.Lock.Lock()
	defer Pool.Lock.Unlock()

	if _, ok := Pool.Pool[domain]; ok {
		return true
	} else {
		return false
	}
}

// 创建新域名
func CreateDomain(domain string) bool {

	// 创建先加锁判断是否存在
	Pool.Lock.Lock()
	defer Pool.Lock.Unlock()

	if _, ok := Pool.Pool[domain]; !ok {
		Pool.Pool[domain] = &Domain{
			IP:          BasePool,
		}
		return true
	}
	return false
}

// 删除域名
func DeleteDomain(domain string) bool {
	Pool.Lock.Lock()
	defer Pool.Lock.Unlock()

	if _, ok := Pool.Pool[domain]; ok {
		delete(Pool.Pool, domain)
		return true
	}
	return false
}

// 检查 ip 是否存在于一个 domain 中
func IPExist(ip string) bool {

	Pool.Lock.Lock()
	defer Pool.Lock.Unlock()

	if _, ok := BasePool[ip]; ok {
		return true
	}

	return false
}

// 向域名中添加 ip
func AddIP(ip string) bool {

	Pool.Lock.Lock()
	defer Pool.Lock.Unlock()

	// 判断是否在 base ip 中
	if _, ok := BasePool[ip]; !ok {
		BasePool[ip] = &IP{
			InitialTime: time.Now(),
		}
		return true
	}

	return false
}

// 删除域名中的 ip
func DeleteIP(ip string) bool {

	Pool.Lock.Lock()
	defer Pool.Lock.Unlock()

	if _, ok := BasePool[ip]; !ok {
		return false
	}

	// 删除 ip
	delete(BasePool, ip)

	return true
}

// 获取域名中的 ip
func GetIP(domain string) string {

	Pool.Lock.Lock()
	defer Pool.Lock.Unlock()

	// 判断域名是否存在
	if _, ok := Pool.Pool[domain]; !ok {
		return ""
	}

	// 遍历 map 查看 ip last use time 大雨 toml 设置的, 如果没有这返回 ""
	max_time := config.Settings.Get("pool.delay").(int64)

	for ip, v := range Pool.Pool[domain].IP {

		// 获取这个 ip 的上次域名使用情况
		if v.DomainUsage[domain].IsZero() {
			v.DomainUsage[domain] = time.Now()
			Pool.Pool[domain].LastUsedTime = time.Now()
			return ip
		}

		if time.Since(v.DomainUsage[domain]).Seconds() > float64(max_time) {
			v.DomainUsage[domain] = time.Now()
			Pool.Pool[domain].LastUsedTime = time.Now()
			return ip
		}
	}

	return ""
}
