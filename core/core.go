/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-25 17:35:17
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"net/http"
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

	// 遍历池中的每个 IP
	for index, ip := range Pool {

		// 获取时间信息
		now := time.Now()
		start := ip.Start
		end := start.Add(900 * time.Second)
		// IP 不在有效期内, 从 IP 池中删除
		if now.After(end) {
			Pool = append(Pool[:index], Pool[index+1:]...)
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
			Pool[index] = ip
			w.Write([]byte(ip.Ip))
			return
			// domain 存在, 判断是 domain info 中的 last use 是否超过 10 s
		} else {
			// 超过 10 s, 更新 last use
			if now.After(v.Add(10 * time.Second)) {
				ip.Domain = append(ip.Domain, domain_info{domain, now})
				Pool[index] = ip
				w.Write([]byte(ip.Ip))
				return
				// 未超过 10 s, 让客户端重新请求
			} else {
				w.Write([]byte("当前请求的 IP 不在有效期内, 请重新请求"))
				return
			}
		}
	}
}

func Server(stop_signal chan string, wg *sync.WaitGroup) {
	wg.Add(1)

loop_server:
	for {
		select {
		case _, ok := <-stop_signal:
			if !ok {
				break loop_server
			}
		default:
			one.Do(func() {
				http.HandleFunc("/get", core)
				http.ListenAndServe(":8080", nil)
			})
		}
	}
	wg.Done()
}
