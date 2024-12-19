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
