/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:29:39
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-25 17:57:03
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package main

import (
	"ip-controller/core"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup

	// 定义通道
	stop_signal_channel := make(chan string)

	// 运行获取 IP 的协程
	go core.RequestIP(stop_signal_channel, &wg)
	// 运行IP HTTP Server
	go core.Server(stop_signal_channel, &wg)

	// 受到 signal 信号后，关闭通道
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		close(stop_signal_channel)
		done <- true
		println("等待其他程序结束")
		wg.Wait()
	}()

	<-done
	println("程序结束")
}
