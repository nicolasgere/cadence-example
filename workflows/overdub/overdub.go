package overdub

import (
	"cadence-example/workflows/overdub/activities"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

/**
 * This is the hello world workflow sample.
 */

// ApplicationName is the task list for this sample
const TaskListName = "helloWorldGroup"
const SignalName = "helloWorldSignal"

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(Workflow)
	activity.Register(activities.SplitText)
	activity.Register(activities.StartOverdub)
	activity.Register(activities.WaitCompletionOverdub)
	activity.Register(activities.DownloadAndMergeAudio)

}

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
}

func Workflow(ctx workflow.Context, text string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	logger := workflow.GetLogger(ctx)
	logger.Info("overdub workflow started")

	var OutputSplitText activities.OutputSplitText
	err := workflow.ExecuteActivity(ctx, activities.SplitText, activities.InputSplitText{
		Text: text,
	}).Get(ctx, &OutputSplitText)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return "", err
	}
	ids := make([]string, len(OutputSplitText.Chunks))
	for i, chunk := range OutputSplitText.Chunks {
		var OutputStartOverdub activities.OutputStartOverdub
		err := workflow.ExecuteActivity(ctx, activities.StartOverdub, &activities.InputStartOverdub{
			Text: chunk,
		}).Get(ctx, &OutputStartOverdub)
		if err != nil {
			logger.Error("Activity failed.", zap.Error(err))
			return "", err
		}
		ids[i] = OutputStartOverdub.Id
	}
	chunkResultChannel := workflow.NewChannel(ctx)
	for chunkIndex := 0; chunkIndex < len(ids); chunkIndex++ {
		i := chunkIndex
		workflow.Go(ctx, func(ctx workflow.Context) {
			var OutputWaitCompletionOverdub activities.OutputWaitCompletionOverdub
			err := workflow.ExecuteActivity(ctx, activities.WaitCompletionOverdub, &activities.InputWaitCompletionOverdub{
				Id:    ids[i],
				Index: i,
			}).Get(ctx, &OutputWaitCompletionOverdub)
			if err == nil {
				chunkResultChannel.Send(ctx, &OutputWaitCompletionOverdub)
			} else {
				chunkResultChannel.Send(ctx, err)
			}
		})
	}
	res := make([]string, len(ids))
	for i := 1; i <= len(ids); i++ {
		var v interface{}
		chunkResultChannel.Receive(ctx, &v)
		switch r := v.(type) {
		case error:
			return "", err
		case *activities.OutputWaitCompletionOverdub:
			res[r.Index] = r.Url
		}
	}
	var OutputDownloadAndMergeAudio activities.OutputDownloadAndMergeAudio

	err = workflow.ExecuteActivity(ctx, activities.DownloadAndMergeAudio, &activities.InputDownloadAndMergeAudio{
		Urls: res,
	}).Get(ctx, &OutputDownloadAndMergeAudio)
	if err != nil {
		return "", err
	}

	return OutputDownloadAndMergeAudio.Path, nil
}
