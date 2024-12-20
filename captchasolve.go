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
	// Attempt to get a token from queue
	token, err := c.getValidTokenFromQueue()
	if err == nil {
		return token, nil
	}

	// Start captcha harvesters
	go c.startHarvesters(ctx)

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

// ClearTokens removes any/all pre-harvested tokens
func (c *captchasolve) ClearTokens() { c.queue.Clear() }

// getValidTokenFromQueue attempts to get a non-expired token from the queue
func (c *captchasolve) getValidTokenFromQueue() (*CaptchaAnswer, error) {
	// Check if pre-harvested tokens are already saved.
	// No need to check the length since the Dequeue method does it under the hood.
	tkn, err := c.queue.Dequeue()
	if err != nil {
		if errors.Is(err, queue.ErrQueueEmpty) {
			return nil, queue.ErrQueueEmpty
		}
		return nil, fmt.Errorf("error dequeueing: %w", err)
	}
	if tkn.IsExpired() {
		return c.getValidTokenFromQueue()
	}
	return tkn, nil
}

// Create channels for results and errors
type result struct {
	token *CaptchaAnswer
	err   error
}

// startHarvesters coordinates concurrent token harvesting from multiple harvesters
func (c *captchasolve) startHarvesters(ctx context.Context) {
	// Create a results channel to collect harvester results
	resultsChan := make(chan result, len(c.harvesters))

	// Use a WaitGroup to manage goroutines
	// TODO: Consider maxing the waitgroup a field in struct
	// wg.Add(c.maxGoroutines) // TODO: Consider adding 1 at a time
	var wg sync.WaitGroup
	for _, harvester := range c.harvesters {
		wg.Add(1)
		go func(h captchatoolsgo.Harvester) {
			defer wg.Done()
			c.harvestToken(ctx, h, resultsChan)
		}(harvester)
	}

	// Close the results channel once all harvesters finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Process the results to retrieve the first valid token
	c.processResults(ctx, resultsChan)
}

func (c *captchasolve) harvestToken(ctx context.Context, h captchatoolsgo.Harvester, resultsChan chan<- result) {
	tkn, err := h.GetTokenWithContext(ctx)
	if err != nil {
		log.Printf("Failed to get a token. Error: %v. Sending to channel...\n", err)
		resultsChan <- result{token: nil, err: fmt.Errorf("error getting token: %w", err)}
		return
	}
	log.Printf("Successfully got token with ID %v! Sending to channel...\n", tkn.Id())
	resultsChan <- result{token: toCaptchaAnswer(tkn), err: nil}
}

func (c *captchasolve) processResults(ctx context.Context, resultsChan <-chan result) (*CaptchaAnswer, error) {
	/*
		Return the first token that is sent to the results channel
		All other tokens that get sent to the results channel, add to queue
	*/
	for {
		select {
		case res, ok := <-resultsChan:
			if !ok {
				return nil, errors.New("no valid tokens found")
			}

			if res.err != nil {
				log.Println("Error on response:", res.err)
				continue
			}

			if res.token != nil {
				log.Println("error - token is nil.")
				continue
			}

			// Add the token to queue
			if err := c.queue.Enqueue(res.token); err != nil {
				return nil, fmt.Errorf("error enqueuing token: %w", err)
			}

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
