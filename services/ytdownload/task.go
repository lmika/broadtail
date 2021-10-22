package ytdownload

import (
	"bufio"
	"context"
	"github.com/lmika/broadtail/jobs"
	"github.com/pkg/errors"
	"os/exec"
)

type YoutubeDownloadTask struct {
}

func (y *YoutubeDownloadTask) String() string {
	return "Download XX"
}

func (y *YoutubeDownloadTask) Execute(ctx context.Context, runContext jobs.RunContext) error {
	cmd := exec.Command("youtube-dl", "--newline", "-f", "mp4[height<=720]", "https://www.youtube.com/watch?v=BaW_jenozKc")

	stderrPipe, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "cannot open pipe to stderr")
	}

	pipeScanner := bufio.NewScanner(stderrPipe)

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "cannot start process")
	}

	for pipeScanner.Scan() {
		runContext.PostUpdate(jobs.Update{Status: pipeScanner.Text()})
	}

	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "caught error waiting for procress")
	}
	return nil
}
