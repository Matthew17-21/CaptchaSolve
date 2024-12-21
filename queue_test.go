package captchasolve

import (
	"testing"

	"github.com/Matthew17-21/CaptchaSolve/internal/queue"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/mock"
)

type mockQueue struct {
	mock.Mock
}

func (m *mockQueue) Enqueue(token *CaptchaAnswer) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockQueue) Dequeue() (*CaptchaAnswer, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CaptchaAnswer), args.Error(1)
}

func (m *mockQueue) Clear() {
	m.Called()
}

func (m *mockQueue) Len() int {
	args := m.Called()
	return args.Int(0) // Retrieve the first return value as an int
}

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
