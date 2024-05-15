/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-12-12 13:51:40
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-15 17:16:28
 * :Description  :动态更新配置文件
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package config

import (
	"ip-controller/internal/log"

	"github.com/pelletier/go-toml"
)

var Settings *toml.Tree

func init() {
	var err error // Declare the err variable
	Settings, err = toml.LoadFile("config/project.toml")
	if err != nil {
		log.Errorln("配置文件读取失败: ", err)
	}
	log.Debug("初始化配置文件完成")
}
