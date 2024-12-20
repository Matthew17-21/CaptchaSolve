package captchasolve

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
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
	// Check if pre-harvested tokens are already saved
	if c.queue.Len() > 0 {
		// Get non-expired token. If token is expired, continue to harvest new one
		tkn, err := c.queue.Dequeue()
		if !errors.Is(err, queue.ErrQueueEmpty) {
			return nil, fmt.Errorf("error getting token from queue: %w", err)
		}
		// TODO: Handle queue.ErrQueueEmpty error
		if err == nil && !tkn.IsExpired() {
			return tkn, nil
		}
	}

	// Create channels for results and errors
	type result struct {
		token *CaptchaAnswer
		err   error
	}
	resultsChan := make(chan result, len(c.harvesters))

	// Use WaitGroup to track goroutines
	// TODO: Consider maxing the waitgroup a field in struct
	var wg sync.WaitGroup
	wg.Add(c.maxGoroutines) // TODO: Consider adding 1 at a time

	// Launch goroutine for each harvester
	for _, harvester := range c.harvesters {
		go func(h captchatoolsgo.Harvester) {
			// TODO: For every API key that still has a balance, attempt to get a captcha token
			tkn, err := h.GetTokenWithContext(ctx)
			if err != nil {
				resultsChan <- result{token: nil, err: fmt.Errorf("error getting token: %w", err)}
				return
			}
			resultsChan <- result{token: toCaptchaAnswer(tkn), err: nil}
		}(harvester)
	}

	// Close results channel after all goroutines complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	/*
		Return the first token that is sent to the results channel
		All other tokens that get sent to the results channel, add to queue
	*/
	var first *CaptchaAnswer
	for {
		select {
		case res, ok := <-resultsChan:
			if !ok {
				// Channel closed
				if first != nil {
					return first, nil
				}
				return nil, errors.New("no valid tokens found")
			}

			if res.err != nil {
				// Log errors and continue to process other results
				log.Println("Error on response:", res.err)
				continue
			}

			if res.token != nil {
				if first == nil {
					// Use the first valid token
					first = res.token
				} else {
					// Enqueue valid tokens that are not the first
					if err := c.queue.Enqueue(res.token); err != nil {
						return nil, fmt.Errorf("error enqueing token: %w", err)
					}
				}
			}
		case <-ctx.Done():
			// Context canceled
			return nil, ctx.Err()
		}
	}
}

// ClearTokens removes any/all pre-harvested tokens
func (c *captchasolve) ClearTokens() { c.queue.Clear() }
