/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-30 17:56:55
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
	"sync"
	"time"
)

type domain_info struct {
	Domain  string
	LastUse time.Time
}

type IP struct {
	Ip     string
	Start  time.Time
	Domain []domain_info
}

type IPResponse struct {
	Http  string `json:"http"`
	Https string `json:"https"`
}

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

/* 定义重要参数 */
var (
	one    sync.Once
	locker sync.Mutex
	Pool   = make([]IP, 0)
)

/* 转换切片为 map */
func convert2map(s []domain_info) map[string]time.Time {
	m := make(map[string]time.Time)
	for _, v := range s {
		m[v.Domain] = v.LastUse
	}
	return m
}

/* 获取IP */
func core(w http.ResponseWriter, r *http.Request) {

	locker.Lock()
	defer locker.Unlock()

	// 获取客户端 Query "domain"
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		w.Write([]byte("domain 不能为空"))
		return
	}

	// 检查池中是否有 IP
	if len(Pool) == 0 {
		w.Write([]byte("IP 池为空"))
		return
	}

	// 取出 IP 池中的第一个 IP
	ip := Pool[0]

	// Exist time
	exist_time := tree.Get("ip.exist_time").(int64)

	// 获取时间信息
	now := time.Now()
	start := ip.Start
	end := start.Add(time.Duration(exist_time) * time.Second)

	// IP 不在有效期内, 从 IP 池中删除
	if now.After(end) {
		Pool = append(Pool[:0], Pool[1:]...)
		w.Write([]byte("当前请求的 IP 不在有效期内, 请重新请求"))
		return
	}

	// 获取 IP 里面的 Domain 切片
	// 转换 Pool slice 为 map
	pool_map := convert2map(ip.Domain)

	// 判断当前请求的 domain 是否在 IP 的 domain 切片中
	// 不存在则添加
	if v, ok := pool_map[domain]; !ok {
		ip.Domain = append(ip.Domain, domain_info{domain, now})
		Pool = append(Pool[:0], Pool[1:]...)
		Pool = append(Pool, ip)
		// 以 json 格式返回 IP
		w.Header().Set("Content-Type", "application/json")
		response := IPResponse{Http: ip.Ip, Https: ip.Ip}
		jsonData, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(jsonData)
		return
		// domain 存在, 判断是 domain info 中的 last use 是否超过 10 s
	} else {
		// 超过 10 s, 更新 last use
		if now.After(v.Add(10 * time.Second)) {
			ip.Domain = append(ip.Domain, domain_info{domain, now})
			Pool = append(Pool[:0], Pool[1:]...)
			Pool = append(Pool, ip)
			// 以 json 格式返回 IP
			w.Header().Set("Content-Type", "application/json")
			response := IPResponse{Http: ip.Ip, Https: ip.Ip}
			jsonData, err := json.Marshal(response)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(jsonData)
			return
			// 未超过 10 s, 让客户端重新请求
		} else {
			w.Write([]byte("当前请求的 IP 不在有效期内, 请重新请求"))
			return
		}
	}
}

func Server() {

	// 获取 port
	port_str := tree.Get("server.port").(int64)
	port := strconv.FormatInt(port_str, 10)

	http.HandleFunc("/get", core)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
