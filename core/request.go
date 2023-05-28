/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:48:07
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-28 12:38:57
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
		if len(Pool) >= 120 {
			log.Debug("IP池已满")
			time.Sleep(1 * time.Second)
			continue
		}

		// 请求代理网址
		url := "http://dw.tiqu.letecs.com/getip_cm?num=1&type=1&pro=&city=0&yys=0&port=1&pack=238257&ts=0&ys=0&cs=0&lb=1&sb=0&pb=4&mr=1&regions=&code=aa604xkzbd"
		r, err := requests.Get(url, nil)
		if err != nil {
			log.Error("请求代理网址失败: ", err)
			continue loop_request_ip
		}

		// 解析 IP
		ip := strings.TrimSpace(r.Text)
		log.Debug("获取到 IP: ", ip)
		// 使用正则表达式，查看是否是 IPv4
		if !IsIPv4(ip) {
			continue loop_request_ip
		}

		// 添加 IP 进入 Pool
		locker.Lock()
		Pool = append(Pool, IP{
			Ip:     ip,
			Start:  time.Now(),
			Domain: []domain_info{},
		})
		locker.Unlock()

		// 休眠 1 秒
		time.Sleep(1 * time.Second)
	}

}
