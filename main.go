/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:29:39
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-25 17:35:49
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package main

import (
	"ip-controller/core"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// 定义通道
	trans_ip_channel := make(chan string)
	stop_signal_channel := make(chan string)

	// 运行获取 IP 的协程
	go core.RequestIP(trans_ip_channel, stop_signal_channel, &wg)
	// 运行IP HTTP Server
	go core.Server(stop_signal_channel, &wg)

	wg.Wait()
}
