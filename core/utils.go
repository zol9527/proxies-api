/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-06-02 17:05:47
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-03 13:00:30
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

/* 计算 IP 池大小 */
func countPool() int64 {
	pool_size := 0
	Pool.Range(func(key, value interface{}) bool {
		pool_size++
		return true
	})
	return int64(pool_size)
}

/* 转换切片为 map */
func to_map(s []DomainIP) map[string]bool {
	m := make(map[string]bool)
	for _, ip := range s {
		m[ip.IP] = true
	}
	return m
}
