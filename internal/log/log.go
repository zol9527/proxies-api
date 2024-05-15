/*
 * :Author       :loyd zeng.yd@outlook.com
 * :Date         :2023-07-13 17:21:23
 * :LastEditors  :Please set LastEditors
 * :LastEditTime :2024-05-15 16:30:20
 * :Description  :
 *
 * Copyright (c) 2023 by loyd, All Rights Reserved.
 */

package log

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// 初始化设置
func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&CustomFormatter{})
	Log.SetLevel(logrus.DebugLevel)
	Log.SetReportCaller(true)
	Debug("日志模块初始化完成")
}

// 定义变量
var Log *logrus.Logger

// 设置打印格式
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

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	Log.Infoln(args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	Log.Debugln(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	Log.Warnln(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	Log.Errorln(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	Log.Fatalln(args...)
}

func Panic(args ...interface{}) {
	Log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	Log.Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	Log.Panicln(args...)
}
