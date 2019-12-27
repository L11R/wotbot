package wargaming

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/L11R/wotbot/internal/domain"
	"go.uber.org/zap"
)

type adapter struct {
	logger *zap.Logger
	config *Config
}

func NewAdapter(logger *zap.Logger, config *Config) domain.Wargaming {
	a := &adapter{
		logger: logger,
		config: config,
	}

	return a
}

func (a *adapter) FindPlayer(nickname string) (string, int, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.worldoftanks.ru/wot/account/list/", nil)
	if err != nil {
		a.logger.Error("Error creating new Wargaming API request!", zap.Error(err))
		return "", 0, domain.ErrInternalWargaming
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	q := req.URL.Query()
	q.Set("application_id", a.config.ApplicationID)
	q.Set("search", nickname)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.logger.Error("Error doing Wargaming API request!", zap.Error(err))
		return "", 0, domain.ErrInternalWargaming
	}
	defer resp.Body.Close()

	var apiResp Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		a.logger.Error("Error decoding Wargaming API response!", zap.Error(err))
		return "", 0, domain.ErrInternalWargaming
	}

	if apiResp.Status != "ok" {
		a.logger.Error("Wargaming API returned an error!", zap.Error(err))
		return "", 0, domain.ErrInternalWargaming
	}

	var pp []PlayerData
	if err := json.Unmarshal(apiResp.Data, &pp); err != nil {
		a.logger.Error("Error decoding Wargaming API response!", zap.Error(err))
		return "", 0, domain.ErrInternalWargaming
	}

	for _, p := range pp {
		if strings.ToLower(p.Nickname) == strings.ToLower(nickname) {
			return p.Nickname, p.AccountID, nil
		}
	}

	return "", 0, domain.ErrPlayerNotFound
}
