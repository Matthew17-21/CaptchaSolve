package captchasolve

import (
	"errors"
	"fmt"

	"github.com/Matthew17-21/CaptchaSolve/internal/queue"
)

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
