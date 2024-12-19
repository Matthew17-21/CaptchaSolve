package captchasolve

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIsExpired(t *testing.T) {
	tests := []struct {
		Id       string
		SolvedAt time.Time
		Expected bool
	}{
		{Id: "1", SolvedAt: time.Now(), Expected: true},
		{Id: "2", SolvedAt: time.Now().Add(time.Minute), Expected: true},
		{Id: "3", SolvedAt: time.Now().Add(time.Hour), Expected: true},
		{Id: "4", SolvedAt: time.Now().Add(-time.Minute), Expected: true},
		{Id: "5", SolvedAt: time.Now().Add(-captchaTokenValidity), Expected: false},
		{Id: "6", SolvedAt: time.Now().Add(-time.Hour), Expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.Id, func(t *testing.T) {
			// Create captcha answer
			ca := CaptchaAnswer{solvedAt: tt.SolvedAt}

			// Run function
			result := ca.IsExpired()

			// Assert
			require.Equal(t, tt.Expected, result)
		})
	}
}
