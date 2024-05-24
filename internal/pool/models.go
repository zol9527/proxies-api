/*
 * :Date         :2024-05-23 14:02:55
 * :LastEditTime :2024-05-24 11:10:02
 * :Description  :
 *
 */

package pool

import (
	"errors"
	"ip-controller/config"
	"sync"
	"time"
)

// 错误类型
var (
	ErrDomainExist    = errors.New("domain exist")
	ErrDomainNotExist = errors.New("domain not exist")
	ErrIPNotInDomain  = errors.New("ip not in domain")
	ErrNoMatchCases   = errors.New("no match cases")
	ErrEmptyInput     = errors.New("empty input")
	ErrIPAlreadyExist = errors.New("ip already exist")
	ErrNoIPMatched    = errors.New("ip not matched")
)

var (
	Pool       *MemoryPool
	BaseIPList []IP // 存储新的 ip
	PLock      sync.RWMutex
)

// 初始化 pool 以及 base pool
func init() {
	Pool = NewMemoryPool()
	BaseIPList = make([]IP, 0)
}

type IP struct {
	Address      string
	InitialTime  time.Time
	LastUsedTime time.Time
}

type Domain struct {
	IPList       []*IP        // 域名列表
	DomainName   string       // 域名
	Lock         sync.RWMutex // 保护域名的锁
	LastUsedTime time.Time    // 最后一次使用时间
}

type MemoryPool struct {
	Lock      sync.RWMutex       // 保护域名的锁
	DomainMap map[string]*Domain // 域名列表
}

// 创建 pool
func NewMemoryPool() *MemoryPool {
	return &MemoryPool{
		Lock:      sync.RWMutex{},
		DomainMap: make(map[string]*Domain),
	}
}

// 检查存在
func (mp *MemoryPool) Exist(domain string, ip string) error {

	switch {

	// 同时存在则去 domain pool list 中查找 ip 是否存在
	case domain != "" && ip != "":
		mp.Lock.RLock()
		defer mp.Lock.RUnlock()

		// 查询 domain 存在
		if _, ok := mp.DomainMap[domain]; ok {
			mp.DomainMap[domain].Lock.RLock()
			defer mp.DomainMap[domain].Lock.RUnlock()
			// 遍历 domain 下面的 ip list
			for _, v := range mp.DomainMap[domain].IPList {
				if v.Address == ip {
					return nil
				}
			}
			return ErrIPNotInDomain
		}
		return ErrDomainNotExist

	// 只检查 ip 是否存在
	case domain == "" && ip != "":

		PLock.RLock()
		defer PLock.RUnlock()

		// 遍历 base pool
		for _, v := range BaseIPList {
			if v.Address == ip {
				return nil
			}
		}
		return ErrIPNotInDomain

	case domain != "" && ip == "":
		mp.Lock.RLock()
		defer mp.Lock.RUnlock()

		if _, ok := mp.DomainMap[domain]; ok {
			return nil
		}

		return ErrDomainNotExist

	case domain == "" && ip == "":
		return ErrEmptyInput
	}

	return ErrNoMatchCases
}

// 删除
func (mp *MemoryPool) Delete(domain string, ip string) error {

	switch {

	// both domain and ip exist
	case domain != "" && ip != "":
		// 先检查 domain 是否存在，再决定是否加锁
		if _, ok := mp.DomainMap[domain]; !ok {
			return nil
		}
		// domain 存在
		mp.DomainMap[domain].Lock.Lock()
		defer mp.DomainMap[domain].Lock.Unlock()
		// 遍历 domain 下面的 pool list, 删除这个 ip
		for i, v := range mp.DomainMap[domain].IPList {
			if v.Address == ip {
				mp.DomainMap[domain].IPList = append(mp.DomainMap[domain].IPList[:i], mp.DomainMap[domain].IPList[i+1:]...)
				return nil
			}
		}
		return nil
	// only ip exist, go to delete ip in base pool and domain ip list
	case ip != "" && domain == "":
		// 删除 ip 从 base pool 中
		PLock.Lock()
		for i, v := range BaseIPList {
			if v.Address == ip {
				BaseIPList = append(BaseIPList[:i], BaseIPList[i+1:]...)
			}
		}
		PLock.Unlock()

		// 遍历 domain, 删除每个 domain 下面的这个 ip
		for _, v := range mp.DomainMap {
			v.Lock.Lock()
			for i, per_ip := range v.IPList {
				if per_ip.Address == ip {
					v.IPList = append(v.IPList[:i], v.IPList[i+1:]...)
				}
			}
			v.Lock.Unlock()
		}
		return nil
	// only domain exist
	case domain != "" && ip == "":
		mp.Lock.Lock()
		defer mp.Lock.Unlock()
		if _, ok := mp.DomainMap[domain]; ok {
			delete(mp.DomainMap, domain)
			return nil
		} else {
			return ErrDomainNotExist
		}
	case domain == "" && ip == "":
		return ErrEmptyInput
	}

	return ErrNoMatchCases
}

