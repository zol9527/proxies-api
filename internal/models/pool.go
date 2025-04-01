package models

import (
	"proxies-api/internal/log"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
)

// Pool 结构体定义了代理IP池的核心数据结构
type Pool struct {
	size           int64                                       // 池的大小
	ip_ttl         int64                                       // 每个IP的生命周期 seconds
	origin_ip_map  *maputil.ConcurrentMap[string, *IpDetails]  // 原始IP映射
	domains        *maputil.ConcurrentMap[string, *DomainInfo] // 域名信息映射
	clean_ticker   *time.Ticker                                // 清理IP的定时器
	ticker_stopped bool                                        // 标记定时器是否已停止
	ticker_mutex   sync.Mutex                                  // 保护定时器操作的互斥锁
	stop_chan      chan struct{}                               // 用于停止定时器的通道
}

/**
 * Function: GetOriginIPs
 *
 * Purpose: 获取池中所有原始IP及其详情的映射
 *
 * Parameters:
 *
 * Returns: map[string]*IpDetails - IP地址到IP详情对象的映射，包含池中所有原始IP信息
 */
func (p *Pool) GetOriginIPs() map[string]*IpDetails {
	ips := make(map[string]*IpDetails)
	p.origin_ip_map.Range(func(key string, value *IpDetails) bool {
		ips[key] = value
		return true
	})
	return ips
}

/**
 * Function: CleanTicker
 *
 * Purpose: 启动一个定时器定期清理池中过期的IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: 无
 */
func (p *Pool) CleanTicker() {
	p.ticker_mutex.Lock()
	defer p.ticker_mutex.Unlock()

	// 如果已经有定时器在运行，则返回
	if p.clean_ticker != nil {
		log.Debug("⚠️ 清理定时器已经在运行中")
		return
	}

	// 创建定时器和停止通道
	p.clean_ticker = time.NewTicker(10 * time.Second)
	p.stop_chan = make(chan struct{})
	p.ticker_stopped = false

	// 在后台启动清理任务
	go func() {
		for {
			select {
			case <-p.clean_ticker.C:
				// 执行清理操作
				p.cleanExpiredIPs()
			case <-p.stop_chan:
				// 收到停止信号，结束goroutine
				log.Debug("🛑 IP清理定时器接收到停止信号")
				return
			}
		}
	}()

	log.Debug("🕒 IP清理定时器已启动，每10秒清理一次过期IP")
}

/**
 * Function: cleanExpiredIPs
 *
 * Purpose: 清理所有域名中过期的IP
 *
 * Parameters:
 *   - 无
 *
 * Returns: 无
 */
func (p *Pool) cleanExpiredIPs() {
	// 遍历所有域名
	p.domains.Range(func(key string, value *DomainInfo) bool {
		// 检查域名的IP池是否为空
		if value.GetIPCount() == 0 {
			log.Debug("❌ 域名 %s 的IP池为空，跳过清理", key)
			return true
		}

		// 清理过期IP
		value.ip_details_map.Range(func(ip string, ipDetails *IpDetails) bool {
			if time.Since(ipDetails.last_use_time).Seconds() > float64(p.ip_ttl) {
				value.Delete(ipDetails)
				log.Debug("✅ 清理过期IP %s 成功", ip)
			}
			return true
		})

		return true
	})
}

/**
 * Function: StopCleanTicker
 *
 * Purpose: 停止IP清理定时器
 *
 * Parameters:
 *   - 无
 *
 * Returns: bool - 是否成功停止定时器
 */
func (p *Pool) StopCleanTicker() bool {
	p.ticker_mutex.Lock()
	defer p.ticker_mutex.Unlock()

	// 检查定时器是否已存在且未停止
	if p.clean_ticker != nil && !p.ticker_stopped {
		// 停止定时器
		p.clean_ticker.Stop()
		// 发送停止信号给goroutine
		close(p.stop_chan)
		// 标记为已停止
		p.ticker_stopped = true
		log.Debug("✅ IP清理定时器已成功停止")
		return true
	}

	log.Debug("⚠️ IP清理定时器未运行或已停止")
	return false
}

/**
 * Function: AddIPToPool
 *
 * Purpose: 向池中添加新的IP和端口，并分配到所有域名
 *
 * Parameters:
 *   - ip (string): 要添加的IP地址
 *   - port (int): 要添加的端口号
 *
 * Returns: 无
 */
