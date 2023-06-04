/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-06-02 16:36:42
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-04 10:49:49
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

type NewIP struct {
	IP          string
	AcquireTime time.Time
}

type DomainIP struct {
	IP          string
	AcquireTime time.Time
	LastUse     time.Time
}

type NewDomainPool struct {
	Domain string
	IPPool []DomainIP
}

var (
	// IP 池
	Pool = sync.Map{}
	// Domain 池
	DomainPool = make(map[string][]DomainIP)
	// 动态读取的配置文件
	Settings *toml.Tree
	Locker   sync.Mutex
	SyncPipe = make(chan NewIP)
)
