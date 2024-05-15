/*
 * :Date         :2024-05-06 16:15:09
 * :LastEditTime :2024-05-15 22:25:05
 * :Description  :
 *
 */

package pool

import (
	"ip-controller/config"
	"sync"
	"time"
)

// 初始化 pool 以及 base pool
var Pool *MemoryPool
var BasePool []IP

func init() {
	Pool = NewMemoryPool()
	BasePool = make([]IP, 0)
}

type IP struct {
	Address     string               // ip 地址
	InitialTime time.Time            // ip 初始化时间
	DomainUsage map[string]time.Time // 这个 ip 被域名使用的时间
}

type Domain struct {
	Lock         *sync.RWMutex // 保护域名的锁
	LastUsedTime *time.Time    // 记录域名最后一次使用时间
	Pool         []*IP         // 记录域名的 ip 情况
}

type MemoryPool struct {
	Lock   *sync.RWMutex      // 部分操作可能设置全部， 所以需要一个锁
	Domain map[string]*Domain //
}

func NewMemoryPool() *MemoryPool {
	return &MemoryPool{
		Lock:   &sync.RWMutex{},
		Domain: make(map[string]*Domain),
	}
}

func GetMemoryPool() *MemoryPool {
	return Pool
}

// 检查当前域名、或 ip 是否存在，如果同时存在那么会检查 domain 下面的 ip
func (mp *MemoryPool) Exist(domain string, ip string) bool {

	exist := false

	switch {
	// 如果 domain 和 ip 都存在
	case domain != "" && ip != "":
		// mp 读锁
		mp.Lock.RLock()
		defer mp.Lock.RUnlock()

		// domain 存在，那么遍历 pool
		if _, ok := mp.Domain[domain]; ok {
			mp.Domain[domain].Lock.RLock()
			for _, v := range mp.Domain[domain].Pool {
				if v.Address == ip {
					exist = true
				}
			}
			mp.Domain[domain].Lock.RUnlock()
		}

	// 如果 domain 不存在， ip 存在
	case domain == "" && ip != "":
		// 读锁
		mp.Lock.RLock()
		defer mp.Lock.RUnlock()

		// 直接检查 base pool
		for _, v := range BasePool {
			if v.Address == ip {
				exist = true
			}
		}

	// 如果 domain 存在， ip 不存在
	case domain != "" && ip == "":
		// 读锁
		mp.Lock.RLock()
		defer mp.Lock.RUnlock()

		if _, ok := mp.Domain[domain]; ok {
			exist = true
		}
	}

	return exist

}

// 创建新域名
func (mp *MemoryPool) Create(domain string) bool {

	// 创建先加锁判断是否存在
	mp.Lock.Lock()
	defer mp.Lock.Unlock()

	// 检查 domain 是否存在, 不存在才添加
	if _, ok := mp.Domain[domain]; !ok {

		// 复制 base pool 并创建新的 tmp pool
		tmpPool := make([]*IP, 0)
		for _, v := range BasePool {
			tmpPool = append(tmpPool, &IP{
				Address:     v.Address,
				InitialTime: v.InitialTime,
				DomainUsage: make(map[string]time.Time),
			})
		}

		// 创建
		mp.Domain[domain] = &Domain{
			Lock: &sync.RWMutex{},
			Pool: tmpPool,
		}

		return true
	}

	// 存在返回 false
	return false
}

// 删除域名或 ip
func (mp *MemoryPool) Delete(domain string, ip string) {

	// 加锁
	mp.Lock.Lock()
	defer mp.Lock.Unlock()

	switch {
	// domain 和 ip 都存在, 去 domain 下面删除 ip
	case domain != "" && ip != "":
		// domain 存在
		if _, ok := mp.Domain[domain]; ok {
			// domain lock
			mp.Domain[domain].Lock.Lock()
			newPool := make([]*IP, 0)
			for _, v := range mp.Domain[domain].Pool {
				if v.Address != ip {
					newPool = append(newPool, v)
				}
			}
			mp.Domain[domain].Pool = newPool
			mp.Domain[domain].Lock.Unlock()
		}
	// domain 存在， ip 不存在，删除 domain
	case domain != "" && ip == "":
		delete(mp.Domain, domain)
	// domain 不存在， ip 存在，删除 base pool 里面的 ip 以及 domain 里面的 ip
	case domain == "" && ip != "":
		// 删除 base pool 里面的 ip
		newBasePool := make([]IP, 0)
		for _, v := range BasePool {
			if v.Address != ip {
				newBasePool = append(newBasePool, v)
			}
		}
		BasePool = newBasePool
		// 删除 domain 里面的 ip, 先判断 domain 是否存在
		if _, ok := mp.Domain[domain]; ok {
			// domain lock
			mp.Domain[domain].Lock.Lock()
			newPool := make([]*IP, 0)
			for _, v := range mp.Domain[domain].Pool {
				if v.Address != ip {
					newPool = append(newPool, v)
				}
			}
			mp.Domain[domain].Pool = newPool
			mp.Domain[domain].Lock.Unlock()
		}
	}
}

// 添加 ip
func (mp *MemoryPool) Add(ip string) {

	// 加锁
	mp.Lock.Lock()
	defer mp.Lock.Unlock()

	// 添加 ip 到 base pool
	BasePool = append(BasePool, IP{
		Address:     ip,
		InitialTime: time.Now(),
		DomainUsage: make(map[string]time.Time),
	})

	// 添加 ip 到 domain pool
	for _, v := range mp.Domain {
		// domain lock
		v.Lock.Lock()
		v.Pool = append(v.Pool, &IP{
			Address:     ip,
			InitialTime: time.Now(),
			DomainUsage: make(map[string]time.Time),
		})
		v.Lock.Unlock()
	}
}

// 获取 ip
func (mp *MemoryPool) Get(domain string) string {

	// 判断 domain 是否存在, 不存在则先创建
	if _, ok := mp.Domain[domain]; !ok {
		mp.Create(domain)
	}

	// 获取设置
	maxTime := config.Settings.Get("pool.delay").(int64) // 延迟时间

	// 设置返回参数
	ipAddress := ""

	// 读锁
	mp.Domain[domain].Lock.RLock()
	// domain 存在，遍历 pool 找到合适的 ip
FindLoop:
	for _, v := range mp.Domain[domain].Pool {

		// 检查 domain usage 可能这个不存在, 因为是第一次使用
		if _, ok := v.DomainUsage[domain]; !ok {
			v.DomainUsage[domain] = time.Now()
			// 更新 domain last used time
			lastUsedTime := time.Now()
			mp.Domain[domain].LastUsedTime = &lastUsedTime
			ipAddress = v.Address
			break FindLoop
		}

		// domain usage 存在, 肯定存在时间, 检查是否超过最大时间
		if time.Since(v.DomainUsage[domain]).Seconds() > float64(maxTime) {
			v.DomainUsage[domain] = time.Now()
			// 更新 domain last used time
			lastUsedTime := time.Now()
			mp.Domain[domain].LastUsedTime = &lastUsedTime
			ipAddress = v.Address
			break FindLoop
		}

	}
	mp.Domain[domain].Lock.RUnlock()

	return ipAddress
}
