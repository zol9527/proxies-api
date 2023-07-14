/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-29 11:42:37
 * :LastEditors  :loyd
 * :LastEditTime :2023-07-14 14:19:46
 * :Description  :读取 toml 配置文件
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"github.com/pelletier/go-toml"
)

var Settings *toml.Tree
