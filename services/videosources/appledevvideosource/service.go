package appledevvideosource

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/videosources/appledevvideosource/httprange"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Service struct {
	client *resty.Client
}

func NewService() *Service {
	return &Service{
		client: resty.New(),
	}
}

func (s *Service) GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error) {
	dom, err := s.getVideoPage(ctx, videoRef)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(dom.Find("li.supplement.details h1").Text())
	description := strings.TrimSpace(dom.Find("li.supplement.details p").Text())

	return &models.Video{
		VideoRef:    videoRef,
		Title:       title,
		Description: description,
		Duration:    0,
	}, nil
}

func (s *Service) getVideoPage(ctx context.Context, videoRef models.VideoRef) (*goquery.Document, error) {
	videoURL := s.GetVideoURL(videoRef)
	resp, err := s.client.R().Get(videoURL)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get page: %v", videoURL)
	} else if !resp.IsSuccess() {
		return nil, errors.Errorf("non-200 response code: %v", resp.Status())
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse dom")
	}

	return dom, nil
}

func (s *Service) DownloadVideo(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions, logline func(line models.LogMessage)) (outputFilename string, err error) {
	logline(models.LogMessage{Message: "Getting video page"})
	dom, err := s.getVideoPage(ctx, videoRef)
	if err != nil {
		return "", err
	}

	logline(models.LogMessage{Message: "Looking for download link"})
	var videoDownloadURL *url.URL
	dom.Find("li.download a").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if selection.Text() == "HD Video" {
			urlStr := strings.TrimSuffix(selection.AttrOr("href", ""), "?dl=1")
			parsedUrl, err := url.Parse(urlStr)
			if err != nil {
				logline(models.LogMessage{Message: "warn: invalid URL: " + err.Error()})
			}
			videoDownloadURL = parsedUrl
			return false
		}
		return true
	})
	if videoDownloadURL == nil {
		return "", errors.Errorf("cannot find download link for video: %v", videoRef.ID)
	}

	logline(models.LogMessage{Message: "Download link found.  Starting download."})
	resp, err := s.client.R().SetDoNotParseResponse(true).SetContext(ctx).Get(videoDownloadURL.String())
	if err != nil {
		return "", errors.Wrap(err, "cannot download video")
	}
	defer resp.RawBody()

	if !resp.IsSuccess() {
		return "", errors.Errorf("non-200 response code: %v", resp.Status())
	}

	targetFilename := filepath.Join(options.TargetDir, filepath.Base(videoDownloadURL.Path))
	partTargetFilename := targetFilename + ".part"

	logline(models.LogMessage{Message: "Downloading"})
	if err := s.downloadFile(ctx, partTargetFilename, videoDownloadURL.String(), logline); err != nil {
		if err2 := os.Remove(partTargetFilename); err2 != nil {
			logline(models.LogMessage{Message: fmt.Sprintf("warn: cannot remove target file '%v': %v", targetFilename, err2)})
		}
		return "", errors.Wrap(err, "unable to download video")
	}

	if err := os.Rename(partTargetFilename, targetFilename); err != nil {
		return "", errors.Wrap(err, "unable to move working download file to final filename")
	}

	logline(models.LogMessage{Message: "Download complete"})

	return targetFilename, nil
}

func (s *Service) downloadFile(ctx context.Context, targetFilename, url string, logline func(line models.LogMessage)) error {
	f, err := os.Create(targetFilename)
	if err != nil {
		return errors.Wrapf(err, "cannot open target file: %v", targetFilename)
	}
	defer f.Close()

	var sendUpdateNext time.Time
	if _, err := httprange.Get(url).WithWriteObserver(func(writtenSoFar int64, totalExpectedSize int64) {
		if time.Now().After(sendUpdateNext) {
			pmille := int((writtenSoFar * 1000) / totalExpectedSize)
			logline(models.LogMessage{
				Permille: pmille,
			})
			sendUpdateNext = time.Now().Add(250 * time.Millisecond)
		}
	}).WriteTo(ctx, f); err != nil {
		return errors.Wrapf(err, "cannot get video from URL %v", url)
	}
	return nil
}

func (s *Service) GetVideoURL(videoRef models.VideoRef) string {
	videoSet, video, _ := strings.Cut(videoRef.ID, ".")
	return fmt.Sprintf("https://developer.apple.com/videos/play/%s/%s/", url.PathEscape(videoSet), url.PathEscape(video))
}
