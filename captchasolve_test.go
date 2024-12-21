package captchasolve

import (
	"context"
	"errors"
	"testing"
	"time"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/mock"
)

func TestNew(t *testing.T) {
	t.Run("creates default instance", func(t *testing.T) {
		// Act
		solver := New()

		// Assert
		require.NotNil(t, solver)
		require.NotNil(t, solver.(*captchasolve).queue)
		require.Empty(t, solver.(*captchasolve).config.harvesters)
	})

	t.Run("applies single option", func(t *testing.T) {
		// Arrange
		expectedMaxCapacity := 1

		// Act
		solver := New(WithMaxCapacity(expectedMaxCapacity))

		// Assert
		require.Equal(t, expectedMaxCapacity, solver.(*captchasolve).config.maxCapacity)
	})

	t.Run("applies multiple options in order", func(t *testing.T) {
		// Arrange
		var harvester captchatoolsgo.Harvester
		expectedMaxCapacity := 30

		// Act
		solver := New(WithMaxCapacity(expectedMaxCapacity), WithHarvester(harvester))

		// Assert
		cfg := solver.(*captchasolve).config
		require.NotEmpty(t, cfg.harvesters)
		require.Equal(t, harvester, cfg.harvesters[0])
		require.Equal(t, expectedMaxCapacity, cfg.maxCapacity)
	})

	t.Run("later options override earlier ones", func(t *testing.T) {
		// Arrange
		firstMaxCap := defaultMaxCapacity
		secondMaxCap := defaultMaxCapacity + 1

		// Act
		solver := New(WithMaxCapacity(firstMaxCap), WithMaxCapacity(secondMaxCap))

		// Assert
		require.Equal(t, secondMaxCap, solver.(*captchasolve).config.maxCapacity)
	})

	t.Run("queue is properly initialized", func(t *testing.T) {
		// Act
		solver := New()

		// Assert
		queue := solver.(*captchasolve).queue
		require.NotNil(t, queue)

		// Verify queue operations work
		answer := &CaptchaAnswer{solvedAt: time.Now(), CaptchaAnswer: captchatoolsgo.CaptchaAnswer{Token: "123"}}
		queue.Enqueue(answer)

		result, err := queue.Dequeue()
		require.NoError(t, err)
		require.Equal(t, answer, result)
	})
}

// MockTokenQueue is a mock implementation of TokenQueue
type MockTokenQueue struct {
	mock.Mock
}

func (m *MockTokenQueue) Enqueue(token *CaptchaAnswer) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenQueue) Dequeue() (*CaptchaAnswer, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CaptchaAnswer), args.Error(1)
}

func (m *MockTokenQueue) Clear() {
	m.Called()
}

func (m *MockTokenQueue) Len() int {
	return 0
}

// MockHarvester is a mock implementation of captchatoolsgo.Harvester
type MockHarvester struct {
	mock.Mock
}

func (m *MockHarvester) GetTokenWithContext(ctx context.Context, additional ...*captchatoolsgo.AdditionalData) (*captchatoolsgo.CaptchaAnswer, error) {
	args := m.Called(ctx, additional)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*captchatoolsgo.CaptchaAnswer), args.Error(1)
}

func TestGetToken(t *testing.T) {
	t.Run("returns valid token from queue immediately", func(t *testing.T) {
		// Arrange
		mockQueue := new(MockTokenQueue)
		expectedToken := &CaptchaAnswer{
			CaptchaAnswer: captchatoolsgo.CaptchaAnswer{Token: "valid-token"},
			solvedAt:      time.Now(),
		}
		mockQueue.On("Dequeue").Return(expectedToken, nil)

		solver := &captchasolve{
			queue: mockQueue,
		}
		solver.logger = NewSilentLogger()

		// Act
		token, err := solver.GetToken(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
		mockQueue.AssertExpectations(t)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		// Arrange
		mockQueue := new(MockTokenQueue)
		mockQueue.On("Dequeue").Return(nil, errors.New("queue empty"))

		solver := &captchasolve{
			queue: mockQueue,
		}
		solver.logger = NewSilentLogger()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Act
		token, err := solver.GetToken(ctx)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled))
		assert.Nil(t, token)
		mockQueue.AssertExpectations(t)
	})

	t.Run("retries queue until valid token is available", func(t *testing.T) {
		// Arrange
		mockQueue := new(MockTokenQueue)
		expectedToken := &CaptchaAnswer{
			CaptchaAnswer: captchatoolsgo.CaptchaAnswer{Token: "valid-token"},
			solvedAt:      time.Now(),
		}

		// First call returns empty, second call returns valid token
		mockQueue.On("Dequeue").Return(nil, errors.New("queue empty")).Once()
		mockQueue.On("Dequeue").Return(expectedToken, nil).Once()

		solver := &captchasolve{
			queue: mockQueue,
		}
		solver.logger = NewSilentLogger()

		// Act
		token, err := solver.GetToken(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
		mockQueue.AssertExpectations(t)
	})

	t.Run("starts harvesters when queue is empty", func(t *testing.T) {
		// Arrange
		mockQueue := new(MockTokenQueue)
		mockHarvester := new(MockHarvester)
		expectedToken := &CaptchaAnswer{
			CaptchaAnswer: captchatoolsgo.CaptchaAnswer{Token: "harvested-token"},
			solvedAt:      time.Now(),
		}

		// Queue is initially empty, then gets token from harvester
		mockQueue.On("Dequeue").Return(nil, errors.New("queue empty")).Once()
		mockQueue.On("Dequeue").Return(expectedToken, nil).Once()

		mockHarvester.On("GetTokenWithContext", mock.Anything, mock.Anything).Return(
			&captchatoolsgo.CaptchaAnswer{Token: "harvested-token"},
			nil,
		)

		solver := &captchasolve{
			queue: mockQueue,
		}
		solver.harvesters = make([]captchatoolsgo.Harvester, 0)
		solver.logger = NewSilentLogger()

		// Act
		token, err := solver.GetToken(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
		mockQueue.AssertExpectations(t)
	})

	t.Run("handles nil additional data", func(t *testing.T) {
		// Arrange
		mockQueue := new(MockTokenQueue)
		expectedToken := &CaptchaAnswer{
			CaptchaAnswer: captchatoolsgo.CaptchaAnswer{Token: "valid-token"},
			solvedAt:      time.Now(),
		}
		mockQueue.On("Dequeue").Return(expectedToken, nil)

		solver := &captchasolve{
			queue: mockQueue,
		}
		solver.logger = NewSilentLogger()

		// Act
		token, err := solver.GetToken(context.Background(), nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
		mockQueue.AssertExpectations(t)
	})
}
