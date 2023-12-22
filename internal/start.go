/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-22 11:50:32
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-22 15:14:41
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package internal

import (
	"ip-controller/internal/http"
	"ip-controller/internal/pool"
)

// 内部任务
type InternalTask struct {}

func (t *InternalTask) Start() {
	go pool.PoolInit()
	go http.Server()
}
