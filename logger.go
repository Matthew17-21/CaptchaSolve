package captchasolve

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
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
