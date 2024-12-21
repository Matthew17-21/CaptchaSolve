package captchasolve

import (
	"context"
	"testing"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (m *mockLogger) Log(message string)               {}
func (m *mockLogger) Debug(format string, args ...any) {}
func (m *mockLogger) Info(format string, args ...any)  {}
func (m *mockLogger) Warn(format string, args ...any)  {}
func (m *mockLogger) Error(format string, args ...any) {}

type mockHarvester struct{}

func (m *mockHarvester) Harvest() string { return "mock-token" }

func (m *mockHarvester) GetBalance() (float32, error) { return 0, nil }

func (m *mockHarvester) GetToken(additional ...*captchatoolsgo.AdditionalData) (*captchatoolsgo.CaptchaAnswer, error) {
	return nil, nil
} // Function to get a captcha token
func (m *mockHarvester) GetTokenWithContext(ctx context.Context, additional ...*captchatoolsgo.AdditionalData) (*captchatoolsgo.CaptchaAnswer, error) {
	return nil, nil
} // Function to get a captcha token

func TestWithMaxCapacity(t *testing.T) {
	cfg := &config{}
	option := WithMaxCapacity(10)
	option(cfg)

	assert.Equal(t, 10, cfg.maxCapacity, "maxCapacity should be set to 10")
}

func TestWithHarvester(t *testing.T) {
	cfg := &config{}
	mockHarvester := &mockHarvester{}
	option := WithHarvester(mockHarvester)
	option(cfg)

	assert.Contains(t, cfg.harvesters, mockHarvester, "harvester should be added to the harvesters slice")
}

func TestWithMaxGoroutines(t *testing.T) {
	cfg := &config{}

	t.Run("valid max goroutines", func(t *testing.T) {
		option := WithMaxGoroutines(5)
		option(cfg)
		assert.Equal(t, 5, cfg.maxCapacity, "maxCapacity should be set to 5")
	})

	t.Run("invalid max goroutines", func(t *testing.T) {
		option := WithMaxGoroutines(0)
		option(cfg)
		assert.Equal(t, 1, cfg.maxCapacity, "maxCapacity should be set to 1 when an invalid value is passed")
	})
}

func TestWithLogger(t *testing.T) {
	cfg := &config{}
	mockLogger := &mockLogger{}
	option := WithLogger(mockLogger)
	option(cfg)

	assert.Equal(t, mockLogger, cfg.logger, "logger should be set to the provided mockLogger")
}
