/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-29 11:42:37
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-29 11:43:41
 * :Description  :读取 toml 配置文件
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import "github.com/pelletier/go-toml"

var (
	tree *toml.Tree
)

func init() {
	tree, _ = toml.LoadFile("project.toml")
}
