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
	return time.Now().Before(expirationTime)
}