package captchasolve

import (
	"time"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
)

// captchaTokenValidity defines the duration for which a captcha token remains valid.
// Once this duration has elapsed since solving, the token is considered expired.
const captchaTokenValidity = 2 * time.Minute

// CaptchaAnswer extends captchatools.CaptchaAnswer by adding metadata about when the captcha was solved.
type CaptchaAnswer struct {
	captchatoolsgo.CaptchaAnswer
	solvedAt time.Time // Timestamp indicating when the captcha was solved
}

// IsExpired checks whether the captcha token has expired based on its solve time.
// It compares the current time with the solve time plus the allowed validity duration.
func (c CaptchaAnswer) IsExpired() bool {
	expirationTime := c.solvedAt.Add(captchaTokenValidity)
	return time.Now().After(expirationTime)
}

// toCaptchaAnswer converts captchatoolsgo.CaptchaAnswer to CaptchaAnswer
func toCaptchaAnswer(c *captchatoolsgo.CaptchaAnswer) *CaptchaAnswer {
	return newCaptchaAnswer(c, time.Now())
}

// newCaptchaAnswer creates a new CaptchaAnswer instance with a specified solved timestamp.
//
// This function handles the actual conversion from the external to internal format,
// embedding the original answer and adding the solved timestamp. It safely handles
// nil inputs by returning an empty CaptchaAnswer struct.
func newCaptchaAnswer(c *captchatoolsgo.CaptchaAnswer, solvedAt time.Time) *CaptchaAnswer {
	if c == nil {
		return &CaptchaAnswer{}
	}
	return &CaptchaAnswer{
		*c,
		solvedAt,
	}
}
