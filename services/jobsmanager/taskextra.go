package jobsmanager

type VideoDownloadTask interface {
	VideoExtID() string
	VideoTitle() string
}
