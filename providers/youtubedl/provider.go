package youtubedl

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type Provider struct {
}

func New() (*Provider, error) {
	provider := &Provider{}
	if err := provider.checkAvailable(); err != nil {
		return nil, err
	}

	return provider, nil
}

func (p *Provider) checkAvailable() error {
	cmd := exec.Command("python3", "/usr/local/bin/youtube-dl", "--version")
	return errors.Wrap(cmd.Run(), "youtube-dl is not available")
}

func (p *Provider) GetVideoMetadata(ctx context.Context, youtubeId string) (*models.Video, error) {
	jsonData, err := p.videoMetadata(ctx, youtubeId)
	if err != nil {
		return nil, err
	}

	// Decode the upload date
	uploadDate, err := jsonData.UploadDate()
	if err != nil {
		log.Printf("invalid upload date '%v': %v", uploadDate, err)
	}

	return &models.Video{
		ExtID:        youtubeId,
		Title:        jsonData.Title,
		ChannelID:    jsonData.ChannelID,
		ChannelName:  jsonData.Channel,
		Description:  jsonData.Description,
		ThumbnailURL: jsonData.ThumbnailURL,
		UploadedOn:   uploadDate,
		Duration:     jsonData.Duration,
	}, nil
}

func (p *Provider) DownloadVideo(ctx context.Context, options models.DownloadOptions, logline func(line string)) (string, error) {
	const filenameFormat = "%(title)s.%(id)s.%(ext)s"

	// Get the expected filename
	downloadUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%v", options.YoutubeID)
	filenameCmd := exec.CommandContext(ctx, "python3", "/usr/local/bin/youtube-dl",
		"--get-filename", "-o", filenameFormat, "--restrict-filenames", downloadUrl)
	out, err := filenameCmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "unable to determine target filename")
	}
	outFilename := strings.TrimSpace(string(out))

	cmd := exec.CommandContext(ctx, "python3", "/usr/local/bin/youtube-dl",
		"-o", filenameFormat, "--restrict-filenames",
		"--newline", "-f", "mp4[height<=720]", downloadUrl)
	cmd.Dir = options.TargetDir

	stderrPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "cannot open pipe to stderr")
	}

	pipeScanner := bufio.NewScanner(stderrPipe)

	if err := cmd.Start(); err != nil {
		return "", errors.Wrap(err, "cannot start process")
	}

	for pipeScanner.Scan() {
		logline(pipeScanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		return "", errors.Wrap(err, "caught error waiting for process")
	}

	outputFile := filepath.Join(options.TargetDir, outFilename)
	return outputFile, nil
}

func (yd *Provider) videoMetadata(ctx context.Context, youtubeVideoID string) (metadataJson, error) {
	cmd := exec.CommandContext(ctx, "python3", "/usr/local/bin/youtube-dl", "--dump-json", "--", youtubeVideoID)
	output, err := cmd.Output()

	if err != nil {
		return metadataJson{}, errors.Wrapf(err, "cannot get metadata from youtube-dl for video %v", youtubeVideoID)
	}

	var jsonData metadataJson

	if err := json.NewDecoder(bytes.NewReader(output)).Decode(&jsonData); err != nil {
		return metadataJson{}, errors.Wrap(err, "unable to decode json")
	}

	return jsonData, nil
}
