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
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
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

func (l *defaultLogger) Debug(args ...any) {
	l.logger.Println(append([]any{"[DEBUG]"}, args...)...)
}

func (l *defaultLogger) Debugf(format string, args ...any) {
	l.logger.Printf("[DEBUG] "+format, args...)
}

func (l *defaultLogger) Info(args ...any) {
	l.logger.Println(append([]any{"[INFO]"}, args...)...)
}

func (l *defaultLogger) Infof(format string, args ...any) {
	l.logger.Printf("[INFO] "+format, args...)
}

func (l *defaultLogger) Warn(args ...any) {
	l.logger.Println(append([]any{"[WARN]"}, args...)...)
}

func (l *defaultLogger) Warnf(format string, args ...any) {
	l.logger.Printf("[WARN] "+format, args...)
}

func (l *defaultLogger) Error(args ...any) {
	l.logger.Println(append([]any{"[ERROR]"}, args...)...)
}

func (l *defaultLogger) Errorf(format string, args ...any) {
	l.logger.Printf("[ERROR] "+format, args...)
}
