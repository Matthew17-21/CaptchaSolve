package captchasolve

import (
	"log"
)

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
}

type logger struct{}

func NewLogger() Logger {
	return &logger{}
}

func (logger) Debug(_ string, _ ...any) {}

func (logger) Info(format string, args ...any) {
	log.Printf(format+"\n", args...)
}

func (logger) Warn(format string, args ...any) {
	log.Printf(format+"\n", args...)
}

func (logger) Error(format string, args ...any) {
	log.Printf(format+"\n", args...)
}

type silentLogger struct {
}

func NewSilentLogger() Logger {
	return &silentLogger{}
}

func (s silentLogger) Debug(_ string, _ ...any) {}

func (s silentLogger) Info(_ string, _ ...any) {}

func (s silentLogger) Warn(_ string, _ ...any) {}

func (s silentLogger) Error(_ string, _ ...any) {}
