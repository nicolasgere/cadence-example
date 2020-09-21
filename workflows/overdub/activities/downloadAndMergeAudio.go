package activities

import (
	"bytes"
	"context"
	"fmt"
	"go.uber.org/cadence/activity"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type InputDownloadAndMergeAudio struct {
	Urls []string
}

type OutputDownloadAndMergeAudio struct {
	Path string
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadAndMergeAudio(ctx context.Context, input *InputDownloadAndMergeAudio) (*OutputDownloadAndMergeAudio, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("download all audio")
	dir, err := ioutil.TempDir("temp", "test")
	if err != nil {
		return nil, err
	}
	paths := []string{}
	for i, url := range input.Urls {
		path := dir + "/" + strconv.Itoa(i) + ".wav"
		paths = append(paths, path)
		err := downloadFile(path, url)
		if err != nil {
			return nil, err
		}
	}
	paths = append(paths, dir+"/final.wav")
	cmd := exec.Command("sox", paths...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return &OutputDownloadAndMergeAudio{
		Path: dir,
	}, nil
}