// 添加
func (mp *MemoryPool) Add(domain, ip string) error {

	switch {

	// both domain and ip exist, 这里会向 domain 下面的 ip list 添加
	case domain != "" && ip != "":
		// 先检查 domain 是否存在，再决定是否加锁
		if _, ok := mp.DomainMap[domain]; !ok {
			return ErrDomainNotExist
		}
		// domain 存在
		mp.DomainMap[domain].Lock.Lock()
		defer mp.DomainMap[domain].Lock.Unlock()
		// 遍历 domain 下面的 pool list, 查看这个 ip 是否存在，不存在再添加
		for _, v := range mp.DomainMap[domain].IPList {
			if v.Address == ip {
				return ErrIPAlreadyExist
			}
		}
		// append ip to ip list
		mp.DomainMap[domain].IPList = append(mp.DomainMap[domain].IPList, &IP{
			Address:     ip,
			InitialTime: time.Now(),
		})
		return nil
	// only ip exist, go to add ip in base pool and domain ip list
	case domain == "" && ip != "":
		// 添加 ip 从 base pool 中
		PLock.Lock()
		BaseIPList = append(BaseIPList, IP{
			Address:     ip,
			InitialTime: time.Now(),
		})
		PLock.Unlock()

		// 添加进入每个 domain
		for _, v := range mp.DomainMap {
			v.Lock.Lock()
			v.IPList = append(v.IPList, &IP{
				Address:     ip,
				InitialTime: time.Now(),
			})
			v.Lock.Unlock()
		}
		return nil

	// only domain exist
	case domain != "" && ip == "":
		mp.Lock.Lock()
		defer mp.Lock.Unlock()
		if _, ok := mp.DomainMap[domain]; !ok {

			// 创建 domain 基础对象
			mp.DomainMap[domain] = &Domain{
				IPList:     make([]*IP, 0),
				DomainName: domain,
			}

			mp.DomainMap[domain].Lock.Lock()
			defer mp.DomainMap[domain].Lock.Unlock()

			// 添加 ip 到 domain, 复制 base pool 到 domain
			PLock.RLock()
			for _, v := range BaseIPList {
				mp.DomainMap[domain].IPList = append(mp.DomainMap[domain].IPList, &IP{
					Address:     v.Address,
					InitialTime: v.InitialTime,
				})
			}
			PLock.RUnlock()
			return nil
		} else {
			return ErrDomainExist
		}
	}

	return ErrNoMatchCases
}

// 获取
func (mp *MemoryPool) Get(domain string) (string, error) {

	// domain 不存在
	if _, ok := mp.DomainMap[domain]; !ok {
		return "", ErrDomainNotExist
	}

	// domain 存在
	mp.DomainMap[domain].Lock.RLock()
	defer mp.DomainMap[domain].Lock.RUnlock()

	// get delay time
	delay := time.Duration(config.Settings.Get("pool.delay").(int64))

	// 遍历 domain 下面的 ip list, 获取符合时间条件的 ip, last use time 到现在必须大于 15 秒
	for _, v := range mp.DomainMap[domain].IPList {

		// last use time is zero, a completely new ip
		if v.LastUsedTime.IsZero() {
			v.LastUsedTime = time.Now()
			return v.Address, nil
		}

		// last use time is not zero, check if it is greater than delay time
		if v.LastUsedTime.Add(delay * time.Second).Before(time.Now()) {
			v.LastUsedTime = time.Now()
			return v.Address, nil
		}
	}

	return "", ErrNoIPMatched
}
