package telegram

import (
	"errors"
	"strings"

	"github.com/L11R/wotbot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func (a *adapter) route(u *tgbotapi.Update) {
	if u.Message == nil { // ignore any non-Message Updates
		return
	}

	switch u.Message.Command() {
	case "start":
		if err := a.handleStart(u); err != nil {
			a.error(u, err)
		}
	case "get":
		if err := a.handleGet(u); err != nil {
			a.error(u, err)
		}
	case "save":
		if err := a.handleSave(u); err != nil {
			a.error(u, err)
		}
	case "me":
		if err := a.handleMe(u); err != nil {
			a.error(u, err)
		}
	case "refresh":
		if err := a.handleRefresh(u); err != nil {
			a.error(u, err)
		}
	default:
		if strings.HasSuffix(u.Message.Command(), "Trend") {
			if err := a.handleTrend(u); err != nil {

			}
		}
	}
}

func (a *adapter) handleStart(u *tgbotapi.Update) error {
	text, err := a.service.GetCreateUserMessage(u.Message.From.ID)
	if err != nil {
		if errors.Is(err, domain.ErrInternalDatabase) {
			return newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	if _, err := a.botAPI.Send(msg); err != nil {
		return newHRError("Невозможно отправить сообщение!", err)
	}

	return nil
}

func (a *adapter) handleGet(u *tgbotapi.Update) error {
	if u.Message.CommandArguments() == "" {
		return newHRError("Никнейм не передан!", domain.ErrBotBadRequest)
	}

	text, err := a.service.GetStatsMessage(u.Message.CommandArguments())
	if err != nil {
		if errors.Is(err, domain.ErrInternalWargaming) {
			return newHRError("Ошибка при обращении к Wargaming API!", err)
		}
		if errors.Is(err, domain.ErrPlayerNotFound) {
			return newHRError("Игрок с данным никнеймом не найден!", err)
		}
		if errors.Is(err, domain.ErrInternalXVM) {
			return newHRError("Ошибка при обращении к XVM!", err)
		}

		return newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	if _, err := a.botAPI.Send(msg); err != nil {
		return newHRError("Невозможно отправить сообщение!", err)
	}

	return nil
}

func (a *adapter) handleSave(u *tgbotapi.Update) error {
	text, err := a.service.GetSaveNicknameMessage(u.Message.From.ID, u.Message.CommandArguments())
	if err != nil {
		if errors.Is(err, domain.ErrInternalWargaming) {
			return newHRError("Ошибка при обращении к Wargaming API!", err)
		}
		if errors.Is(err, domain.ErrPlayerNotFound) {
			return newHRError("Игрок с данным никнеймом не найден!", err)
		}
		if errors.Is(err, domain.ErrInternalXVM) {
			return newHRError("Ошибка при обращении к XVM!", err)
		}
		if errors.Is(err, domain.ErrInternalDatabase) {
			return newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	if _, err := a.botAPI.Send(msg); err != nil {
		return newHRError("Невозможно отправить сообщение!", err)
	}

	return nil
}

func (a *adapter) handleRefresh(u *tgbotapi.Update) error {
	text, err := a.service.GetRefreshMessage(u.Message.From.ID)
	if err != nil {
		if errors.Is(err, domain.ErrInternalXVM) {
			return newHRError("Ошибка при обращении к XVM!", err)
		}
		if errors.Is(err, domain.ErrNicknameNotSaved) {
			return newHRError("Сначала сохрани свой никнейм!", err)
		}
		if errors.Is(err, domain.ErrInternalDatabase) {
			return newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	if _, err := a.botAPI.Send(msg); err != nil {
		return newHRError("Невозможно отправить сообщение!", err)
	}

	return nil
}

func (a *adapter) handleMe(u *tgbotapi.Update) error {
	text, err := a.service.GetMeMessage(u.Message.From.ID)
	if err != nil {
		if errors.Is(err, domain.ErrTrendImageNotFound) {
			return newHRError("График не найден!", err)
		}
		if errors.Is(err, domain.ErrInternalDatabase) {
			return newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	if _, err := a.botAPI.Send(msg); err != nil {
		return newHRError("Невозможно отправить сообщение!", err)
	}

	return nil
}

func (a *adapter) handleTrend(u *tgbotapi.Update) error {
	img, err := a.service.GetTrendImage(u.Message.From.ID, "#"+u.Message.Command())
	if err != nil {
		if errors.Is(err, domain.ErrNicknameNotSaved) {
			return newHRError("Сначала сохрани свой никнейм!", err)
		}
		if errors.Is(err, domain.ErrInternalDatabase) {
			return newHRError("Ошибка при работе с базой! Обратитесь к администратору бота.", err)
		}

		return newHRError("Произошла неизвестная ошибка!", err)
	}

	msg := tgbotapi.NewPhotoUpload(u.Message.Chat.ID, tgbotapi.FileBytes{
		Name:  u.Message.Command(),
		Bytes: img,
	})
	if _, err := a.botAPI.Send(msg); err != nil {
		return newHRError("Невозможно отправить сообщение!", err)
	}

	return nil
}

func (a *adapter) error(update *tgbotapi.Update, err error) {
	if update == nil || err == nil {
		// Why did you call this function?
		return
	}

	// Log error
	a.logger.Error("Error occurred in handler!", zap.Error(err))

	// Send human readable representation of error to user to let him know
	if hrerr, ok := err.(*hrError); ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, hrerr.Human())
		_, err := a.botAPI.Send(msg)
		if err != nil {
			a.logger.Error("Error sending message with human readable error!", zap.Error(err))
		}
	} else {
		// ... do nothing? Unreadable error useless for people
	}
}
