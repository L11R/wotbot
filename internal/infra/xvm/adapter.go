package xvm

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"net/http"
	"time"

	"github.com/L11R/wotbot/internal/domain"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type adapter struct {
	logger *zap.Logger
}

func NewAdapter(logger *zap.Logger) domain.XVM {
	a := &adapter{
		logger: logger,
	}

	return a
}

func (a *adapter) GetStats(accountID int, withTrend bool) ([]*domain.Stat, error) {
	req, err := http.NewRequest(http.MethodGet, "https://stats.modxvm.com/ru/stat/players/"+fmt.Sprint(accountID), nil)
	if err != nil {
		a.logger.Error("Error creating new XVM stats request!", zap.Error(err))
		return nil, domain.ErrInternalXVM
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.logger.Error("Error doing XVM stats request!", zap.Error(err))
		return nil, domain.ErrInternalXVM
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		a.logger.Error("Error creating document from reader!", zap.Error(err))
		return nil, domain.ErrInternalXVM
	}

	var ss []*domain.Stat
	doc.Find(".stats-summary a").Each(func(i int, selection *goquery.Selection) {
		id, ok := selection.Attr("href")
		if !ok {
			return
		}

		ss = append(ss, &domain.Stat{
			HtmlID: id,
			Name:   selection.Find(".h5").Text(),
			Value:  selection.Find(".h2").Text(),
		})
	})

	if withTrend {
		tt := chromedp.Tasks{
			chromedp.EmulateViewport(1920, 7666),
			chromedp.Navigate("https://stats.modxvm.com/ru/stat/players/" + fmt.Sprint(accountID)),
		}
		for i := range ss {
			tt = append(
				tt,
				chromedp.Screenshot(ss[i].HtmlID, &ss[i].TrendImage, chromedp.NodeVisible, chromedp.ByID),
			)
		}

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ctx, cancel = chromedp.NewContext(timeoutCtx)
		defer cancel()

		if err := chromedp.Run(ctx, tt); err != nil {
			a.logger.Error("Error while taking screenshot!", zap.Error(err))
			return nil, domain.ErrInternalXVM
		}
	}

	return ss, nil
}
