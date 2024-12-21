package captchasolve

import captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"

const defaultMaxCapacity = 25

type config struct {
	// maxCapacity defines the maximum number of captcha tokens that can be held in the system.
	maxCapacity int

	// maxGoroutines specifies the maximum number of goroutines allowed to run concurrently
	// for solving captchas. This helps control resource usage and parallel processing.
	maxGoroutines int

	// harvesters is a slice of Harvester instances from the captchatools package.
	// These are used to fetch or generate captcha tokens as needed.
	harvesters []captchatools.Harvester

	// logger is an instance of the Logger interface used for logging system events,
	// debugging information, and error messages.
	logger Logger
}

func defaultConfig() config {
	return config{
		maxCapacity: defaultMaxCapacity,
		harvesters:  make([]captchatools.Harvester, 0),
		logger:      NewSilentLogger(),
	}
}
