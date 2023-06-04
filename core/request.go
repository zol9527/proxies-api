/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:48:07
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-03 15:51:03
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"net"
	"strings"
	"time"

	"github.com/wangluozhe/requests"
)

var (
	total int = 0
)

func IsIPv4(ip string) bool {
	// IP 可能是 ip:port 的形式，用:分割
	ip_v4 := strings.Split(ip, ":")[0]
	result := net.ParseIP(ip_v4)
	return result != nil
}

/* 获取 IP  */
func RequestIP() {

loop_request_ip:

	for {
		// 获取 IP 池中 IP 数量
		ip_pool_size := countPool()

		// 保证 IP 池的数量不超过最大值
		if ip_pool_size >= Settings.Get("pool.size").(int64) {
			time.Sleep(2 * time.Second)
			continue
		}

		// 请求代理网址
		r, err := requests.Get(Settings.Get("proxy.url").(string), nil)
		if err != nil {
			log.Error("请求代理网址失败: ", err)
			continue loop_request_ip
		}

		// 解析返回的数据
		ip := strings.TrimSpace(r.Text)
		log.Debug("服务器返回: ", ip)

		// 验证 IP 是否合法
		if !IsIPv4(ip) {
			time.Sleep(2 * time.Second)
			continue loop_request_ip
		}

		total += 1
		log.Debugf("获取到 IP: %s, 累积数量: %d", ip, total)

		// 添加 IP 到 IP 池
		now := time.Now()
		new_ip := NewIP{
			IP:          ip,
			AcquireTime: now,
		}
		SyncPipe <- new_ip

		// 休眠 1 秒
		time.Sleep(1 * time.Second)
	}

}
