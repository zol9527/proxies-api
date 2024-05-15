/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:29:39
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-15 17:16:59
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package main

import (
	"fmt"
	_ "ip-controller/config"
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
