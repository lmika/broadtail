package ytdownload

import (
	"regexp"
	"strconv"
)

type progress struct {
	Percent float64
	Size    string
	Rate    string
	ETA     string
}

func parseProgress(message string) (progress, bool) {
	//log.Printf("Parsing progress: '%v'", message)
	groups := progressRegexp.FindStringSubmatch(message)
	if len(groups) != 5 {
		//log.Println("bad progress: ", len(groups))
		return progress{}, false
	}

	percentFloat, err := strconv.ParseFloat(groups[1], 64)
	if err != nil {
		return progress{}, false
	}

	return progress{
		Percent: percentFloat,
		Size:    groups[2],
		Rate:    groups[3],
		ETA:     groups[4],
	}, true
}

// [download] 2.1% of 86.31MiB at 84.91KiB/s ETA 16:59
var progressRegexp = regexp.MustCompile(`\[download\]\s+([0-9.]+)% of ([0-9A-Za-z.]+) at ([0-9A-Za-z.]+)/s ETA ([0-9:.]+)`)
