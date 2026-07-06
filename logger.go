package telegram_api

import "log"

// 全局logger实例
var logger Logger = defaultLogger{}

// SetLogger 设置自定义logger
func SetLogger(l Logger) {
	if l != nil {
		logger = l
	}
}

// Logger 日志接口
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

// defaultLogger 默认日志实现
type defaultLogger struct{}

// Debug 输出Debug级别日志
func (l defaultLogger) Debug(args ...interface{}) { log.Println("[DEBUG]", args) }

// Info 输出Info级别日志
func (l defaultLogger) Info(args ...interface{}) { log.Println("[INFO]", args) }

// Warn 输出Warn级别日志
func (l defaultLogger) Warn(args ...interface{}) { log.Println("[WARN]", args) }

// Error 输出Error级别日志
func (l defaultLogger) Error(args ...interface{}) { log.Println("[ERROR]", args) }
