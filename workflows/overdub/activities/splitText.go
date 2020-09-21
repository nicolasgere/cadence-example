package activities

import (
	"context"
	"go.uber.org/cadence/activity"
	"gopkg.in/neurosnap/sentences.v1/english"
)

type InputSplitText struct {
	Text string
}
type OutputSplitText struct {
	Chunks []string
}

func SplitText(ctx context.Context, input *InputSplitText) (*OutputSplitText, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("starting splitting text")
	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		return nil, err
	}
	sentences := tokenizer.Tokenize(input.Text)
	res := make([]string, len(sentences))
	for i, s := range sentences {
		res[i] = s.Text
	}
	return &OutputSplitText{
		Chunks: res,
	}, nil
}
