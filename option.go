package captchasolve

import captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"

type ClientOption func(c *config)

// WithMaxCapacity sets the maximum number of tokens to be saved in the underlying data structure
func WithMaxCapacity(i int) ClientOption {
	return func(c *config) {
		c.maxCapacity = i
	}
}

// WithHarvester uses a given captcha harvester in the client
func WithHarvester(h captchatoolsgo.Harvester) ClientOption {
	return func(c *config) {
		c.harvesters = append(c.harvesters, h)
	}
}

// WithMaxGoroutines sets the max number of goroutines to be used while getting captcha tokens
func WithMaxGoroutines(max int) ClientOption {
	// Make sure it is a valid amount
	if max < 1 {
		max = 1
	}
	return func(c *config) {
		c.maxCapacity = max
	}
}

// WithLogger is a functional option for configuring a client with a custom logger.
// It accepts a Logger instance and returns a ClientOption function that sets the
// provided Logger in the client's configuration.
func WithLogger(l Logger) ClientOption {
	return func(c *config) {
		c.logger = l
	}
}
