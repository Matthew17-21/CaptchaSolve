package captchasolve

import captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"

type config struct {
	maxCapacity int // The max amount of captcha tokens to be held
	harvesters  []captchatools.Harvester
}
