package ytdownload

import (
	"context"
	"fmt"
	"sync"

	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
	"github.com/pkg/errors"
)

type YoutubeDownloadTask struct {
	YoutubeId string
	TargetDir string
	Provider  Provider

	// The following fields are protected by the mutex
	detailsMutex *sync.Mutex
	youtubeTitle string
}

const maxTitleLength = 50

func (y *YoutubeDownloadTask) Init() {
	y.detailsMutex = new(sync.Mutex)
}

func (y *YoutubeDownloadTask) String() string {
	y.detailsMutex.Lock()
	defer y.detailsMutex.Unlock()

	if y.youtubeTitle != "" {
		return fmt.Sprintf("Downloading '%s'", summariseTitle(y.youtubeTitle, maxTitleLength))
	}

	return "Downloading " + y.YoutubeId
}

func (y *YoutubeDownloadTask) Execute(ctx context.Context, runContext jobs.RunContext) error {
	runContext.PostUpdate(jobs.Update{Status: "Fetching video metadata"})
	metadata, err := y.Provider.GetVideoMetadata(ctx, y.YoutubeId)
	if err != nil {
		return errors.Wrap(err, "cannot get metadata")
	}

	y.setTitle(metadata.Title)

	return y.Provider.DownloadVideo(ctx, models.DownloadOptions{
		YoutubeID: y.YoutubeId,
		TargetDir: y.TargetDir,
	}, func(line string) {
		runContext.PostUpdate(jobs.Update{Status: line})
	})
}

func (y *YoutubeDownloadTask) setTitle(newTitle string) {
	y.detailsMutex.Lock()
	defer y.detailsMutex.Unlock()

	y.youtubeTitle = newTitle
}

func (y *YoutubeDownloadTask) VideoExtID() string {
	return y.YoutubeId
}

func (y *YoutubeDownloadTask) VideoTitle() string {
	return y.youtubeTitle
}

func summariseTitle(t string, maxLen int) string {
	if len(t) > maxLen {
		return t[:maxLen-3] + "..."
	}

	return t
}
