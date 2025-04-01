package models

import (
	"fmt"
	"testing"
	"time"
)

// Function: TestGetIPCount
//
// Purpose: 测试DomainInfo结构体的GetIPCount方法
//
// Parameters:
//   - t (*testing.T): 测试上下文对象
//
// Returns:
func TestGetIPCount(t *testing.T) {
	// 测试用例1: 空IP池
	t.Run("EmptyIPPool", func(t *testing.T) {
		// 创建一个空的DomainInfo实例
		domainInfo := NewDomainInfo("test.com")

		// 验证IP数量为0
		if count := domainInfo.GetIPCount(); count != 0 {
			t.Errorf("❌ 期望IP数量为0，但得到 %d", count)
		}
	})

	// 测试用例2: 含有多个IP的池
	t.Run("MultipleIPs", func(t *testing.T) {
		// 创建DomainInfo实例
		domainInfo := NewDomainInfo("example.com")

		// 添加3个IP
		ip1 := &IpDetails{
			ip:            "192.168.1.1",
			last_use_time: time.Now(),
		}
		ip2 := &IpDetails{
			ip:            "192.168.1.2",
			last_use_time: time.Now(),
		}
		ip3 := &IpDetails{
			ip:            "192.168.1.3",
			last_use_time: time.Now(),
		}

		domainInfo.Add(ip1)
		domainInfo.Add(ip2)
		domainInfo.Add(ip3)

		// 验证IP数量为3
		if count := domainInfo.GetIPCount(); count != 3 {
			t.Errorf("❌ 期望IP数量为3，但得到 %d", count)
		}

		// 删除一个IP后再次验证
		domainInfo.Delete(ip2)
		if count := domainInfo.GetIPCount(); count != 2 {
			t.Errorf("❌ 删除一个IP后期望数量为2，但得到 %d", count)
		}
	})

	// 测试用例3: 测试IP重复添加
	t.Run("DuplicateIPs", func(t *testing.T) {
		domainInfo := NewDomainInfo("duplicate.com")

		ip := &IpDetails{
			ip:            "192.168.1.1",
			last_use_time: time.Now(),
		}

		// 多次添加相同IP
		domainInfo.Add(ip)
		domainInfo.Add(ip)

		// OrderedMap会覆盖相同键的值，所以数量应该为1
		if count := domainInfo.GetIPCount(); count != 1 {
			t.Errorf("❌ 添加重复IP后期望数量为1，但得到 %d", count)
		}
	})
}

// Function: TestGet
//
// Purpose: 测试DomainInfo结构体的Get方法
//
// Parameters:
//   - t (*testing.T): 测试上下文对象
//
// Returns:
func TestGet(t *testing.T) {
	// 测试用例1: 空IP池
	t.Run("EmptyIPPool", func(t *testing.T) {
		// 创建一个空的DomainInfo实例
		domain_info := NewDomainInfo("empty.com")

		// 尝试获取IP
		ip, err := domain_info.Get()

		// 验证返回错误且IP为nil
		if ip != nil {
			t.Errorf("❌ 期望IP为nil，但得到 %v", ip)
		}
		if err == nil {
			t.Error("❌ 期望返回错误，但得到nil")
		}
		expected_err_msg := fmt.Sprintf("当前域名%s%s", "empty.com", ERROR_NO_IP_AVAILABLE)
		if err.Error() != expected_err_msg {
			t.Errorf("❌ 期望错误信息为 '%s'，但得到 '%s'", expected_err_msg, err.Error())
		}
	})

	// 测试用例2: 单个IP
	t.Run("SingleIP", func(t *testing.T) {
		// 创建DomainInfo实例
		domain_info := NewDomainInfo("single.com")

		// 添加一个IP
		old_time := time.Now().Add(-1 * time.Hour) // 1小时前
		ip_details := &IpDetails{
			ip:            "192.168.1.1",
			last_use_time: old_time,
		}
		domain_info.Add(ip_details)

		// 记录当前域名的最后使用时间
		old_domain_use_time := domain_info.GetLastUseTime()

		// 获取IP
		result_ip, err := domain_info.Get()

		// 验证获取的IP
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if result_ip == nil {
			t.Error("❌ 期望返回IP，但得到nil")
		}
		if result_ip.ip != "192.168.1.1" {
			t.Errorf("❌ 期望IP为 '192.168.1.1'，但得到 '%s'", result_ip.ip)
		}

		// 验证last_use_time已更新
		if !result_ip.last_use_time.After(old_time) {
			t.Errorf("❌ IP的最后使用时间应该更新")
		}

		// 验证域名的最后使用时间也已更新
		if !domain_info.GetLastUseTime().After(old_domain_use_time) {
			t.Error("❌ 域名的最后使用时间应该更新")
		}
	})

	// 测试用例3: 多个IP并验证轮换
	t.Run("MultipleIPs", func(t *testing.T) {
		// 创建DomainInfo实例
		domain_info := NewDomainInfo("multiple.com")

		// 添加三个IP
		ip1 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now()}
		ip2 := &IpDetails{ip: "192.168.1.2", last_use_time: time.Now()}
		ip3 := &IpDetails{ip: "192.168.1.3", last_use_time: time.Now()}

		domain_info.Add(ip1) // 此时顺序: [ip1]
		domain_info.Add(ip2) // 此时顺序: [ip1, ip2]
		domain_info.Add(ip3) // 此时顺序: [ip1, ip2, ip3]

		// 获取第一个IP
		first_ip, err := domain_info.Get()
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if first_ip.ip != "192.168.1.1" {
			t.Errorf("❌ 期望首个IP为 '192.168.1.1'，但得到 '%s'", first_ip.ip)
		}

		// 获取第二个IP
		second_ip, err := domain_info.Get()
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if second_ip.ip != "192.168.1.2" {
			t.Errorf("❌ 期望第二个IP为 '192.168.1.2'，但得到 '%s'", second_ip.ip)
		}

		// 获取第三个IP
		third_ip, err := domain_info.Get()
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if third_ip.ip != "192.168.1.3" {
			t.Errorf("❌ 期望第三个IP为 '192.168.1.3'，但得到 '%s'", third_ip.ip)
		}

		// 再次获取，应该回到第一个IP（轮换）
		fourth_ip, err := domain_info.Get()
		if err != nil {
			t.Errorf("❌ 没有期望错误，但得到 %v", err)
		}
		if fourth_ip.ip != "192.168.1.1" {
			t.Errorf("❌ 期望第四个IP为 '192.168.1.1'，但得到 '%s'", fourth_ip.ip)
		}
	})
}

