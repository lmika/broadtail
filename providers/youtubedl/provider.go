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
	youtubeDLCommand []string
}

func New(youtubeDLCommand []string) (*Provider, error) {
	provider := &Provider{
		youtubeDLCommand: youtubeDLCommand,
	}
	if err := provider.checkAvailable(); err != nil {
		return nil, err
	}

	return provider, nil
}

func (p *Provider) checkAvailable() error {
	cmd := p.buildYoutubeDLCommand(context.Background(), "--version")
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
		VideoRef:     models.VideoRef{Source: models.YoutubeVideoRefSource, ID: youtubeId},
		Title:        jsonData.Title,
		ChannelID:    jsonData.ChannelID,
		ChannelName:  jsonData.Channel,
		Description:  jsonData.Description,
		ThumbnailURL: jsonData.ThumbnailURL,
		UploadedOn:   uploadDate,
		Duration:     jsonData.Duration,
	}, nil
}

// "python3", "/usr/local/bin/youtube-dl"

func (p *Provider) DownloadVideo(ctx context.Context, youtubeId string, options models.DownloadOptions, logline func(line string)) (string, error) {
	const filenameFormat = "%(title)s.%(id)s.%(ext)s"

	// Get the expected filename
	downloadUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%v", youtubeId)
	filenameCmd := p.buildYoutubeDLCommand(ctx, "-f", "mp4[height<=720]",
		"--get-filename", "-o", filenameFormat, "--restrict-filenames", downloadUrl)
	out, err := filenameCmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "unable to determine target filename")
	}
	outFilename := strings.TrimSpace(string(out))

	logline("Target dir: " + options.TargetDir)
	cmd := p.buildYoutubeDLCommand(ctx, "-o", filenameFormat, "--restrict-filenames",
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
	cmd := yd.buildYoutubeDLCommand(ctx, "--dump-json", "--", youtubeVideoID)
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

func (yt *Provider) buildYoutubeDLCommand(ctx context.Context, args ...string) *exec.Cmd {
	fullCmd := append(append([]string{}, yt.youtubeDLCommand...), args...)

	return exec.CommandContext(ctx, fullCmd[0], fullCmd[1:]...)
}
