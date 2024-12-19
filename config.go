package captchasolve

import captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"

const defaultMaxCapacity = 25

type config struct {
	maxCapacity int // The max amount of captcha tokens to be held
	harvesters  []captchatools.Harvester
}

func defaultConfig() config {
	return config{
		maxCapacity: defaultMaxCapacity,
		harvesters:  make([]captchatools.Harvester, 0),
	}
}
