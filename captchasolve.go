package captchasolve

import (
	"context"
	"errors"

	"github.com/Matthew17-21/CaptchaSolve/internal/queue"
)

type CaptchaSolve interface {
	GetToken(context.Context) (*CaptchaAnswer, error)
	ClearTokens() // Clears all pre-harvested tokens
}

type captchasolve struct {
	config
	queue *queue.SliceQueue[*CaptchaAnswer]
}

func New(opts ...ClientOption) CaptchaSolve {
	// Create default config
	cfg := defaultConfig()

	// Set any/all options
	for _, optFunc := range opts {
		optFunc(&cfg)
	}

	// Return the instance
	return &captchasolve{
		queue:  queue.NewSliceQueue[*CaptchaAnswer](),
		config: cfg,
	}
}

func (c *captchasolve) GetToken(ctx context.Context) (*CaptchaAnswer, error) {
	return nil, errors.New("not implemented")
	// TODO: Check if pre-harvested tokens are already saved
	// TODO: For every API key that still has a balance, attempt to get a captcha token
}

// ClearTokens removes any/all pre-harvested tokens
func (c *captchasolve) ClearTokens() { c.queue.Clear() }
