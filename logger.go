package telegram_api

import (
	"log"
	"os"
)

var logger Logger = &defaultLogger{
	logger: log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
}

func SetLogger(l Logger) {
	if l != nil {
		logger = l
	}
}

func GetLogger() Logger {
	return logger
}

// Logger 日志接口
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// defaultLogger 默认日志实现
type defaultLogger struct {
	logger *log.Logger
}

// NewDefaultLogger 创建默认日志器
func NewDefaultLogger() Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}
}

func (l *defaultLogger) Debug(args ...interface{}) {
	l.logger.Println(append([]interface{}{"[DEBUG]"}, args...)...)
}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	l.logger.Printf("[DEBUG] "+format, args...)
}

func (l *defaultLogger) Info(args ...interface{}) {
	l.logger.Println(append([]interface{}{"[INFO]"}, args...)...)
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf("[INFO] "+format, args...)
}

func (l *defaultLogger) Warn(args ...interface{}) {
	l.logger.Println(append([]interface{}{"[WARN]"}, args...)...)
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	l.logger.Printf("[WARN] "+format, args...)
}

func (l *defaultLogger) Error(args ...interface{}) {
	l.logger.Println(append([]interface{}{"[ERROR]"}, args...)...)
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+format, args...)
}
