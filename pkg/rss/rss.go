package rss

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

var gmtTimeZoneLocation *time.Location

func init() {
	loc, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
	gmtTimeZoneLocation = loc
}

var DafaultReader = New(&http.Client{})

func ReadFeed(url string, etag string, lastModified time.Time) (*Feed, error) {
	return DafaultReader.ReadFeed(url, etag, lastModified, context.Background())
}
func ReadFeedWithContext(url string, etag string, lastModified time.Time, ctx context.Context) (*Feed, error) {
	return DafaultReader.ReadFeed(url, etag, lastModified, ctx)
}

var ErrNotModified = errors.New("not modified")

func New(client *http.Client) *reader {
	return &reader{
		feedReader: gofeed.NewParser(),
		client:     client,
	}
}

type reader struct {
	feedReader *gofeed.Parser
	client     *http.Client
}

type Feed struct {
	*gofeed.Feed

	ETag         string
	LastModified time.Time
}

func (r *reader) ReadFeed(url string, etag string, lastModified time.Time, ctx context.Context) (*Feed, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "Gofeed/1.0")

	if etag != "" {
		req.Header.Set("If-None-Match", fmt.Sprintf(`"%s"`, etag))
	}

	req.Header.Set("If-Modified-Since", lastModified.In(gmtTimeZoneLocation).Format(time.RFC1123))

	resp, err := r.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode == http.StatusNotModified {
		return nil, ErrNotModified
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, gofeed.HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	feed := &Feed{}

	feedBody, err := r.feedReader.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	feed.Feed = feedBody

	if eTag := resp.Header.Get("Etag"); eTag != "" {
		feed.ETag = eTag
	}

	if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
		parsed, err := time.ParseInLocation(time.RFC1123, lastModified, gmtTimeZoneLocation)
		if err == nil {
			feed.LastModified = parsed
		}
	}

	return feed, nil
}
