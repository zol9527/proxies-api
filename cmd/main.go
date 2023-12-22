/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:29:39
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-22 15:38:22
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package main

import (
	"fmt"
	"ip-controller/config"
	"ip-controller/internal/http"
	"ip-controller/internal/pool"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	println("程序开始")

	go func() {
		fmt.Println("hhaha")
	}()

	// 依次启动
	go config.Start()
	go pool.PoolInit()
	go http.Server()

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
