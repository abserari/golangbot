package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abserari/golangbot/pkg/rss"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	feed, err := rss.ReadFeedWithContext("https://yusank.space/index.xml", "61d83d75-18aa", time.Time{}, ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed.Title, feed.ETag)
	// for _, v := range feed.Items {
	// 	log.Println(v.Title, v.Link)
	// 	log.Println(v.PublishedParsed.Date())
	// 	log.Println(v.Description)
	// }
}
