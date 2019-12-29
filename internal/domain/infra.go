package domain

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type Service interface {
	GetCreateUserMessage(telegramID int) (string, error)
	GetSaveNicknameMessage(telegramID int, nickname string) (string, error)
	GetRefreshMessage(telegramID int) (string, error)
	GetMeMessage(telegramID int, chatType string) (string, error)
	GetTrendImage(telegramID int, htmlID string) ([]byte, error)
	GetStatsMessage(nickname string) (string, error)
}

type Wargaming interface {
	FindPlayer(nickname string) (string, int, error)
}

type XVM interface {
	GetStats(accountID int, withTrend bool) ([]*Stat, error)
}

type Database interface {
	GetUserByTelegramID(telegramID int) (*User, error)
	UpsertUser(user *User) (*User, error)
	GetStatsByUserID(userID int) ([]*Stat, error)
	UpdateStatsByUserID(userID int, stats []*Stat) ([]*Stat, error)
}

type service struct {
	logger    *zap.Logger
	database  Database
	wargaming Wargaming
	xvm       XVM
}

func NewService(logger *zap.Logger, database Database, wargaming Wargaming, xvm XVM) Service {
	s := &service{
		logger:    logger,
		database:  database,
		wargaming: wargaming,
		xvm:       xvm,
	}

	return s
}

func (s *service) GetCreateUserMessage(telegramID int) (string, error) {
	user, err := s.database.UpsertUser(&User{
		TelegramID: telegramID,
	})
	if err != nil {
		s.logger.Error("Error upserting user!", zap.Int("telegram_id", telegramID), zap.Error(err))
		return "", err
	}

	msg := `Привет.

Команды:
/get <i>nickname</i> — запрашивает и отображает статистику игрока.
/save <i>nickname</i> — позволяет сохранить свой никнейм.
/me — выводит расширенную статистику по сохранённому никнейму.
/refresh — обновляет кэш.`

	if user.Nickname != nil {
		msg += fmt.Sprintf("\n\nКстати, ты уже сохранил свой никнейм, приветствую <b>%s</b>!", *user.Nickname)
	}

	return msg, nil
}

func (s *service) GetSaveNicknameMessage(telegramID int, nickname string) (string, error) {
	nickname, accountID, err := s.wargaming.FindPlayer(nickname)
	if err != nil {
		s.logger.Error("Error getting account_id!", zap.String("nickname", nickname), zap.Error(err))
		return "", err
	}

	if _, err = s.database.UpsertUser(&User{
		TelegramID:  telegramID,
		Nickname:    &nickname,
		WargamingID: &accountID,
	}); err != nil {
		s.logger.Error(
			"Error upserting user!",
			zap.Int("telegram_id", telegramID),
			zap.String("nickname", nickname),
			zap.Int("wargaming_id", accountID),
			zap.Error(err),
		)
		return "", err
	}

	stats, err := s.xvm.GetStats(accountID, true)
	if err != nil {
		s.logger.Error("Error getting XVM stats!", zap.Int("wargaming_id", accountID), zap.Error(err))
		return "", err
	}

	user, err := s.database.GetUserByTelegramID(telegramID)
	if err != nil {
		s.logger.Error("Error getting user by telegram_id!", zap.Int("telegram_id", telegramID), zap.Error(err))
		return "", err
	}

	if _, err = s.database.UpdateStatsByUserID(user.ID, stats); err != nil {
		s.logger.Error("Error updating stats by user_id!", zap.Int("user_id", telegramID), zap.Error(err))
		return "", err
	}

	return "Твой никнейм сохранён, ты можешь посмотреть свою статистику здесь: /me", nil
}

