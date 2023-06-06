/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-06 11:42:49
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

	// 确保 domain str 在 DomainPool 中
	d_pool, ok := DomainPool[domainStr]
	if !ok {
		DomainPool[domainStr] = make([]DomainIP, 0)
		w.Write([]byte("当前 IP 池为空"))
		return
	}

	// 确保 d_pool 中有一个以上的 IP
	if len(d_pool) == 0 {
		w.Write([]byte("当前 IP 池为空"))
		return
	}

	// 从 d_pool 中获取第一个 IP
	ip := d_pool[0]

	// 验证 ip 的 LastUse 是否超时
	delay := Settings.Get("pool.delay").(int64)
	// 计算当前 ip 是否超时
	now := time.Now()
	// 大于 delay, 可以继续使用
	if now.Sub(ip.LastUse).Seconds() > float64(delay) {
		// 调整 ip 位置
		ip.LastUse = now
		new_ips := append(d_pool[:0], d_pool[1:]...)
		new_ips = append(new_ips, ip)
		DomainPool[domainStr] = new_ips
		// 返回 IP
		w.Header().Set("Content-Type", "application/json")
		response := IPResponse{Http: ip.IP, Https: ip.IP}
		jsonData, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(jsonData)
		return
	} else {
		// 小于 delay, 返回错误
		// 调整位置后继续
		new_ips := append(d_pool[:0], d_pool[1:]...)
		new_ips = append(new_ips, ip)
		DomainPool[domainStr] = new_ips
		w.Write([]byte("当前 IP 未到可用时间"))
		return
	}
}

func Server() {

	// 获取 port
	port_str := Settings.Get("server.port").(int64)
	port := strconv.FormatInt(port_str, 10)

	http.HandleFunc("/get", core)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
