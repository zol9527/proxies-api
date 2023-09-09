/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-09-08 17:47:26
 * :LastEditors  :loyd
 * :LastEditTime :2023-09-10 02:42:11
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package create

import "time"

type IP struct {
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	AuquireTime time.Time `json:"auquire_time"`
	LastuseTime time.Time `json:"expire_time"`
}

// IP interface
type IPConfigInterface interface {
	SetIP(ip string) IPConfigInterface
	SetPort(port int) IPConfigInterface
	SetAuquireTime(auquire_time time.Time) IPConfigInterface
	SetLastuseTime(expire_time time.Time) IPConfigInterface
	Create() *IP
}

type NewIP struct {
	ip IP
}

func (c *NewIP) SetIP(ip string) IPConfigInterface {
	c.ip.IP = ip
	return c
}

func (c *NewIP) SetPort(port int) IPConfigInterface {
	c.ip.Port = port
	return c
}

func (c *NewIP) SetAuquireTime(auquire_time time.Time) IPConfigInterface {
	c.ip.AuquireTime = auquire_time
	return c
}

func (c *NewIP) SetLastuseTime(expire_time time.Time) IPConfigInterface {
	c.ip.LastuseTime = expire_time
	return c
}

func (c *NewIP) Create() *IP {
	return &c.ip
}
