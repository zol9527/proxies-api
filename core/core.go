/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 15:57:08
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-02 14:51:57
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type IP struct {
	IP    string
	Start time.Time
}

type DomainIP struct {
	IP      string
	Start   time.Time
	LastUse time.Time
}

type Domain struct {
	Domain string
	IPPool []DomainIP
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
	locker     = sync.Mutex{}
	Pool       = make([]IP, 0)
	DomainPool = sync.Map{}
)

/* 同步IP */
func syncIP(domainStr string) {

	// 从 Domain Pool 里面获取 Domain 对应的信息
	domain, ok := DomainPool.LoadAndDelete(domainStr)
	if !ok {
		// 复制 Pool里面的 IP 到 Domain 里面的 IPPool
		copied_pool := make([]DomainIP, 0)
		for _, ip := range Pool {
			copied_pool = append(copied_pool, DomainIP{ip.IP, ip.Start, time.Now().Add(-time.Duration(1000) * time.Second)})
		}
		DomainPool.Store(domainStr, Domain{domainStr, copied_pool})
		// 为防止上面 domain，虽然设置了，但是变量里面还没有设置, 这里重新获取
		domain, _ = DomainPool.Load(domainStr)
	}

	// 同步 主要 IP 到 Domain Pool 当中
	// Domain 池中的所有 IP 转换为 map
	domain_map := to_map(domain.(Domain).IPPool)

	// 获取新的 IP 准备添加到 Domain Pool 中
	tmp_pool := make([]DomainIP, 0)
	for _, ip := range Pool {
		if _, ok := domain_map[ip.IP]; !ok {
			tmp_pool = append(tmp_pool, DomainIP{ip.IP, ip.Start, time.Now()})
		}
	}
	// 添加新获取的 IP 到 Domain Pool 中
	new_total_pool := append(domain.(Domain).IPPool, tmp_pool...)
	DomainPool.Store(domainStr, Domain{domainStr, new_total_pool})
}

/* 选择合适的IP */
func selectIP(domainStr string) (string, error) {

	// 取出 Domain Pool 中的第一个 IP
	domain, _ := DomainPool.Load(domainStr)
	ip := domain.(Domain).IPPool[0]

	// 获取配置文件中 IP 最大存活时间
	exist_time := tree.Get("ip.exist_time").(int64)

	// 获取时间信息
	now := time.Now()
	start := ip.Start
	end := start.Add(time.Duration(exist_time) * time.Second)

	// IP 不在有效期内, 从 IP 池中删除
	if now.After(end) {
		new_domain_pool := append(domain.(Domain).IPPool[:0], domain.(Domain).IPPool[1:]...)
		DomainPool.Store(domainStr, Domain{domainStr, new_domain_pool})
		return "", errors.New("IP 不在有效期内")
	}

	// IP 在有效期内，使用 IP 前判断 last use 是否超过 pool.delay 设置
	delay := tree.Get("pool.delay").(int64)
	if now.After(ip.LastUse.Add(time.Duration(delay) * time.Second)) {
		// 超过 10 s, 更新 last use
		ip.LastUse = now
		// 更新到 Domain Pool 中
		new_domain_pool := append(domain.(Domain).IPPool[:0], domain.(Domain).IPPool[1:]...)
		new_domain_pool = append(new_domain_pool, ip)
		DomainPool.Store(domainStr, Domain{domainStr, new_domain_pool})
		return ip.IP, nil
		// 以 json 格式返回 IP
	} else {
		// w.Write([]byte("当前请求的 IP 不在有效期内, 请重新请求"))
		errMsg := fmt.Sprintf("当前请求的 IP 小于 %d 秒", delay)
		return "", errors.New(errMsg)
	}
}

/* 转换切片为 map */
func to_map(s []DomainIP) map[string]bool {
	m := make(map[string]bool)
	for _, ip := range s {
		m[ip.IP] = true
	}
	return m
}

/* 获取IP */
func core(w http.ResponseWriter, r *http.Request) {

	locker.Lock()
	defer locker.Unlock()

	// 获取客户端 Query "domain"
	domainStr := r.URL.Query().Get("domain")
	if domainStr == "" {
		w.Write([]byte("domain 不能为空"))
		return
	}

	// 检查池中是否有 IP
	if len(Pool) == 0 {
		w.Write([]byte("IP 池为空"))
		return
	}

	// 同步 IP 到 Domain Pool 中
	syncIP(domainStr)

	// 获取一个 IP
	ip, err := selectIP(domainStr)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// 返回 IP
	w.Header().Set("Content-Type", "application/json")
	response := IPResponse{Http: ip, Https: ip}
	jsonData, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(jsonData)
}

func Server() {

	// 获取 port
	port_str := tree.Get("server.port").(int64)
	port := strconv.FormatInt(port_str, 10)

	http.HandleFunc("/get", core)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
