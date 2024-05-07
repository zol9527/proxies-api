/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-12 11:57:21
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-07 18:00:12
 * :Description  :POOL 部分启动
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

func PoolInit() {
	go AcquireIP()
	go Cleaner()
}
