// Author       :loyd
// Date         :2025-03-29 02:00:56
// LastEditors  :loyd
// LastEditTime :2025-03-30 19:14:09
// Description  :日志工具封装
//

package log

import (
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 常量定义
const (
	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"
)

// LogConfig 日志配置结构体
type LogConfig struct {
	LogLevel   string // 日志级别，可选值: debug, info, warn, error
	LogFile    string // 日志文件路径，为空则只输出到控制台
	MaxSize    int    // 单个日志文件最大大小(MB)，超过后会轮转
	MaxBackups int    // 最大备份日志文件数量
	MaxAge     int    // 日志文件保留天数
	Compress   bool   // 是否压缩旧日志文件
}

// 定义私有的全局logger实例
var log = logrus.New()

// 初始化日志
func init() {
	// 默认配置 - 在SetupLogger被显式调用前使用
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.SetReportCaller(true)
}

// SetupLogger 配置日志系统
func SetupLogger(config LogConfig) {
	// 创建日志目录
	if config.LogFile != "" {
		log_dir := path.Dir(config.LogFile)
		if _, err := os.Stat(log_dir); os.IsNotExist(err) {
			err := os.MkdirAll(log_dir, 0755)
			if err != nil {
				log.Fatalf("❌ 创建日志目录失败: %v", err)
			}
		}

		// 配置日志轮转
		lumberjack_logger := &lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    config.MaxSize, // 单位MB
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge, // 天数
			Compress:   config.Compress,
		}

		// 同时输出到控制台和文件
		multi_writer := io.MultiWriter(os.Stdout, lumberjack_logger)
		log.SetOutput(multi_writer)
	}

	// 设置日志格式
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", filename + ":" + strconv.Itoa(f.Line)
		},
	})

	// 设置日志级别
	switch config.LogLevel {
	case LOG_LEVEL_DEBUG:
		log.SetLevel(logrus.DebugLevel)
	case LOG_LEVEL_INFO:
		log.SetLevel(logrus.InfoLevel)
	case LOG_LEVEL_WARN:
		log.SetLevel(logrus.WarnLevel)
	case LOG_LEVEL_ERROR:
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	Info("🚀 日志系统初始化成功")
}

// 对外暴露的日志方法，直接调用

// Debug 输出调试日志
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info 输出信息日志
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn 输出警告日志
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error 输出错误日志
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatal 输出致命错误日志
func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// WithFields 创建带字段的日志条目
func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}
