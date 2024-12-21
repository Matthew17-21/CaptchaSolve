package captchasolve

import (
	"testing"

	"github.com/Matthew17-21/CaptchaSolve/internal/queue"
	"github.com/stretchr/testify/require"
)

func TestClearTokens(t *testing.T) {

	// Create new CaptchaSolve instance
	cs := captchasolve{queue: queue.NewSliceQueue[*CaptchaAnswer]()}

	// Push to queue
	const numElems int = 5
	for i := 0; i < numElems; i++ {
		cs.queue.Enqueue(&CaptchaAnswer{})
	}

	// Run method
	cs.ClearTokens()

	// Assert
	require.Empty(t, cs.queue.Len())
}
