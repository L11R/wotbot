package telegram

import (
	"errors"
	"github.com/L11R/wotbot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (a *adapter) route(u *tgbotapi.Update) {
	if u.Message == nil { // ignore any non-Message Updates
		return
	}

	var (
		sentMsg *tgbotapi.Message
		err     error
	)

	defer func(err *error) {
		if r := recover(); r != nil {
			a.logger.Error("panic recoved!", zap.Any("panic", r))
			return
		}

		if err != nil && *err != nil {
			sentMsg = a.error(u, *err)
		}

		if u.Message.Chat.Type == "supergroup" {
			// Pass copies to goroutine
			if u != nil && sentMsg != nil {
				go func(update tgbotapi.Update, msg tgbotapi.Message) {
					ticker := time.NewTicker(10 * time.Second)
					<-ticker.C
					a.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)
					a.deleteMessage(msg.Chat.ID, msg.MessageID)
				}(*u, *sentMsg)
			}
		}
	}(&err)

	switch u.Message.Command() {
	case "start":
		sentMsg, err = a.handleStart(u)
	case "get":
		sentMsg, err = a.handleGet(u)
	case "save":
		sentMsg, err = a.handleSave(u)
	case "me":
		sentMsg, err = a.handleMe(u)
	case "refresh":
		sentMsg, err = a.handleRefresh(u)
	default:
		if strings.HasSuffix(u.Message.Command(), "Trend") ||
			strings.Contains(u.Message.Command(), "ByVehicle") {
			sentMsg, err = a.handleTrend(u)
		}
	}
}

func (a *adapter) handleStart(u *tgbotapi.Update) (*tgbotapi.Message, error) {
	text, err := a.service.GetCreateUserMessage(u.Message.From.ID)
	if err != nil {
		if errors.Is(err, domain.ErrInternalDatabase) {
			return nil, newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return nil, newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	sentMsg, err := a.botAPI.Send(msg)
	if err != nil {
		return nil, newHRError("Невозможно отправить сообщение!", err)
	}

	return &sentMsg, nil
}

func (a *adapter) handleGet(u *tgbotapi.Update) (*tgbotapi.Message, error) {
	if u.Message.CommandArguments() == "" {
		return nil, newHRError("Никнейм не передан!", domain.ErrBotBadRequest)
	}

	text, err := a.service.GetStatsMessage(u.Message.CommandArguments())
	if err != nil {
		if errors.Is(err, domain.ErrInternalWargaming) {
			return nil, newHRError("Ошибка при обращении к Wargaming API!", err)
		}
		if errors.Is(err, domain.ErrPlayerNotFound) {
			return nil, newHRError("Игрок с данным никнеймом не найден!", err)
		}
		if errors.Is(err, domain.ErrInternalXVM) {
			return nil, newHRError("Ошибка при обращении к XVM!", err)
		}

		return nil, newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	sentMsg, err := a.botAPI.Send(msg)
	if err != nil {
		return nil, newHRError("Невозможно отправить сообщение!", err)
	}

	return &sentMsg, nil
}

func (a *adapter) handleSave(u *tgbotapi.Update) (*tgbotapi.Message, error) {
	if u.Message.CommandArguments() == "" {
		return nil, newHRError("Никнейм не передан!", domain.ErrBotBadRequest)
	}

	text, err := a.service.GetSaveNicknameMessage(u.Message.From.ID, u.Message.CommandArguments())
	if err != nil {
		if errors.Is(err, domain.ErrInternalWargaming) {
			return nil, newHRError("Ошибка при обращении к Wargaming API!", err)
		}
		if errors.Is(err, domain.ErrPlayerNotFound) {
			return nil, newHRError("Игрок с данным никнеймом не найден!", err)
		}
		if errors.Is(err, domain.ErrInternalXVM) {
			return nil, newHRError("Ошибка при обращении к XVM!", err)
		}
		if errors.Is(err, domain.ErrInternalDatabase) {
			return nil, newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return nil, newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	sentMsg, err := a.botAPI.Send(msg)
	if err != nil {
		return nil, newHRError("Невозможно отправить сообщение!", err)
	}

	return &sentMsg, nil
}

func (a *adapter) handleRefresh(u *tgbotapi.Update) (*tgbotapi.Message, error) {
	text, err := a.service.GetRefreshMessage(u.Message.From.ID)
	if err != nil {
		if errors.Is(err, domain.ErrInternalXVM) {
			return nil, newHRError("Ошибка при обращении к XVM!", err)
		}
		if errors.Is(err, domain.ErrNicknameNotSaved) {
			return nil, newHRError("Сначала сохрани свой никнейм!", err)
		}
		if errors.Is(err, domain.ErrInternalDatabase) {
			return nil, newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return nil, newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	sentMsg, err := a.botAPI.Send(msg)
	if err != nil {
		return nil, newHRError("Невозможно отправить сообщение!", err)
	}

	return &sentMsg, nil
}

func (a *adapter) handleMe(u *tgbotapi.Update) (*tgbotapi.Message, error) {
	text, err := a.service.GetMeMessage(u.Message.From.ID)
	if err != nil {
		if errors.Is(err, domain.ErrTrendImageNotFound) {
			return nil, newHRError("График не найден!", err)
		}
		if errors.Is(err, domain.ErrInternalDatabase) {
			return nil, newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return nil, newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	sentMsg, err := a.botAPI.Send(msg)
	if err != nil {
		return nil, newHRError("Невозможно отправить сообщение!", err)
	}

	return &sentMsg, nil
}

func (a *adapter) handleTrend(u *tgbotapi.Update) (*tgbotapi.Message, error) {
	img, err := a.service.GetTrendImage(u.Message.From.ID, "#"+u.Message.Command())
	if err != nil {
		if errors.Is(err, domain.ErrNicknameNotSaved) {
			return nil, newHRError("Сначала сохрани свой никнейм!", err)
		}
		if errors.Is(err, domain.ErrInternalDatabase) {
			return nil, newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return nil, newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewPhotoUpload(u.Message.Chat.ID, tgbotapi.FileBytes{
		Name:  u.Message.Command(),
		Bytes: img,
	})
	sentMsg, err := a.botAPI.Send(msg)
	if err != nil {
		return nil, newHRError("Невозможно отправить сообщение!", err)
	}

	return &sentMsg, nil
}

func (a *adapter) error(update *tgbotapi.Update, err error) *tgbotapi.Message {
	if update == nil || err == nil {
		// Why did you call this function?
		return nil
	}

	// Log error
	a.logger.Error("Error occurred in handler!", zap.Error(err))

	// Send human readable representation of error to user to let him know
	if hrerr, ok := err.(*hrError); ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, hrerr.Human())
		sentMsg, err := a.botAPI.Send(msg)
		if err != nil {
			a.logger.Error("Error sending message with human readable error!", zap.Error(err))
			return nil
		}

		return &sentMsg
	}

	return nil
}

func (a *adapter) deleteMessage(chatID int64, messageID int) {
	if _, err := a.botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}); err != nil {
		a.logger.Error(
			"Error deleting the message!",
			zap.Int64("chat_id", chatID),
			zap.Int("message_id", messageID),
			zap.Error(err),
		)
	}
}
