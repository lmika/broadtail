package ytdownload

import (
	"bufio"
	"context"
	"fmt"
	jobs2 "github.com/lmika/broadtail/providers/jobs"
	"github.com/pkg/errors"
	"os/exec"
)

type YoutubeDownloadTask struct {
	YoutubeId string
	TargetDir string
}

func (y *YoutubeDownloadTask) String() string {
	return "Downloading " + y.YoutubeId
}

func (y *YoutubeDownloadTask) Execute(ctx context.Context, runContext jobs2.RunContext) error {
	downloadUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%v", y.YoutubeId)
	cmd := exec.CommandContext(ctx, "python3", "/usr/local/bin/youtube-dl", "--newline", "-f", "mp4[height<=720]", downloadUrl)
	cmd.Dir = y.TargetDir

	stderrPipe, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "cannot open pipe to stderr")
	}

	pipeScanner := bufio.NewScanner(stderrPipe)

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "cannot start process")
	}

	for pipeScanner.Scan() {
		runContext.PostUpdate(jobs2.Update{Status: pipeScanner.Text()})
	}

	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "caught error waiting for process")
	}
	return nil
}
