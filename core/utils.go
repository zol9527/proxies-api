/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-06-02 17:05:47
 * :LastEditors  :loyd
 * :LastEditTime :2023-09-10 02:14:04
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"net"
	"strings"
)

/* 判断是否是 ip */
func isIPv4(ip string) bool {
	// IP 可能是 ip:port 的形式，用:分割
	ip_v4 := strings.Split(ip, ":")[0]
	result := net.ParseIP(ip_v4)
	return result != nil
}
