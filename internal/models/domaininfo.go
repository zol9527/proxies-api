package models

import (
	"fmt"
	"proxies-api/internal/log"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
)

// 常量定义
const (
	ERROR_NO_IP_AVAILABLE = "没有可用的IP"
)

// DomainInfo 域名信息结构体
type DomainInfo struct {
	lock           sync.RWMutex                            // 读写锁，用于保护并发访问
	domain         string                                  // 域名
	last_use_time  time.Time                               // 上次使用时间
	delay          []int                                   // 延迟时间
	ip_details_map *maputil.OrderedMap[string, *IpDetails] // IP池，存储与域名相关的IP详情
}

// Function: NewDomainInfo
//
// Purpose: 创建并初始化一个新的域名信息对象
//
// Parameters:
//   - domain (string): 域名字符串
//
// Returns: *DomainInfo - 初始化完成的DomainInfo实例指针
func NewDomainInfo(domain string) *DomainInfo {
	return &DomainInfo{
		domain:         domain,
		last_use_time:  time.Now(),
		ip_details_map: maputil.NewOrderedMap[string, *IpDetails](),
		lock:           sync.RWMutex{},
	}
}

// Function: SetDelay
//
// Purpose: 设置域名的延迟时间
//
// Parameters:
//   - delay (int): 要设置的延迟时间（毫秒）
//
// Returns:
func (d *DomainInfo) SetDelay(delay []int) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 超过 2 个延迟时间则不设置
	if len(delay) > 2 {
		log.Error("❌ 域名 %s 设置延迟失败，延迟时间超过限制", d.domain)
		return false
	}
	// 设置延迟时间
	d.delay = delay
	log.Debug("⏱️ 域名 %s 设置延迟时间: %v", d.domain, delay)
	return true
}

// Function: GetDomain
//
// Purpose: 获取DomainInfo对象存储的域名
//
// Parameters:
//
// Returns: string - 域名字符串
func (d *DomainInfo) GetDomain() string {
	return d.domain
}

// Function: GetLastUseTime
//
// Purpose: 获取域名最后一次被使用的时间
//
// Parameters:
//
// Returns: time.Time - 域名最后使用的时间点
func (d *DomainInfo) GetLastUseTime() time.Time {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.last_use_time
}

// Function: GetIPCount
//
// Purpose: 获取域名下可用IP的数量
//
// Parameters:
//
// Returns: int - 可用IP的总数量
func (d *DomainInfo) GetIPCount() int {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.ip_details_map.Len()
}

// Function: Get
//
// Purpose: 从域名的IP池中获取一个IP并将其移至队列末尾
//
// Parameters:
//
// Returns:
//   - *IpDetails - 获取到的IP详情对象，如果没有可用IP则为nil
//   - error - 可能的错误，如无可用IP时返回错误信息
func (d *DomainInfo) Get() (*IpDetails, error) {
	// 使用读锁进行并发保护
	d.lock.RLock()
	defer d.lock.RUnlock()

	// 检查IP池是否为空
	if d.ip_details_map.Len() == 0 {
		log.Error("❌ 域名 %s 没有可用IP", d.domain)
		return nil, fmt.Errorf("当前域名%s%s", d.domain, ERROR_NO_IP_AVAILABLE)
	}

	// 获取第一个元素
	elem, ok := d.ip_details_map.Front()
	if !ok {
		log.Error("❌ 域名 %s 获取IP失败", d.domain)
		return nil, fmt.Errorf("当前域名%s%s", d.domain, ERROR_NO_IP_AVAILABLE)
	}

	// 移动到末尾（删除并重新添加）
	ip_details := elem.Value
	d.ip_details_map.Delete(elem.Key)
	d.ip_details_map.Set(ip_details.ip, ip_details)

	// 更新使用时间
	d.last_use_time = time.Now()
	ip_details.last_use_time = time.Now()

	log.Debug("🔄 域名 %s 获取IP: %s 成功", d.domain, ip_details.ip)
	return ip_details, nil
}

// Function: Add
//
// Purpose: 向域名的IP池中添加新的IP
//
// Parameters:
//   - ip (*IpDetails): 要添加的IP详情对象指针
//
// Returns:
func (d *DomainInfo) Add(ip *IpDetails) {
	// 使用写锁进行并发保护
	d.lock.Lock()
	defer d.lock.Unlock()

	// 添加IP到有序映射中
	d.ip_details_map.Set(ip.ip, ip)

	log.Debug("🚀 域名 %s 添加IP %s，当前可用IP数量: %d",
		d.domain,
		ip.ip,
		d.ip_details_map.Len())
}

// Function: Delete
//
// Purpose: 从域名的IP池中删除指定IP
//
// Parameters:
//   - ip (*IpDetails): 要删除的IP详情对象指针
//
// Returns:
func (d *DomainInfo) Delete(ip *IpDetails) {
	// 使用写锁进行并发保护
	d.lock.Lock()
	defer d.lock.Unlock()

	// 检查IP是否存在并删除
	if _, exists := d.ip_details_map.Get(ip.ip); exists {
		d.ip_details_map.Delete(ip.ip)
		log.Debug("✅ 域名 %s 删除IP %s，当前可用IP数量: %d",
			d.domain,
			ip.ip,
			d.ip_details_map.Len())
	} else {
		log.Debug("❌ 域名 %s 删除IP %s 失败，IP不存在",
			d.domain,
			ip.ip)
	}
}

// Function: DeleteByIP
//
// Purpose: 通过IP地址字符串从域名的IP池中删除指定IP
//
// Parameters:
//   - ipAddress (string): 要删除的IP地址字符串
//
// Returns: bool - 删除操作是否成功，成功返回true，失败返回false
func (d *DomainInfo) DeleteByIP(ipAddress string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	if _, exists := d.ip_details_map.Get(ipAddress); exists {
		d.ip_details_map.Delete(ipAddress)
		log.Debug("✅ 域名 %s 删除IP %s，当前可用IP数量: %d",
			d.domain,
			ipAddress,
			d.ip_details_map.Len())
		return true
	}

	log.Debug("❌ 域名 %s 删除IP %s 失败，IP不存在", d.domain, ipAddress)
	return false
}

// Function: GetAllIPs
//
// Purpose: 获取域名下所有的IP详情列表
//
// Parameters:
//
// Returns: []*IpDetails - IP详情对象指针的切片，包含所有IP信息
func (d *DomainInfo) GetAllIPs() []*IpDetails {
	d.lock.RLock()
	defer d.lock.RUnlock()

	result := make([]*IpDetails, 0, d.ip_details_map.Len())
	for elem := range d.ip_details_map.Iter() {
		result = append(result, elem.Value)
	}

	return result
}
