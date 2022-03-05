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
	runContext.PostUpdate(jobs.Update{Summary: "Initialising", Percent: 0.0})
	runContext.PostMessage("Fetching video metadata")

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
			runContext.PostMessage("warn: " + err.Error())
		}
	}

	// Download the video
	var outputFilename string
	for attempt := 1; attempt <= 3; attempt++ {
		runContext.PostMessage(fmt.Sprintf("Downloading video: attempt %d of 3", attempt))

		var err error
		outputFilename, err = y.DownloadProvider.DownloadVideo(ctx, models.DownloadOptions{
			YoutubeID: y.YoutubeId,
			TargetDir: y.TargetDir,
		}, func(line string) {
			if prog, ok := parseProgress(line); ok {
				runContext.PostUpdate(jobs.Update{
					Percent: prog.Percent,
					Summary: fmt.Sprintf("%.1f%% - ETA %v", prog.Percent, prog.ETA),
				})
			}
			runContext.PostMessage(line)
		})
		if err != nil {
			// Check that the context hasn't been cancelled
			if errors.Is(ctx.Err(), context.Canceled) {
				return err
			}

			// PARSE UPDATE
			runContext.PostMessage("Download error: " + err.Error())
			if attempt >= 3 {
				return errors.New("too many failed attempts")
			} else {
				runContext.PostMessage("Will sleep for 10 seconds, then try again")
				time.Sleep(10 * time.Second)
			}
		} else {
			break
		}
	}

	runContext.PostUpdate(jobs.Update{Summary: "Finalising", Percent: 100.0})

	// Check that the video is present
	stat, err := os.Stat(outputFilename)
	if err != nil {
		return errors.Wrap(err, "cannot stat saved file")
	}

	// If setting the owner
	if err := y.changeToTargetOwner(outputFilename); err != nil {
		runContext.PostMessage("warn: " + err.Error())
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
		runContext.PostMessage("warn: cannot delete existing video details: " + err.Error())
	}

	if err := y.VideoStore.Save(savedVideo); err != nil {
		runContext.PostMessage("warn: cannot save video details: " + err.Error())
	}

	//runContext.PostUpdate(jobs.Update{Summary: "Done", Percent: 100.0})
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
