package captchasolve

import (
	"testing"
	"time"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
	"github.com/stretchr/testify/require"
)

func TestIsExpired(t *testing.T) {
	tests := []struct {
		Id       string
		SolvedAt time.Time
		Expected bool
	}{
		{Id: "1", SolvedAt: time.Now(), Expected: false},
		{Id: "2", SolvedAt: time.Now().Add(time.Minute), Expected: false},
		{Id: "3", SolvedAt: time.Now().Add(time.Hour), Expected: false},
		{Id: "4", SolvedAt: time.Now().Add(-time.Minute), Expected: false},
		{Id: "5", SolvedAt: time.Now().Add(-captchaTokenValidity), Expected: true},
		{Id: "6", SolvedAt: time.Now().Add(-time.Hour), Expected: true},
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

func TestToCaptchaAnswer(t *testing.T) {
	// Setup
	input := &captchatoolsgo.CaptchaAnswer{
		Token:     "token",
		UserAgent: "ua",
	}

	// Test execution
	result := toCaptchaAnswer(input)

	// Assertions
	require.NotNil(t, result)
	require.Equal(t, input.Token, result.Token)
	require.Equal(t, input.UserAgent, result.UserAgent)

	// Verify that solvedAt time is recent
	timeDiff := time.Since(result.solvedAt)
	require.Less(t, timeDiff, 2*time.Second)
}

func TestNewCaptchaAnswer(t *testing.T) {
	// Setup
	input := &captchatoolsgo.CaptchaAnswer{
		Token:     "token",
		UserAgent: "ua",
	}
	fixedTime := time.Now()

	// Test execution
	result := newCaptchaAnswer(input, fixedTime)

	// Assertions
	require.NotNil(t, result)
	require.Equal(t, input.Token, result.Token)
	require.Equal(t, input.UserAgent, result.UserAgent)
	require.Equal(t, fixedTime, result.solvedAt)
}

func TestToCaptchaAnswerWithNilInput(t *testing.T) {
	// Test execution
	result := toCaptchaAnswer(nil)

	// Assertions
	require.NotNil(t, result)
	require.Empty(t, result.Id())
	require.Empty(t, result.Token)
	require.True(t, result.solvedAt.IsZero())
}

func TestNewCaptchaAnswerWithNilInput(t *testing.T) {

	// Setup & test execution
	result := newCaptchaAnswer(nil, time.Now())

	// Assertions
	require.NotNil(t, result)
	require.Empty(t, result.Id())
	require.Empty(t, result.Token)
	require.True(t, result.solvedAt.IsZero())
}