func (p *Pool) AddIPToPool(ip string, port int) {
	// 添加原始IP到池中
	ipDetails := NewIpDetails(ip, port)
	p.origin_ip_map.Set(ip, ipDetails)
	log.Debug("✅ 添加原始IP %s:%d 成功", ip, port)

	// 添加到全部 domains
	p.domains.Range(func(key string, value *DomainInfo) bool {
		value.Add(ipDetails)
		return true
	})
}

/**
 * Function: DeleteIPFromPool
 *
 * Purpose: 从池中删除指定IP，并从所有域名中移除该IP
 *
 * Parameters:
 *   - ip (string): 要删除的IP地址
 *
 * Returns:
 */
func (p *Pool) DeleteIPFromPool(ip string) {
	// 删除原始IP从池中
	if ipDetails, ok := p.origin_ip_map.Get(ip); ok {
		p.origin_ip_map.Delete(ip)
		log.Debug("✅ 删除原始IP %s 成功", ip)

		// 从所有域名中删除
		p.domains.Range(func(key string, value *DomainInfo) bool {
			value.Delete(ipDetails)
			return true
		})
	} else {
		log.Debug("❌ 原始IP %s 不存在，无法删除", ip)
	}
}

/**
 * Function: GetSize
 *
 * Purpose: 获取IP池的大小
 *
 * Parameters:
 *
 * Returns: int64 - IP池的大小，表示可以容纳的最大域名数量
 */
func (p *Pool) GetSize() int64 {
	return p.size
}

/**
 * Function: GetDomains
 *
 * Purpose: 获取IP池中所有域名及其信息的映射
 *
 * Parameters:
 *
 * Returns: map[string]*DomainInfo - 域名到域名信息对象的映射
 */
func (p *Pool) GetDomains() map[string]*DomainInfo {
	domains := make(map[string]*DomainInfo)
	p.domains.Range(func(key string, value *DomainInfo) bool {
		domains[key] = value
		return true
	})
	return domains
}

/**
 * Function: GetDomainInfo
 *
 * Purpose: 获取指定域名的详细信息
 *
 * Parameters:
 *   - domain (string): 要查询的域名
 *
 * Returns:
 *   - *DomainInfo - 域名信息对象，如不存在则为nil
 *   - bool - 表示域名是否存在于池中
 */
func (p *Pool) GetDomainInfo(domain string) (*DomainInfo, bool) {
	domainInfo, ok := p.domains.Get(domain)
	if !ok {
		log.Debug("❌ 域名 %s 不存在", domain)
		return nil, false
	}
	log.Debug("✅ 获取域名 %s 的信息成功", domain)
	return domainInfo, true
}

/**
 * Function: AddDomain
 *
 * Purpose: 向IP池中添加新的域名
 *
 * Parameters:
 *   - domain (string): 要添加的域名
 *
 * Returns:
 */
func (p *Pool) AddDomain(domain string) {
	// 检查域名是否已存在
	ok := p.domains.Has(domain)
	if ok {
		log.Debug("❌ 域名 %s 已存在，无法重复添加", domain)
		return
	}

	// 创建新的域名信息对象并添加到池中
	domainInfo := NewDomainInfo(domain)
	p.domains.Set(domain, domainInfo)
	log.Debug("✅ 添加域名 %s 成功", domain)
}

/**
 * Function: DeleteDomain
 *
 * Purpose: 从IP池中删除指定域名及其所有IP信息
 *
 * Parameters:
 *   - domain (string): 要删除的域名
 *
 * Returns:
 */
func (p *Pool) DeleteDomain(domain string) {
	// 检查域名是否存在
	ok := p.domains.Has(domain)
	if !ok {
		log.Debug("❌ 域名 %s 不存在，无法删除", domain)
		return
	}

	// 删除域名信息对象
	p.domains.Delete(domain)
	log.Debug("✅ 删除域名 %s 成功", domain)
}

/**
 * Function: NewPool
 *
 * Purpose: 创建并初始化一个新的池对象
 *
 * Parameters:
 *   - size (int64): 池的大小
 *   - lifetime_per_ip (int64): 每个IP的生命周期(秒)
 *   - urlResource (string): 资源地址字符串
 *
 * Returns: *Pool - 初始化完成的Pool实例指针
 */
func NewPool(size int64, lifetime_per_ip int64) *Pool {
	return &Pool{
		size:          size,
		ip_ttl:        lifetime_per_ip,
		domains:       maputil.NewConcurrentMap[string, *DomainInfo](int(size)),
		origin_ip_map: maputil.NewConcurrentMap[string, *IpDetails](100),
	}
}
