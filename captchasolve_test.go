package captchasolve

import (
	"testing"
	"time"

	captchatoolsgo "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
	"github.com/stretchr/testify/require"
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
		answer := CaptchaAnswer{solvedAt: time.Now(), CaptchaAnswer: captchatoolsgo.CaptchaAnswer{Token: "123"}}
		queue.Enqueue(answer)

		result, err := queue.Dequeue()
		require.NoError(t, err)
		require.Equal(t, answer, result)
	})
}
