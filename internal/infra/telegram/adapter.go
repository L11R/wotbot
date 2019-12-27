package telegram

import (
	"github.com/L11R/wotbot/internal/domain"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type Adapter interface {
	ListenAndServe() error
	Shutdown()
}

type adapter struct {
	logger  *zap.Logger
	config  *Config
	botAPI  *tgbotapi.BotAPI
	service domain.Service
}

func NewAdapter(logger *zap.Logger, config *Config, service domain.Service) (Adapter, error) {
	a := &adapter{
		logger:  logger,
		config:  config,
		service: service,
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = config.Debug

	a.botAPI = bot

	return a, nil
}

func (a *adapter) ListenAndServe() error {
	a.logger.Info("Starting listening and serving Telegram Bot updates.")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	uu, err := a.botAPI.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for u := range uu {
		go a.route(&u)
	}

	return nil
}

func (a *adapter) Shutdown() {
	a.botAPI.StopReceivingUpdates()
}
