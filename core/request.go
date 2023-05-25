/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-24 14:48:07
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-25 17:50:44
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"strings"
	"sync"
	"time"

	"github.com/wangluozhe/requests"
)

/* 获取 IP  */
func RequestIP(stop_signal chan string, wg *sync.WaitGroup) {

	wg.Add(1)
loop_request_ip:

	for {
		url := ""
		r, err := requests.Get(url, nil)
		if err != nil {
			// todo 使用日志管理
			println("requests.Get error: ", err.Error())
			continue loop_request_ip
		}

		// 解析 IP
		ip := strings.TrimSpace(r.Text)

		// 添加进入 Pool
		locker.Lock()
		Pool = append(Pool, IP{
			Ip:     ip,
			Start:  time.Now(),
			Domain: []domain_info{},
		})
		locker.Unlock()

		// 检查是否有停止信号
		select {
		case _, ok := <-stop_signal:
			if !ok {
				break loop_request_ip
			}
		default:
		}

		// 休眠 1 秒
		time.Sleep(1 * time.Second)
	}

	wg.Done()
}
