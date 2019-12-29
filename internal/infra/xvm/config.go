package xvm

import "time"

type Config struct {
	ChromeDevtoolsURL string        `long:"chrome-devtools-url" env:"CHROME_DEVTOOLS_URL" description:"Chrome Devtools URL" required:"yes"`
	HTTPTimeout       time.Duration `long:"http-timeout" env:"HTTP_TIMEOUT" description:"HTTP XVM webpage call timeout" default:"10s"`
	DevtoolsTimeout   time.Duration `long:"devtools-timeout" env:"DEVTOOLS_TIMEOUT" description:"Devtools XVM webpage call timeout" default:"10s"`
}
