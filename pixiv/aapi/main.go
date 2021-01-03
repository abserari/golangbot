package main

import (
	"fmt"

	"github.com/everpcpc/pixiv"
)

func main() {
	account, err := pixiv.Login("abserari", "5Ek8ZQf4Z4vEp2U")
	app := pixiv.NewApp()
	illusts, err := app.IllustRanking("day_male_r18", "", "", "")
	fmt.Println(account, err, illusts)
}
