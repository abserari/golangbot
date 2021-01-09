package main

import (
	"context"
	"log"
	"os"

	"github.com/NateScarlet/pixiv/pkg/artwork"
	"github.com/NateScarlet/pixiv/pkg/client"
)

var (
	PHPSESSID = os.Getenv("PIXIV_PHPSESSID")
)

func main() {
	// 使用 PHPSESSID Cookie 登录 (推荐)。
	c := &client.Client{}
	c.SetPHPSESSID("38448763_xzUGIUnmDwHQkGzanCoAXROOrN6G7JWR")
	// // 通过账号密码登录(可能触发 reCAPTCHA)。
	c.SetDefaultHeader("User-Agent", client.DefaultUserAgent)

	// 所有查询从 context 获取客户端设置, 如未设置将使用默认客户端。
	var ctx = context.Background()

	ctx = client.With(ctx, c)

	// 搜索画作
	// result, _ := artwork.Search(ctx, "パチュリー・ノーレッジ")
	// fmt.Println(result.JSON)                                        // json return data.
	// []artwork.Artwork，只有部分数据，通过 `Fetch` `FetchPages` 方法获取完整数据。
	result, _ := artwork.Search(ctx, "パチュリー・ノーレッジ", artwork.SearchOptionPage(2), artwork.SearchOptionMode("r18"), artwork.SearchOptionOrder("date"))
	data := result.Artworks()
	for _, i := range data {
		var found bool
		for _, v := range i.Tags {
			if v != "R-18" && v != "R-18G" {
				continue
			}
			found = true
		}
		log.Println(found)

	}
	// rank := &artwork.Rank{Mode: "daily_r18"}
	// for {
	// 	err := rank.Fetch(ctx)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		time.Sleep(time.Second)
	// 		continue
	// 	}
	// 	break
	// }

	// fmt.Println(rank.Items[0].Image, rank.Items[1].Image, rank.Items[2].Image)
}
