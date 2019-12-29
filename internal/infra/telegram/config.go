package telegram

import "time"

type Config struct {
	Token        string        `short:"t" long:"token" env:"TOKEN" description:"Telegram Bot API token" required:"yes"`
	Debug        bool          `long:"debug" env:"DEBUG" description:"Debug logs for Telegram Bot API adapter"`
	AutoDeleting time.Duration `long:"auto-deleting" env:"AUTO_DELETING" description:"Messages auto-deleting in supergroups" default:"1m"`
}
