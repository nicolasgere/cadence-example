package activities

import (
	"cadence-example/services/overdub"
	"context"
	"go.uber.org/cadence/activity"
	"time"
)

type InputWaitCompletionOverdub struct {
	Id    string
	Index int
}
type OutputWaitCompletionOverdub struct {
	Index int
	Url   string
}

func WaitCompletionOverdub(ctx context.Context, input *InputWaitCompletionOverdub) (*OutputWaitCompletionOverdub, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("waiting overdub: " + input.Id)
	delay := 2 * time.Second
	for {
		r, err := overdub.GetOverdubWithApi(input.Id)
		if err != nil {
			continue
		}
		if r.Url != "" {
			return &OutputWaitCompletionOverdub{
				Url:   r.Url,
				Index: input.Index,
			}, nil
		}
		time.Sleep(delay)
	}
}
