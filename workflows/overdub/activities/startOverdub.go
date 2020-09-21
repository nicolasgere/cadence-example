package activities

import (
	"cadence-example/services/overdub"
	"context"
	"go.uber.org/cadence/activity"
)

type InputStartOverdub struct {
	Text string
}
type OutputStartOverdub struct {
	Id string
}

func StartOverdub(ctx context.Context, input *InputStartOverdub) (*OutputStartOverdub, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("sending overdub request: " + input.Text)
	r, err := overdub.StartOverdubWithApi(input.Text)
	if err != nil {
		return nil, err
	}
	return &OutputStartOverdub{
		Id: r.Id,
	}, nil
}
