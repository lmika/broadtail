package youtubedl

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type progress struct {
	Percent float64
	Size    string
	Rate    string
	ETA     time.Duration
}

func parseProgress(message string) (progress, bool) {
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
		ETA:     parseETA(groups[4]),
	}, true
}

func parseETA(eta string) (dur time.Duration) {
	if eta == "" {
		return time.Duration(-1)
	}
	
	toks := strings.Split(eta, ":")
	if len(toks) == 0 || len(toks) > 3 {
		return time.Duration(-1)
	}

	if len(toks) == 3 {
		dur += parseNumAndShift(&toks, time.Hour)
	}
	if len(toks) == 2 {
		dur += parseNumAndShift(&toks, time.Minute)
	}
	dur += parseNumAndShift(&toks, time.Second)
	return dur
}

func parseNumAndShift(toks *[]string, mup time.Duration) time.Duration {
	num, _ := strconv.ParseInt((*toks)[0], 10, 32)
	*toks = (*toks)[1:]
	return time.Duration(num) * mup
}

// [download] 2.1% of 86.31MiB at 84.91KiB/s ETA 16:59
var progressRegexp = regexp.MustCompile(`\[download\]\s+([0-9.]+)% of ([0-9A-Za-z.]+) at ([0-9A-Za-z.]+)/s ETA ([0-9:.]+)`)
