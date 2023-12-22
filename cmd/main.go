/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:29:39
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-22 15:15:57
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package main

import (
	"ip-controller/config"
	"ip-controller/internal"
	"os"
	"os/signal"
	"syscall"
	"sync"
)

// 依次启动
func StartChain() {

	type task interface{
		Start()
	}

	// 任务序列
	task_list := []task{
		&config.ConfigTask{},
	}

	// lock
	var l sync.Mutex

	// 遍历启动
	for _, t := range task_list {
		t.Start(l)
	}

}

func main() {

	println("程序开始")

	// 依次启动
	StartChain()

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
