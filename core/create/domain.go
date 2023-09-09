/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-09-08 17:49:43
 * :LastEditors  :loyd
 * :LastEditTime :2023-09-10 02:01:23
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package create

import (
	"errors"
	"ip-controller/log"
	"strconv"
	"sync"
	"time"
)

type Domain struct {
	Lock        sync.Mutex
	Domain      string
	LastUseTime time.Time
	Pool        []*IP
}

// -------------->> Function Split <<--------------
// Topic: Domain 中Pool的操作
// Description: 为了保证Pool的线程安全，需要加锁
// -------------->> Function Split <<--------------
func (d *Domain) GetIP() (*IP, error) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	// pop 出pool的第一个ip，然后放置于最后
	if len(d.Pool) == 0 {
		msg := "当前域名" + d.Domain + "没有可用ip"
		log.Debug(msg)
		return nil, errors.New(msg)
	} else {
		msg := "当前域名" + d.Domain + "已获取IP" + d.Pool[0].IP + "，可用IP数量" + strconv.Itoa(len(d.Pool)-1)
		ip := d.Pool[0]
		d.Pool = append(d.Pool[1:], ip)
		log.Debug(msg)
		return ip, nil
	}
}

func (d *Domain) AddIP(ip *IP) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Pool = append(d.Pool, ip)
}

func (d *Domain) DeleteIP(ip *IP) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	for i, v := range d.Pool {
		if v == ip {
			d.Pool = append(d.Pool[:i], d.Pool[i+1:]...)
		}
	}
}

// -------------->> Function Split <<--------------
// Topic: 创建 domain
// Description:
// -------------->> Function Split <<--------------
type DomainConfigInterface interface {
	SetDomain(domain string) DomainConfigInterface
	SetLastUseTime(last_use_time time.Time) DomainConfigInterface
	SetPool(pool []*IP) DomainConfigInterface
	Create() *Domain
}

type NewDomain struct {
	domain Domain
}

func (c *NewDomain) SetDomain(domain string) DomainConfigInterface {
	c.domain.Domain = domain
	return c
}

func (c *NewDomain) SetLastUseTime(last_use_time time.Time) DomainConfigInterface {
	c.domain.LastUseTime = last_use_time
	return c
}

func (c *NewDomain) SetPool(pool []*IP) DomainConfigInterface {
	c.domain.Pool = pool
	return c
}

func (c *NewDomain) Create() *Domain {
	return &c.domain
}
