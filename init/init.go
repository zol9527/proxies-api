/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-07-13 17:53:54
 * :LastEditors  :loyd
 * :LastEditTime :2023-09-08 17:43:46
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package init

import (
	"ip-controller/core"
	_ "ip-controller/log"

	"github.com/pelletier/go-toml"
)

func initSettings() {
	core.Settings, _ = toml.LoadFile("project.toml")
}

func init() {
	initSettings()
}
