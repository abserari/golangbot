package workflows

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abserari/golangbot/pkg/rss"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func RssWorker() {
	go func() {
		w := worker.New(C, "cron", worker.Options{})

		w.RegisterWorkflow(CronRssFeedWorkflow)

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	}()
}

// CronResult is used to return data from one cron run to the next
type CronResult struct {
	RunTime      time.Time
	Etag         string
	LastModified time.Time
}

var Localexecutor *TgBot = NewTgBot()

// SampleCronWorkflow executes on the given schedule
// The schedule is provided when starting the Workflow
func CronRssFeedWorkflow(ctx workflow.Context) (*CronResult, error) {
	workflow.GetLogger(ctx).Info("Cron workflow started.", "StartTime", workflow.Now(ctx))

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}

	ctx1 := workflow.WithActivityOptions(ctx, ao)

	// Start from 0 for first cron job
	lastRunTime := time.Time{}
	lastEtag := ""
	lastModefiedTime := time.Time{}
	// Check to see if there was a previous cron job
	if workflow.HasLastCompletionResult(ctx) {
		var lastResult CronResult
		if err := workflow.GetLastCompletionResult(ctx, &lastResult); err == nil {
			lastRunTime = lastResult.RunTime
			lastEtag = lastResult.Etag
			lastModefiedTime = lastResult.LastModified
		}
	}
	thisRunTime := workflow.Now(ctx)

	var feed rss.Feed

	lao := workflow.LocalActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy:         NoRetry,
	}
	ctx1 = workflow.WithLocalActivityOptions(ctx1, lao)

	err := workflow.ExecuteLocalActivity(ctx1, RssFeed, lastRunTime, thisRunTime, lastModefiedTime, lastEtag).Get(ctx, &feed)
	if err != nil {
		// Cron job failed
		// Next cron will still be scheduled by the Server
		workflow.GetLogger(ctx).Error("Cron job failed.", "Error", err)
		return nil, err
	}

	// not modified
	if feed.ETag == lastEtag {
		return &CronResult{RunTime: thisRunTime, Etag: feed.ETag, LastModified: feed.LastModified}, nil
	}

	var url string

	err = workflow.ExecuteLocalActivity(ctx1, Localexecutor.SendMessageToTelegraph, PrintFeedToTelegraph(feed)).Get(ctx, &url)
	if err != nil {
		workflow.GetLogger(ctx).Error("PrintFeedToTelegraph failed.", "Error", err)
		return nil, err
	}
	workflow.GetLogger(ctx).Info("URL:", url)

	return &CronResult{RunTime: thisRunTime, Etag: feed.ETag, LastModified: feed.LastModified}, nil
}

// DoSomething is an Activity
func RssFeed(ctx context.Context, lastRunTime, thisRunTime, lastmodified time.Time, etag string) (*rss.Feed, error) {
	activity.GetLogger(ctx).Info("Cron job running.", "lastRunTime_exclude", lastRunTime, "thisRunTime_include", thisRunTime)

	return rss.ReadFeedWithContext("https://yusank.space/index.xml", "", time.Time{}, ctx)

}

func PrintFeedToTelegraph(feed rss.Feed) string {

	header := `<figure>
	<img src="/file/6a5b15e7eb4d7329ca7af.jpg"/>
</figure>
<p><i>Hello</i>, my name is <b>Tech Cats</b>, <u>look at me</u>!</p>
`

	article := header

	for _, v := range feed.Items {
		pub := fmt.Sprintf(`<a href="%s">%s</a> 
<p>%s</p>`, v.Link, v.Title, v.Description)

		article += pub
	}

	return article
}
