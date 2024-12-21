package captchasolve

import (
	"context"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
	"github.com/Matthew17-21/CaptchaSolve/internal/queue"
)

type CaptchaSolve interface {
	// GetToken retrieves a valid captcha token using the configured harvesters.
	// It first attempts to return a pre-harvested token from the queue, and if none
	// are available, starts harvesting new tokens.
	//
	// The context parameter can be used to cancel the token retrieval operation.
	// Additional data can be provided if required by the harvesting service.
	//
	// Returns a valid CaptchaAnswer and nil error if successful, or nil and an error
	// if token retrieval fails or is cancelled.
	GetToken(context.Context, ...*captchatoolsgo.AdditionalData) (*CaptchaAnswer, error)

	// ClearTokens removes all pre-harvested tokens from the internal queue.
	// This is useful when you want to ensure fresh tokens are retrieved on
	// subsequent GetToken calls or when you need to clear potentially stale tokens.
	ClearTokens() // Clears all pre-harvested tokens
}

type captchasolve struct {
	config
	queue *queue.SliceQueue[*CaptchaAnswer]
}

// New initializes a CaptchaSolve with default configuration and then applies any provided
// option functions to customize the configuration. It creates an empty token queue and
// returns the fully initialized instance ready for use.
//
// Example:
//
//	solver := New(
//	    WithMaxGoroutines(5),
//	    WithLogger(customLogger),
//	)
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

// GetToken retrieves a valid captcha token, either from the pre-harvested queue or by
// starting new harvesters if needed.
//
// The function first attempts to get a pre-harvested token from the queue. If none are
// available, it starts background harvesters to generate new tokens. It continuously
// checks the queue for new tokens until either:
//   - A valid token is found
//   - The context is cancelled
//
// The function is non-blocking on harvester initialization, allowing multiple concurrent
// calls to GetToken. Harvesters run in the background and add tokens to the queue as
// they become available.
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
