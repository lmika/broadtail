package appledevvideosource

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

func (s *Service) DownloadVideo(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions, logline func(line string)) (outputFilename string, err error) {
	logline("Getting video page")
	dom, err := s.getVideoPage(ctx, videoRef)
	if err != nil {
		return "", err
	}

	logline("Looking for download link")
	var videoDownloadURL *url.URL
	dom.Find("li.download a").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if selection.Text() == "HD Video" {
			urlStr := strings.TrimSuffix(selection.AttrOr("href", ""), "?dl=1")
			parsedUrl, err := url.Parse(urlStr)
			if err != nil {
				logline("warn: invalid URL: " + err.Error())
			}
			videoDownloadURL = parsedUrl
			return false
		}
		return true
	})
	if videoDownloadURL == nil {
		return "", errors.Errorf("cannot find download link for video: %v", videoRef.ID)
	}

	logline("Download link found.  Starting download.")
	resp, err := s.client.R().SetDoNotParseResponse(true).Get(videoDownloadURL.String())
	if err != nil {
		return "", errors.Wrap(err, "cannot download video")
	}
	defer resp.RawBody()

	if !resp.IsSuccess() {
		return "", errors.Errorf("non-200 response code: %v", resp.Status())
	}

	targetFilename := filepath.Join(options.TargetDir, filepath.Base(videoDownloadURL.Path))
	f, err := os.Create(targetFilename)
	if err != nil {
		return "", errors.Wrap(err, "cannot open target file")
	}
	defer f.Close()

	logline("Downloading")
	if _, err := io.Copy(f, resp.RawBody()); err != nil {
		return "", errors.Wrap(err, "cannot download video")
	}
	logline("Download complete")
	return targetFilename, nil
}

func (s *Service) GetVideoURL(videoRef models.VideoRef) string {
	videoSet, video, _ := strings.Cut(videoRef.ID, ".")
	return fmt.Sprintf("https://developer.apple.com/videos/play/%s/%s/", url.PathEscape(videoSet), url.PathEscape(video))
}
