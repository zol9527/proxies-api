/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-12 11:54:39
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-22 17:45:23
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import "time"

func Dispatcher() {

	// 插入新 IP 进入每个域
	for {
		select {
		case ip := <-IPPipeline:
			// 添加
			AddIP(ip)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
