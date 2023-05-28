/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-28 01:10:23
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-28 01:26:09
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"os"

	"github.com/sirupsen/logrus"
)

// 定义日志, 打印到 stdout
var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Out = os.Stdout
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.DebugLevel
}

// var log = &logrus.Logger{
// 	Out:       os.Stderr,
// 	Formatter: new(logrus.TextFormatter),
// 	Hooks:     make(logrus.LevelHooks),
// 	Level:     logrus.DebugLevel,
// }
