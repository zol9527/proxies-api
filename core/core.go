/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-03 13:29:13
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// type IP struct {
// 	IP    string
// 	Start time.Time
// }

// type DomainIP struct {
// 	IP      string
// 	Start   time.Time
// 	LastUse time.Time
// }

// type Domain struct {
// 	Domain string
// 	IPPool []DomainIP
// }

type IPResponse struct {
	Http  string `json:"http"`
	Https string `json:"https"`
}

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

/* 获取IP */
func core(w http.ResponseWriter, r *http.Request) {

	Locker.Lock()
	defer Locker.Unlock()

	// 获取客户端 Query "domain"
	domainStr := r.URL.Query().Get("domain")
	if domainStr == "" {
		w.Write([]byte("domain 不能为空"))
		return
	}

	// 获取 Domain pool
	d_pool, _ := DomainPool.LoadOrStore(domainStr, make([]DomainIP, 0))
	if len(d_pool.([]DomainIP)) == 0 {
		w.Write([]byte("当前 IP 池为空"))
		return
	}

	// 从 d_pool 中获取第一个 IP
	ip := d_pool.([]DomainIP)[0]

	// 验证 ip 的 LastUse 是否超时
	delay := Settings.Get("pool.delay").(int64)
	// 计算当前 ip 是否超时
	now := time.Now()
	// 大于 delay, 可以继续使用
	if now.Sub(ip.LastUse).Seconds() > float64(delay) {
		// 调整 ip 位置
		ip.LastUse = now
		new_ips := append(d_pool.([]DomainIP)[:0], d_pool.([]DomainIP)[1:]...)
		new_ips = append(new_ips, ip)
		DomainPool.Store(domainStr, new_ips)
		// 返回 IP
		w.Header().Set("Content-Type", "application/json")
		response := IPResponse{Http: ip.IP, Https: ip.IP}
		jsonData, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(jsonData)
	} else {
		// 小于 delay, 返回错误
		// 调整位置后继续
		new_ips := append(d_pool.([]DomainIP)[:0], d_pool.([]DomainIP)[1:]...)
		new_ips = append(new_ips, ip)
		DomainPool.Store(domainStr, new_ips)
		w.Write([]byte("当前 IP 未到可用时间"))
	}
}

func Server() {

	// 获取 port
	port_str := Settings.Get("server.port").(int64)
	port := strconv.FormatInt(port_str, 10)

	http.HandleFunc("/get", core)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
