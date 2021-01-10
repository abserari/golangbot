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

	// 所有查询从 context 获取客户端设置, 如未设置将使用默认客户端。
	var ctx = context.Background()

	ctx = papi.With(ctx, c.Papi)

	// 搜索画作
	// result, _ := artwork.Search(ctx, "パチュリー・ノーレッジ")
	// fmt.Println(result.JSON) // json return data.
	// result.Artworks() // []artwork.Artwork，只有部分数据，通过 `Fetch` `FetchPages` 方法获取完整数据。
	// artwork.Search(ctx, "パチュリー・ノーレッジ", artwork.SearchOptionPage(2)) // 获取第二页
	rank := &artwork.Rank{Mode: "daily_r18"}
	err := rank.Fetch(ctx)
	if err != nil {
		c.Logger.Info("login failed with cookies: " + config.Cookies)
		panic(err)
	}
	c.Logger.Info("login success")
	return &c
}

// Search
// offset would only be 1+. For example, if set offset = 0, would be 1.
// []artwork.Artwork，只有部分数据，通过 `Fetch` `FetchPages` 方法获取完整数据。
// search mode
// - s_tc: title & word
// - s_tag: partial consistent
func (c *Client) Search(ctx context.Context,
	query string, order artwork.Order,
	contentRating artwork.ContentRating,
	searchmode artwork.SearchMode,
	offset int) (art []artwork.Artwork, err error) {
	ctx = papi.With(ctx, c.Papi)
	// 搜索画作
	result, err := artwork.Search(ctx,
		query,
		artwork.SearchOptionPage(offset),
		artwork.SearchOptionOrder(order),
		artwork.SearchOptionMode(searchmode),
		artwork.SearchOptionContentRating(contentRating))
	if err != nil {
		return
	}
	art = result.Artworks()
	return
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
func (c *Client) Ranking(ctx context.Context, mode, content string, date time.Time, page int) ([]artwork.RankItem, error) {
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

	rank := &artwork.Rank{Mode: mode, Content: content, Date: date, Page: page}
	err := rank.Fetch(ctx)
	if err != nil {
		c.Logger.With(
			zap.String("mode", mode),
			zap.Int("page", page),
			zap.String("content", content),
			zap.Time("usedate", date),
		).Info("Got message")
		return nil, err
	}

	return rank.Items, nil
}
