package captchasolve

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
)

type mockHarvester struct {
	mock.Mock
}

func (m *mockHarvester) GetTokenWithContext(ctx context.Context, additional ...*captchatoolsgo.AdditionalData) (*captchatoolsgo.CaptchaAnswer, error) {
	args := m.Called(ctx, additional)
	return args.Get(0).(*captchatoolsgo.CaptchaAnswer), args.Error(1)
}

func TestStartHarvesters(t *testing.T) {
	ctx := context.Background()

	mockLogger := &mockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	mockHarvester := &mockHarvester{}
	mockHarvester.On("GetTokenWithContext", mock.Anything, mock.Anything).Return(&captchatoolsgo.CaptchaAnswer{}, nil)

	mockQueue := &mockQueue{}
	mockQueue.On("Enqueue", mock.Anything).Return(nil)

	c := &captchasolve{
		config: config{
			logger:        mockLogger,
			harvesters:    []captchatoolsgo.Harvester{mockHarvester},
			maxGoroutines: 1,
		},
		queue: mockQueue,
	}

	c.startHarvesters(ctx)

	mockLogger.AssertExpectations(t)
	mockHarvester.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

func TestHarvestToken_Success(t *testing.T) {
	ctx := context.Background()
	resultsChan := make(chan result, 1)

	mockLogger := &mockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	mockHarvester := &mockHarvester{}
	mockHarvester.On("GetTokenWithContext", mock.Anything, mock.Anything).Return(&captchatoolsgo.CaptchaAnswer{}, nil)

	c := &captchasolve{
		config: config{
			logger: mockLogger,
		},
	}

	c.harvestToken(ctx, mockHarvester, resultsChan)

	res := <-resultsChan
	assert.NoError(t, res.err)
	assert.NotNil(t, res.token)

	mockLogger.AssertExpectations(t)
	mockHarvester.AssertExpectations(t)
}

func TestProcessResults_Error(t *testing.T) {
	ctx := context.Background()
	resultsChan := make(chan result, 1)

	mockLogger := &mockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	mockQueue := &mockQueue{}

	c := &captchasolve{
		config: config{
			logger: mockLogger,
		},
		queue: mockQueue,
	}

	resultsChan <- result{token: nil, err: errors.New("test error")}
	close(resultsChan)

	_, err := c.processResults(ctx, resultsChan)
	assert.Error(t, err)

	mockLogger.AssertExpectations(t)
}
