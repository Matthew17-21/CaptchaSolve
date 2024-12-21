package captchasolve

import "github.com/test-go/testify/mock"

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *mockLogger) Error(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *mockLogger) Warn(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *mockLogger) Debug(format string, args ...interface{}) {
	m.Called(format, args)
}
