/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-07-13 17:53:54
 * :LastEditors  :loyd
 * :LastEditTime :2023-07-14 17:19:17
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package init

import (
	"bytes"
	"fmt"
	"ip-controller/core"
	"ip-controller/log"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (m *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	newLog := fmt.Sprintf("[%s]: %s - %s\n", entry.Level, timestamp, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func initLog() {
	// 初始化设置
	log.Log = logrus.New()
	log.Log.SetOutput(os.Stdout)
	log.Log.SetFormatter(&CustomFormatter{})
	log.Log.SetLevel(logrus.DebugLevel)
	log.Log.SetReportCaller(true)

	// 输出日志
	log.Log.Debug("日志模块初始化完成")
}

func initSettings() {
	core.Settings, _ = toml.LoadFile("project.toml")

	log.Log.Debug("配置文件初始化完成")
}

func init() {
	initLog()
	initSettings()
}
