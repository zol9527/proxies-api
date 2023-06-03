/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-05-28 01:10:23
 * :LastEditors  :loyd
 * :LastEditTime :2023-06-03 12:43:52
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
	newLog := fmt.Sprintf("[%s]: %s - %s\n", entry.Level, timestamp, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

// 定义日志, 打印到 stdout
var log *logrus.Logger

func init() {
	// log = &logrus.Logger{
	// 	Out:       os.Stdout,
	// 	Formatter: &MyFormatter{},
	// 	Level:     logrus.DebugLevel,
	// }
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&MyFormatter{})
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
}
