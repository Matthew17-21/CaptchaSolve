package captchasolve

import (
	"context"
	"errors"
	"fmt"
	"sync"

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

// ClearTokens removes any/all pre-harvested tokens
func (c *captchasolve) ClearTokens() { c.queue.Clear() }

// getValidTokenFromQueue attempts to get a non-expired token from the queue
func (c *captchasolve) getValidTokenFromQueue() (*CaptchaAnswer, error) {
	// Check if pre-harvested tokens are already saved.
	// No need to check the length since the Dequeue method does it under the hood.
	c.logger.Info("Attempting to get a valid token from queue...")
	tkn, err := c.queue.Dequeue()
	if err != nil {
		if errors.Is(err, queue.ErrQueueEmpty) {
			c.logger.Info("Can't get token - queue is empty.")
			return nil, queue.ErrQueueEmpty
		}
		c.logger.Error("Unknown error getting token from queue:", err)
		return nil, fmt.Errorf("error dequeueing: %w", err)
	}
	if tkn.IsExpired() {
		c.logger.Info("Token is expired. Getting new one from queue...")
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
func (c *captchasolve) startHarvesters(ctx context.Context, additional ...*captchatoolsgo.AdditionalData) {
	// Create a results channel to collect harvester results
	resultsChan := make(chan result, len(c.harvesters))

	// Create a semaphore channel to limit concurrent goroutines
	sem := make(chan struct{}, c.maxGoroutines)

	// Create harvesters
	c.logger.Info("Creating %d harvesters...", len(c.harvesters))
	var wg sync.WaitGroup
	for i, harvester := range c.harvesters {
		wg.Add(1)
		c.logger.Info("Created harvester #%d", i+1)

		// Acquire semaphore
		sem <- struct{}{} // Will block if maxConcurrent goroutines are running

		go func(h captchatoolsgo.Harvester) {
			defer func() {
				<-sem // Release semaphore when done
				wg.Done()
			}()
			c.harvestToken(ctx, h, resultsChan, additional...)
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

func (c *captchasolve) harvestToken(ctx context.Context, h captchatoolsgo.Harvester, resultsChan chan<- result, additional ...*captchatoolsgo.AdditionalData) {
	c.logger.Info("Attempting to get a token from harvester...")
	tkn, err := h.GetTokenWithContext(ctx, additional...)
	if err != nil {
		c.logger.Error("Failed to get a token. Error: %v", err)
		resultsChan <- result{token: nil, err: fmt.Errorf("error getting token: %w", err)}
		return
	}
	c.logger.Info("Successfully got token with ID %v!", tkn.Id())
	resultsChan <- result{token: toCaptchaAnswer(tkn), err: nil}
}

func (c *captchasolve) processResults(ctx context.Context, resultsChan <-chan result) (*CaptchaAnswer, error) {
	/*
		Return the first token that is sent to the results channel
		All other tokens that get sent to the results channel, add to queue
	*/
	c.logger.Info("Processing results...")
	for {
		select {
		case res, ok := <-resultsChan:
			if !ok {
				return nil, errors.New("no valid tokens found")
			}

			if res.err != nil {
				c.logger.Error("Error on response:", res.err)
				continue
			}

			if res.token == nil {
				c.logger.Warn("error - token is nil. Retrying...")
				continue
			}

			// Add the token to queue
			if err := c.queue.Enqueue(res.token); err != nil {
				return nil, fmt.Errorf("error enqueuing token: %w", err)
			}

		case <-ctx.Done():
			c.logger.Warn("Context is done. Stopping...")
			return nil, ctx.Err()
		}
	}
}
