package ytdownload

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
	"github.com/pkg/errors"
)

type YoutubeDownloadTask struct {
	YoutubeId   string
	TargetDir   string
	TargetOwner string

	DownloadProvider DownloadProvider
	VideoStore       VideoStore
	Feed             *models.Feed

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
	metadata, err := y.DownloadProvider.GetVideoMetadata(ctx, y.YoutubeId)
	if err != nil {
		return errors.Wrap(err, "cannot get metadata")
	}

	y.setTitle(metadata.Title)

	// Create the directory
	if y.TargetDir != "" {
		if err := os.MkdirAll(y.TargetDir, 0755); err != nil {
			return errors.Wrapf(err, "cannot create target directory: %v", y.TargetDir)
		}

		// FIXME: this needs to be recursive up to the library dir
		if err := y.changeToTargetOwner(y.TargetDir); err != nil {
			runContext.PostUpdate(jobs.Update{Status: "warn: " + err.Error()})
		}
	}

	// Download the video
	outputFilename, err := y.DownloadProvider.DownloadVideo(ctx, models.DownloadOptions{
		YoutubeID: y.YoutubeId,
		TargetDir: y.TargetDir,
	}, func(line string) {
		runContext.PostUpdate(jobs.Update{Status: line})
	})
	if err != nil {
		return err
	}

	// Check that the video is present
	stat, err := os.Stat(outputFilename)
	if err != nil {
		return errors.Wrap(err, "cannot stat saved file")
	}

	// If setting the owner
	if err := y.changeToTargetOwner(outputFilename); err != nil {
		runContext.PostUpdate(jobs.Update{Status: "warn: " + err.Error()})
	}

	// Save the downloaded file details
	videoExtId := models.ExtIDPrefixYoutube + y.YoutubeId
	savedVideo := models.SavedVideo{
		ID:       uuid.New(),
		ExtID:    videoExtId,
		Title:    metadata.Title,
		SavedOn:  time.Now(),
		Location: outputFilename,
		FileSize: stat.Size(),
	}
	if y.Feed != nil {
		savedVideo.FeedID = y.Feed.ID
		savedVideo.Source = y.Feed.Name
	} else {
		savedVideo.Source = "Manual download"
	}

	if err := y.VideoStore.DeleteWithExtID(videoExtId); err != nil {
		runContext.PostUpdate(jobs.Update{Status: "warn: cannot delete existing video details: " + err.Error()})
	}

	if err := y.VideoStore.Save(savedVideo); err != nil {
		runContext.PostUpdate(jobs.Update{Status: "warn: cannot save video details: " + err.Error()})
	}

	return nil
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

func (y *YoutubeDownloadTask) changeToTargetOwner(filename string) error {
	if y.TargetOwner == "" {
		return nil
	}

	targetUser, err := user.Lookup(y.TargetOwner)
	if err != nil {
		return errors.Wrapf(err, "unable to find target owner: %v", y.TargetOwner)
	}

	uid, err := strconv.Atoi(targetUser.Uid)
	if err != nil {
		return errors.Wrapf(err, "target uid not an int: %v", targetUser.Uid)
	}

	gid, err := strconv.Atoi(targetUser.Gid)
	if err != nil {
		return errors.Wrapf(err, "target primary gid not an int: %v", targetUser.Gid)
	}

	return errors.Wrapf(os.Chown(filename, uid, gid), "unable to chown file '%v' to user '%v'", filename, y.TargetOwner)
}

func summariseTitle(t string, maxLen int) string {
	if len(t) > maxLen {
		return t[:maxLen-3] + "..."
	}

	return t
}
