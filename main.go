/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:29:39
 * :LastEditors  :loyd
 * :LastEditTime :2023-07-14 11:09:22
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package main

import (
	"ip-controller/core"
	_ "ip-controller/init"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	println("程序开始")

	// 运行获取 IP 的协程
	go core.RequestIP()
	// 运行IP HTTP Server
	go core.Server()
	go core.SyncIP()
	go core.SyncSettings()

	// 受到 signal 信号后，关闭通道
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
		println("等待其他程序结束")
	}()

	<-done
	println("程序结束")
}