// Function: TestAdd
//
// Purpose: 测试DomainInfo结构体的Add方法
//
// Parameters:
//   - t (*testing.T): 测试上下文对象
//
// Returns:
func TestAdd(t *testing.T) {
	// 测试用例1: 向空的DomainInfo添加IP
	t.Run("AddToEmpty", func(t *testing.T) {
		// 创建一个空的DomainInfo实例
		domainInfo := NewDomainInfo("test.com")

		// 初始IP数量应为0
		if count := domainInfo.GetIPCount(); count != 0 {
			t.Errorf("❌ 初始IP数量期望为0，但得到 %d", count)
		}

		// 添加一个IP
		ip := &IpDetails{
			ip:            "192.168.1.1",
			last_use_time: time.Now(),
		}
		domainInfo.Add(ip)

		// 验证IP已被添加
		if count := domainInfo.GetIPCount(); count != 1 {
			t.Errorf("❌ 添加IP后期望数量为1，但得到 %d", count)
		}

		// 验证可以获取到添加的IP
		allIPs := domainInfo.GetAllIPs()
		if len(allIPs) != 1 || allIPs[0].ip != "192.168.1.1" {
			t.Errorf("❌ 无法找到添加的IP '192.168.1.1'")
		}
	})

	// 测试用例2: 添加多个不同的IP
	t.Run("AddMultiple", func(t *testing.T) {
		// 创建DomainInfo实例
		domainInfo := NewDomainInfo("multiple.com")

		// 添加三个不同的IP
		ip1 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now()}
		ip2 := &IpDetails{ip: "192.168.1.2", last_use_time: time.Now()}
		ip3 := &IpDetails{ip: "192.168.1.3", last_use_time: time.Now()}

		domainInfo.Add(ip1)
		domainInfo.Add(ip2)
		domainInfo.Add(ip3)

		// 验证IP数量
		if count := domainInfo.GetIPCount(); count != 3 {
			t.Errorf("❌ 添加三个IP后期望数量为3，但得到 %d", count)
		}

		// 获取所有IP并验证是否都被正确添加
		allIPs := domainInfo.GetAllIPs()
		if len(allIPs) != 3 {
			t.Errorf("❌ GetAllIPs期望返回3个IP，但返回了 %d 个", len(allIPs))
		}

		// 创建一个map来检查所有IP是否存在
		ipMap := make(map[string]bool)
		for _, ip := range allIPs {
			ipMap[ip.ip] = true
		}

		// 验证所有添加的IP都在集合中
		if !ipMap["192.168.1.1"] || !ipMap["192.168.1.2"] || !ipMap["192.168.1.3"] {
			t.Errorf("❌ 无法找到全部添加的IP")
		}
	})

	// 测试用例3: 重复添加相同IP
	t.Run("AddDuplicate", func(t *testing.T) {
		// 创建DomainInfo实例
		domainInfo := NewDomainInfo("duplicate.com")

		// 添加一个IP
		ip1 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now()}
		domainInfo.Add(ip1)

		// 验证IP数量为1
		if count := domainInfo.GetIPCount(); count != 1 {
			t.Errorf("❌ 添加一个IP后期望数量为1，但得到 %d", count)
		}

		// 创建新的IpDetails但使用相同的IP地址
		ip2 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now().Add(1 * time.Hour)}
		domainInfo.Add(ip2)

		// 验证IP数量仍为1（因为使用了相同的键）
		if count := domainInfo.GetIPCount(); count != 1 {
			t.Errorf("❌ 添加重复IP后期望数量仍为1，但得到 %d", count)
		}

		// 验证存储的是最新添加的IP对象
		result, _ := domainInfo.Get()
		if result != ip2 {
			t.Errorf("❌ 期望存储的是最新添加的IP对象，但不是")
		}
	})

	// 测试用例4: 测试添加和删除后的IP位置顺序
	t.Run("AddDeleteOrder", func(t *testing.T) {
		// 创建DomainInfo实例
		domainInfo := NewDomainInfo("order.com")

		// 添加三个不同的IP
		ip1 := &IpDetails{ip: "192.168.1.1", last_use_time: time.Now()}
		ip2 := &IpDetails{ip: "192.168.1.2", last_use_time: time.Now()}
		ip3 := &IpDetails{ip: "192.168.1.3", last_use_time: time.Now()}

		domainInfo.Add(ip1)
		domainInfo.Add(ip2)
		domainInfo.Add(ip3)

		// 验证IP数量为3
		if count := domainInfo.GetIPCount(); count != 3 {
			t.Errorf("❌ 添加三个IP后期望数量为3，但得到 %d", count)
		}

		// 获取所有IP验证添加顺序
		allIPs := domainInfo.GetAllIPs()
		if len(allIPs) != 3 {
			t.Errorf("❌ 期望IP数量为3，但得到 %d", len(allIPs))
		}

		// 确认IP顺序为 ip1, ip2, ip3
		expectedOrder := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
		for i, ip := range allIPs {
			if ip.ip != expectedOrder[i] {
				t.Errorf("❌ 位置 %d 期望IP为 %s，但得到 %s", i, expectedOrder[i], ip.ip)
			}
		}

		// 删除中间的IP (ip2)
		domainInfo.Delete(ip2)

		// 验证IP数量为2
		if count := domainInfo.GetIPCount(); count != 2 {
			t.Errorf("❌ 删除一个IP后期望数量为2，但得到 %d", count)
		}

		// 获取所有IP验证删除后的顺序
		allIPs = domainInfo.GetAllIPs()

		// 确认IP顺序为 ip1, ip3
		expectedOrder = []string{"192.168.1.1", "192.168.1.3"}
		for i, ip := range allIPs {
			if ip.ip != expectedOrder[i] {
				t.Errorf("❌ 删除中间IP后，位置 %d 期望IP为 %s，但得到 %s", i, expectedOrder[i], ip.ip)
			}
		}

		// 添加新的IP
		ip4 := &IpDetails{ip: "192.168.1.4", last_use_time: time.Now()}
		domainInfo.Add(ip4)

		// 验证IP数量为3
		if count := domainInfo.GetIPCount(); count != 3 {
			t.Errorf("❌ 添加新IP后期望数量为3，但得到 %d", count)
		}

		// 获取所有IP验证添加后的顺序
		allIPs = domainInfo.GetAllIPs()

		// 确认IP顺序为 ip1, ip3, ip4
		expectedOrder = []string{"192.168.1.1", "192.168.1.3", "192.168.1.4"}
		for i, ip := range allIPs {
			if ip.ip != expectedOrder[i] {
				t.Errorf("❌ 添加新IP后，位置 %d 期望IP为 %s，但得到 %s", i, expectedOrder[i], ip.ip)
			}
		}

		// 验证轮换获取IP的顺序
		first, err := domainInfo.Get()
		if err != nil || first.ip != "192.168.1.1" {
			t.Errorf("❌ 期望获取第一个IP为 192.168.1.1，但得到 %v", first)
		}

		second, err := domainInfo.Get()
		if err != nil || second.ip != "192.168.1.3" {
			t.Errorf("❌ 期望获取第二个IP为 192.168.1.3，但得到 %v", second)
		}

		third, err := domainInfo.Get()
		if err != nil || third.ip != "192.168.1.4" {
			t.Errorf("❌ 期望获取第三个IP为 192.168.1.4，但得到 %v", third)
		}
	})

}