func (s *service) GetRefreshMessage(telegramID int) (string, error) {
	user, err := s.database.GetUserByTelegramID(telegramID)
	if err != nil {
		s.logger.Error("Error getting user!", zap.Int("telegram_id", telegramID), zap.Error(err))
		return "", err
	}

	if user.WargamingID == nil {
		s.logger.Error("User Wargaming ID is null, he didn't save nickname!")
		return "", ErrNicknameNotSaved
	}

	stats, err := s.xvm.GetStats(*user.WargamingID, true)
	if err != nil {
		s.logger.Error("Error getting XVM stats!", zap.Int("wargaming_id", *user.WargamingID), zap.Error(err))
		return "", err
	}

	if _, err = s.database.UpdateStatsByUserID(user.ID, stats); err != nil {
		s.logger.Error("Error updating stats by user_id!", zap.Int("user_id", telegramID), zap.Error(err))
		return "", err
	}

	return "Статистика обновлена!", nil
}

func (s *service) GetMeMessage(telegramID int, chatType string) (string, error) {
	user, err := s.database.GetUserByTelegramID(telegramID)
	if err != nil {
		s.logger.Error("Error getting user!", zap.Int("telegram_id", telegramID), zap.Error(err))
		return "", err
	}

	ss, err := s.database.GetStatsByUserID(user.ID)
	if err != nil {
		s.logger.Error("Error getting stats by user_id!", zap.Int("user_id", telegramID), zap.Error(err))
		return "", err
	}

	if user.Nickname == nil || user.WargamingID == nil {
		s.logger.Error("User nickname or Wargaming ID are null, he didn't save nickname!")
		return "", ErrNicknameNotSaved
	}

	msg := fmt.Sprintf("<b>Игрок:</b> %s <a href=\"https://stats.modxvm.com/ru/stat/players/%d\">(на сайте)</a>\n\n", *user.Nickname, *user.WargamingID)
	for _, s := range ss {
		if s.Value != nil {
			if chatType == "private" {
				msg += fmt.Sprintf("<b>%s:</b> %s %s\n", s.Name, *s.Value, strings.Replace(s.HtmlID, "#", "/", 1))
			} else {
				msg += fmt.Sprintf("<b>%s:</b> %s\n", s.Name, *s.Value)
			}
		}
	}

	if chatType == "private" {
		msg += "\n<b>Техника:</b>\n"
		for _, s := range ss {
			if s.Value == nil {
				msg += fmt.Sprintf("<b>%s:</b> %s\n", s.Name, strings.Replace(s.HtmlID, "#", "/", 1))
			}
		}
	}

	return msg, nil
}

func (s *service) GetTrendImage(telegramID int, htmlID string) ([]byte, error) {
	user, err := s.database.GetUserByTelegramID(telegramID)
	if err != nil {
		s.logger.Error("Error getting user!", zap.Int("telegram_id", telegramID), zap.Error(err))
		return nil, err
	}

	ss, err := s.database.GetStatsByUserID(user.ID)
	if err != nil {
		s.logger.Error("Error getting stats by user_id!", zap.Int("user_id", telegramID), zap.Error(err))
		return nil, err
	}

	for i := range ss {
		if ss[i].HtmlID == htmlID {
			return ss[i].Image, nil
		}
	}

	return nil, ErrTrendImageNotFound
}

func (s *service) GetStatsMessage(nickname string) (string, error) {
	nickname, accountID, err := s.wargaming.FindPlayer(nickname)
	if err != nil {
		s.logger.Error("Error getting account_id!", zap.String("nickname", nickname), zap.Error(err))
		return "", err
	}

	ss, err := s.xvm.GetStats(accountID, false)
	if err != nil {
		s.logger.Error("Error getting stats!", zap.Int("account_id", accountID), zap.Error(err))
		return "", err
	}

	msg := fmt.Sprintf("<b>Игрок:</b> %s\n\n", nickname)
	for _, s := range ss {
		if s.Value != nil {
			msg += fmt.Sprintf("<b>%s:</b> %s\n", s.Name, *s.Value)
		}
	}

	if len(ss) == 0 {
		msg += fmt.Sprintf("Показатели на найдены.")
	}

	return msg, nil
}
