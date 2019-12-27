package xvm

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/L11R/wotbot/internal/domain"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type adapter struct {
	logger *zap.Logger
	config *Config
}

func NewAdapter(logger *zap.Logger, config *Config) domain.XVM {
	a := &adapter{
		logger: logger,
		config: config,
	}

	return a
}

//noinspection GoUnhandledErrorResult
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

		name := selection.Find(".h5").Text()
		value := selection.Find(".h2").Text()

		ss = append(ss, &domain.Stat{
			Type:   domain.TrendStat,
			Name:   name,
			Value:  &value,
			HtmlID: id,
		})
	})

	// Silly, but working way; Feel free to suggest proper solution
	doc.Find("#battlesByVehicleType").Parent().Parent().Find("div").Each(func(i int, selection *goquery.Selection) {
		id, ok := selection.Find("canvas").Attr("id")
		if !ok {
			return
		}
		id = "#" + id

		// Parse chart title from JS script
		r := regexp.MustCompile(`text: "(.+)"`)
		match := r.FindStringSubmatch(selection.Find("script").Text())
		if len(match) != 2 {
			return
		}
		name := match[1]

		ss = append(ss, &domain.Stat{
			Type:   domain.VehicleStat,
			Name:   name,
			HtmlID: id,
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
				chromedp.Screenshot(ss[i].HtmlID, &ss[i].Image, chromedp.NodeVisible, chromedp.ByID),
			)
		}

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		debuggerURL, err := a.getWebSocketDebuggerURL()
		if err != nil {
			return nil, err
		}
		remoteAllocatorCtx, cancel := chromedp.NewRemoteAllocator(timeoutCtx, debuggerURL)
		defer cancel()
		ctx, cancel = chromedp.NewContext(remoteAllocatorCtx)
		defer cancel()

		if err := chromedp.Run(ctx, tt); err != nil {
			a.logger.Error("Error while taking screenshot!", zap.Error(err))
			return nil, domain.ErrInternalXVM
		}
	}

	return ss, nil
}

func (a *adapter) getWebSocketDebuggerURL() (string, error) {
	// Chrome could not resolve local docker service name :\
	u, err := url.Parse(a.config.ChromeDevtoolsURL)
	if err != nil {
		a.logger.Error("Error parsing Chrome Devtools URL!", zap.Error(err))
		return "", domain.ErrInternalXVM
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", u.Host)
	if err != nil {
		a.logger.Error("Error resolving Chrome Devtools host!", zap.Error(err))
		return "", domain.ErrInternalXVM
	}
	u.Host = tcpAddr.String()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		a.logger.Error("Error creating new Chrome Devtools request!", zap.Error(err))
		return "", domain.ErrInternalXVM
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.logger.Error("Error doing XVM stats request!", zap.Error(err))
		return "", domain.ErrInternalXVM
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	var target []struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		a.logger.Error("Error decoding Chrome Devtools response!", zap.Error(err))
		return "", domain.ErrInternalXVM
	}

	if len(target) != 1 {
		a.logger.Error("Received an empty targets array!")
		return "", domain.ErrInternalXVM
	}

	return target[0].WebSocketDebuggerURL, nil
}
