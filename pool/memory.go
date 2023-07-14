/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-07-13 16:46:06
 * :LastEditors  :loyd
 * :LastEditTime :2023-07-14 17:58:18
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package pool

import (
	"sync"
)

// 格式: {'xxx:xxx:xxx:xxx': '2023-07-13 16:46:06'}
var MemoryPool sync.Map

// 格式:
//
//	sync.Map{
//		"www.xxx.com": sync.Map{
//			domain string
//			last_use time.Time
//			ip_pool sync.Map{
//				ip sync.Map{
//					acquire_time time.Time
//					last_use time.Time
//				}
//			}
//		}
//	}
var MemorySecondaryPool sync.Map
