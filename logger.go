package telegram_api

import "log"

var logger Logger = defaultLogger{}

func SetLogger(l Logger) {
	if l != nil {
		logger = l
	}
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type defaultLogger struct{}

func (l defaultLogger) Debug(args ...interface{}) { log.Println("[DEBUG]", args) }
func (l defaultLogger) Info(args ...interface{})  { log.Println("[INFO]", args) }
func (l defaultLogger) Warn(args ...interface{})  { log.Println("[WARN]", args) }
func (l defaultLogger) Error(args ...interface{}) { log.Println("[ERROR]", args) }
