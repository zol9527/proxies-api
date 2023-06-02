/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:48:07
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-01 16:56:00
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
		// 查看 POOL 里面的数量

		if int64(len(Pool)) >= tree.Get("pool.size").(int64) {
			time.Sleep(1 * time.Second)
			continue
		}

		// 请求代理网址
		r, err := requests.Get(tree.Get("proxy.url").(string), nil)
		if err != nil {
			log.Error("请求代理网址失败: ", err)
			continue loop_request_ip
		}

		// 解析 IP
		ip := strings.TrimSpace(r.Text)
		log.Debug("服务器返回: ", ip)
		// 使用正则表达式，查看是否是 IPv4
		if !IsIPv4(ip) {
			time.Sleep(1 * time.Second)
			continue loop_request_ip
		}
		total += 1
		log.Debugf("获取到 IP: %s, 累积数量: %d", ip, total)

		// 添加 IP 进入 Pool
		locker.Lock()
		Pool = append(Pool, IP{
			IP:    ip,
			Start: time.Now(),
		})
		locker.Unlock()

		// 休眠 1 秒
		time.Sleep(1 * time.Second)
	}

}
