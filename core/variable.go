/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-06-02 16:36:42
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-03 13:01:28
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"sync"
	"time"

	"github.com/pelletier/go-toml"
)

type IP struct {
	IP    string
	Start time.Time
}

type DomainIP struct {
	IP      string
	Start   time.Time
	LastUse time.Time
}

type Domain struct {
	Domain string
	IPPool []DomainIP
}

var (
	// IP 池
	Pool = sync.Map{}
	// Domain 池 {"domain": []DomainIP}
	DomainPool = sync.Map{}
	// 动态读取的配置文件
	Settings *toml.Tree
	Locker   sync.Mutex
)
