/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-07-13 16:46:06
 * :LastEditors  :loyd
 * :LastEditTime :2023-09-10 02:12:03
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import (
	"sync"
)

// -------------->> Function Split <<--------------
// Topic: 原始 IP
// Description: {'xxx:xxx:xxx:xxx': '2023-07-13 16:46:06'}
// -------------->> Function Split <<--------------
var OriginalPool sync.Map

// 获取原始 IP 池中 IP 数量
func GetOriginalPoolSize() int64 {
	size := 0
	OriginalPool.Range(func(k, v interface{}) bool {
		size++
		return true
	})
	return int64(size)
}

// -------------->> Function Split <<--------------
// Topic: IP 管理
// Description: 存储不同域名下的 IP
// -------------->> Function Split <<--------------
var MemoryPool sync.Map

// 获取池中 域名的数量
func GetDomainSize() int64 {
	size := 0
	MemoryPool.Range(func(k, v interface{}) bool {
		size++
		return true
	})
	return int64(size)
}
