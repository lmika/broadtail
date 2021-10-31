package youtubedl

import (
	"time"
)

type metadataJson struct {
	UploadDateStr   string       `json:"upload_date"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	ThumbnailURL string       `json:"thumbnail"`
	Duration     int          `json:"duration"`
	//Formats      []formatJson `json:"formats"`
}

func (r metadataJson) UploadDate() (time.Time, error) {
	return time.Parse("20060102", r.UploadDateStr)
}

//func (r formatJson) extractUrlExpiryTimestamp() (time.Time, error) {
//	parsedUrl, err := url.Parse(r.URL)
//	if err != nil {
//		return time.Time{}, errors.Wrapf(err, "unable to parse URL: %v", r.URL)
//	}
//
//	expiryTimestampStr := parsedUrl.Query().Get("expire")
//	if expiryTimestampStr == "" {
//		return time.Time{}, errors.Wrap(err, "no expiry timestamp present")
//	}
//
//	unixSec, err := strconv.ParseInt(expiryTimestampStr, 10, 64)
//	if err != nil {
//		return time.Time{}, errors.Wrapf(err, "expiry is not an integer: %v", expiryTimestampStr)
//	}
//
//	return time.Unix(unixSec, 0).In(time.UTC), nil
//}
