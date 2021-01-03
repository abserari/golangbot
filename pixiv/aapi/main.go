package main

import (
	"fmt"
	"os"

	"github.com/everpcpc/pixiv"
)

var (
	username = os.Getenv("PIXIV_USERNAME")
	password = os.Getenv("PIXIV_PASSWORD")
)

func main() {
	_, _ = pixiv.Login(username, password)
	app := pixiv.NewApp()
	// searchTarget - Search type
	//	"partial_match_for_tags"  - The label part is consistent
	//	"exact_match_for_tags"    - The labels are exactly the same
	//	"title_and_caption"       - Title description
	//
	// sort: [date_desc, date_asc]
	//
	// duration: [within_last_day, within_last_week, within_last_month]
	searchResult, _ := app.SearchIllust("r18", "partial_match_for_tags", "date_desc", "within_last_week", "", 0)
	for _, v := range searchResult.Illusts {
		fmt.Println(v.Images.Large, v.PageCount, v.MetaPages)
	}
}
