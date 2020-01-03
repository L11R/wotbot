package kttc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/L11R/wotbot/internal/domain"
	"go.uber.org/zap"
)

type adapter struct {
	logger *zap.Logger
	config *Config
}

func NewAdapter(logger *zap.Logger, config *Config) domain.KTTC {
	a := &adapter{
		logger: logger,
		config: config,
	}

	return a
}

func (a *adapter) GetStats(accountID int) ([]*domain.KTTCStat, error) {
	req, err := http.NewRequest(http.MethodGet, "https://kttc.ru/wot/ru/statistics/user/get-by-battles/"+fmt.Sprint(accountID)+"/", nil)
	if err != nil {
		a.logger.Error("Error creating new KTTC stats request!", zap.Error(err))
		return nil, domain.ErrInternalKTTC
	}

	ctx, cancel := context.WithTimeout(context.Background(), a.config.HTTPTimeout)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.logger.Error("Error doing KTTC stats request!", zap.Error(err))
		return nil, domain.ErrInternalKTTC
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	var apiResp Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		a.logger.Error("Error decoding KTTC API response!", zap.Error(err))
		return nil, domain.ErrInternalKTTC
	}

	if !apiResp.Success {
		a.logger.Error("KTTC API returned an error!", zap.Error(fmt.Errorf(apiResp.Message)))
		return nil, domain.ErrInternalKTTC
	}

	var sbb StatsByBattles
	if err := json.Unmarshal(apiResp.Data, &sbb); err != nil {
		a.logger.Error("Error decoding KTTC API response!", zap.Error(err))
		return nil, domain.ErrInternalKTTC
	}

	kttcStats := make([]*domain.KTTCStat, 0)
	if latest, ok := sbb["1000"]; ok {
		wn8 := &domain.KTTCStat{
			Name:  "WN8",
			Value: latest.WN8,
			Delta: &latest.Deltas.WN8.Value,
		}

		switch {
		case /*latest.WN8 >= 448 &&*/ latest.WN8 < 978:
			wn8.Color = "❤️"
		case latest.WN8 >= 978 && latest.WN8 < 1574:
			wn8.Color = "💛"
		case latest.WN8 >= 1574 && latest.WN8 < 2371:
			wn8.Color = "💚"
		case latest.WN8 >= 2371 && latest.WN8 < 3188:
			wn8.Color = "💙"
		case latest.WN8 >= 3188:
			wn8.Color = "💜"
		}

		kttcStats = append(kttcStats, wn8)

		wtr := &domain.KTTCStat{
			Name:  "WTR",
			Value: float64(latest.WTR),
			Delta: &latest.Deltas.WTR.Value,
		}

		switch {
		case /*latest.WTR >= 2730 &&*/ latest.WTR < 4703:
			wtr.Color = "❤️"
		case latest.WTR >= 4703 && latest.WTR < 6800:
			wtr.Color = "💛"
		case latest.WTR >= 6800 && latest.WTR < 9050:
			wtr.Color = "💚"
		case latest.WTR >= 9050 && latest.WTR < 10496:
			wtr.Color = "💙"
		case latest.WTR >= 10496:
			wtr.Color = "💜"
		}

		kttcStats = append(kttcStats, wtr)

		winrate := &domain.KTTCStat{
			Name:  "Процент побед",
			Value: latest.Winrate,
			Delta: &latest.Deltas.Winrate.Value,
		}

		switch {
		case /*latest.Winrate >= 46.44 &&*/ latest.Winrate < 49.2:
			winrate.Color = "❤️"
		case latest.Winrate >= 49.2 && latest.Winrate < 52.54:
			winrate.Color = "💛"
		case latest.Winrate >= 52.54 && latest.Winrate < 57.81:
			winrate.Color = "💚"
		case latest.Winrate >= 57.81 && latest.Winrate < 63.81:
			winrate.Color = "💙"
		case latest.Winrate >= 63.81:
			winrate.Color = "💜"
		}

		kttcStats = append(kttcStats, winrate)

		damaged := &domain.KTTCStat{
			Name:  "Средний урон",
			Value: latest.Damaged,
			Delta: &latest.Deltas.Damaged.Value,
		}

		switch {
		case /*latest.Damaged >= 500 &&*/ latest.Damaged < 750:
			damaged.Color = "❤️"
		case latest.Damaged >= 750 && latest.Damaged < 1000:
			damaged.Color = "💛"
		case latest.Damaged >= 1000 && latest.Damaged < 1800:
			damaged.Color = "💚"
		case latest.Damaged >= 1800 && latest.Damaged < 2500:
			damaged.Color = "💙"
		case latest.Damaged >= 2500:
			damaged.Color = "💜"
		}

		kttcStats = append(kttcStats, damaged)

		hitsPercentage := &domain.KTTCStat{
			Name:  "Процент попадений",
			Value: latest.HitsPercentage,
			Delta: &latest.Deltas.HitsPercentage.Value,
		}

		switch {
		case /*latest.HitsPercentage >= 47.5 &&*/ latest.HitsPercentage < 60.5:
			hitsPercentage.Color = "❤️"
		case latest.HitsPercentage >= 60.5 && latest.HitsPercentage < 68.5:
			hitsPercentage.Color = "💛"
		case latest.HitsPercentage >= 68.5 && latest.HitsPercentage < 74.5:
			hitsPercentage.Color = "💚"
		case latest.HitsPercentage >= 74.5 && latest.HitsPercentage < 78.5:
			hitsPercentage.Color = "💙"
		case latest.HitsPercentage >= 78.5:
			hitsPercentage.Color = "💜"
		}

		kttcStats = append(kttcStats, hitsPercentage)
	} else {
		a.logger.Error("Error getting stats for the latest 1000 battles!")
		return nil, domain.ErrInternalKTTC
	}

	return kttcStats, nil
}
