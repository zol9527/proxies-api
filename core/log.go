/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-28 01:10:23
 * :LastEditors  :loyd
 * :LastEditTime :2023-05-29 12:00:31
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package core

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	newLog := fmt.Sprintf("[%s]: %s - [%s:%d] %s\n", entry.Level, timestamp, entry.Caller.File, entry.Caller.Line, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

// 定义日志, 打印到 stdout
var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Out = os.Stdout
	log.Formatter = new(MyFormatter)
	log.Level = logrus.DebugLevel
}
