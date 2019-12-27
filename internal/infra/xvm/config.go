package xvm

type Config struct {
	ChromeDevtoolsURL string `long:"chrome-devtools-url" env:"CHROME_DEVTOOLS_URL" description:"Chrome Devtools Websocket URL" required:"yes"`
}
