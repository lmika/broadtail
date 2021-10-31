package ytdownload

import (
	"context"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
)

type YoutubeDownloadTask struct {
	YoutubeId string
	TargetDir string
	Provider Provider
}

func (y *YoutubeDownloadTask) String() string {
	return "Downloading " + y.YoutubeId
}

func (y *YoutubeDownloadTask) Execute(ctx context.Context, runContext jobs.RunContext) error {
	return y.Provider.DownloadVideo(ctx, models.DownloadOptions{
		YoutubeID: y.YoutubeId,
		TargetDir: y.TargetDir,
	}, func(line string) {
		runContext.PostUpdate(jobs.Update{Status: line})
	})
}
