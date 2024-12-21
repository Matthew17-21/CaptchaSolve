package captchasolve

import (
	"context"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
	"github.com/Matthew17-21/CaptchaSolve/internal/queue"
)

type CaptchaSolve interface {
	GetToken(context.Context, ...*captchatoolsgo.AdditionalData) (*CaptchaAnswer, error)
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

func (c *captchasolve) GetToken(ctx context.Context, additional ...*captchatoolsgo.AdditionalData) (*CaptchaAnswer, error) {
	// Attempt to get a token from queue
	token, err := c.getValidTokenFromQueue()
	if err == nil {
		return token, nil
	}

	// Start captcha harvesters
	go c.startHarvesters(ctx, additional...)

	// While ctx not cancelled, return first token from queue
	for {
		// Make sure ctx is not cancelled
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default: // So it doesn;t block
		}

		// Return the first token from queue
		token, err := c.getValidTokenFromQueue()
		if err == nil {
			return token, nil
		}
	}
}

// Create channels for results and errors
type result struct {
	token *CaptchaAnswer
	err   error
}
