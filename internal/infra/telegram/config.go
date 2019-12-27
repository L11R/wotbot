package telegram

type Config struct {
	Token string `short:"t" long:"token" env:"TOKEN" description:"Telegram Bot API token" required:"yes"`
	Debug bool   `long:"debug" env:"DEBUG" description:"Debug logs for Telegram Bot API adapter"`
}
