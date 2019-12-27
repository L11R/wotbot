package configs

import (
	"github.com/L11R/wotbot/internal/infra/database"
	"github.com/L11R/wotbot/internal/infra/xvm"
	"os"

	"github.com/L11R/wotbot/internal/infra/telegram"
	"github.com/L11R/wotbot/internal/infra/wargaming"
	"github.com/jessevdk/go-flags"
)

type Config struct {
	Database  *database.Config  `group:"Database args" namespace:"database" env-namespace:"WOT_DATABASE"`
	Telegram  *telegram.Config  `group:"Telegram args" namespace:"telegram" env-namespace:"WOT_TELEGRAM"`
	Wargaming *wargaming.Config `group:"Wargaming args" namespace:"wargaming" env-namespace:"WOT_WARGAMING"`
	XVM       *xvm.Config       `group:"XVM args" namespace:"xvm" env-namespace:"WOT_XVM"`

	Verbose []bool `short:"v" long:"verbose" env:"WOT_VERBOSE" description:"Verbose logs"`
}

func Parse() (*Config, error) {
	var config Config
	p := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)

	_, err := p.ParseArgs(os.Args)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
