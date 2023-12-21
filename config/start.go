/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-12 13:51:40
 * :LastEditors  :loyd
 * :LastEditTime :2023-12-12 14:00:31
 * :Description  :动态更新配置文件
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package config

import (
	"ip-controller/core/log"
	"time"

	"github.com/pelletier/go-toml"
)

var Settings *toml.Tree

func SettingsReload() {

	timer := time.NewTicker(5 * time.Second)
	
	for range timer.C {
		var err error // Declare the err variable
		Settings, err = toml.LoadFile("project.toml")
		if err != nil {
			log.Errorln("配置文件读取失败: ", err)
		}
		log.Debug("配置文件已更新")
	}

}
