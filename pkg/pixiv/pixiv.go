package pixiv

import (
	"context"
	"time"

	"github.com/NateScarlet/pixiv/pkg/artwork"
	papi "github.com/NateScarlet/pixiv/pkg/client"
	aapi "github.com/everpcpc/pixiv"
	"go.uber.org/zap"
)

type Client struct {
	Papi   *papi.Client
	Aapi   *aapi.AppPixivAPI
	Logger *zap.Logger
}
type ClientConfig struct {
	Username string
	Password string
	Cookies  string
	Logger   *zap.Logger
}

func NewClient(config *ClientConfig) *Client {
	var c Client
	// 使用 PHPSESSID Cookie 登录 (推荐)。
	c.Papi = &papi.Client{}
	c.Papi.SetPHPSESSID(config.Cookies)
	c.Papi.SetDefaultHeader("User-Agent", papi.DefaultUserAgent)
	// 通过账号密码登录(可能触发 reCAPTCHA)。

	_, _ = aapi.Login(config.Username, config.Password)
	c.Aapi = aapi.NewApp()

	c.Logger = config.Logger
	return &c
}

// Ranking return ranking with ppai without login status.
// required, possible rank modes:
// 	- daily (default)
//  - weekly
//  - monthly
//  - rookie
//  - original
//  - male
//  - female
//  - daily_r18
//  - weekly_r18
//  - male_r18
//  - female_r18
//  - r18g
// optional, possible rank content:
//  - all (default)
//  - illust
//  - ugoira
//  - manga
// date: YYYYMMDD (default is yesterday)
// page: pages
func (c *Client) Ranking(ctx context.Context, mode, content, date string, page int) ([]artwork.RankItem, error) {
	ctx = papi.With(ctx, c.Papi)
	switch mode {
	case "daily":
	case "weekly":
	case "monthly":
	case "rookie":
	case "original":
	case "male":
	case "female":
	case "daily_r18":
	case "weekly_r18":
	case "male_r18":
	case "female_r18":
	case "r18g":
	default:
		mode = "daily"
	}
	switch content {
	case "all":
	case "illust":
	case "ugoira":
	case "manga":
	default:
		content = "all"
	}

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		t = time.Now().Local().AddDate(0, 0, -1)
	}
	rank := &artwork.Rank{Mode: mode, Content: content, Date: t, Page: page}
	err = rank.Fetch(ctx)
	if err != nil {
		c.Logger.With(
			zap.String("mode", mode),
			zap.Int("page", page),
			zap.String("content", content),
			zap.Time("usedate", t),
		).Info("Got message")
		return nil, err
	}

	return rank.Items, nil
}
