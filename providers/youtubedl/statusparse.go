package youtubedl

import (
		"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/muesli/ansi"
)

type progress struct {
	Percent float64
	Size    string
	Rate    string
	ETA     time.Duration
}

func parseProgress(message string) (progress, bool) {
	msg := strings.TrimSpace(stripAnsiCharacters(message))
	groups := progressRegexp.FindStringSubmatch(msg)
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
// [download] 96.0% of 20.55MiB at 2.04MiB/s ETA 00:00
var progressRegexp = regexp.MustCompile(`\[download\]\s+([0-9.]+)%\s+of\s+([0-9A-Za-z.]+)\s+at\s+([0-9A-Za-z.]+)/s ETA ([0-9:.]+)`)

func stripAnsiCharacters(line string) string {
	nonAnsiSequences := new(strings.Builder)
	inAnsi := false

	for _, c := range line {
		if c == ansi.Marker {
			inAnsi = true
		} else if inAnsi {
			if ansi.IsTerminator(c) {
				inAnsi = false
			}
		} else {
			nonAnsiSequences.WriteRune(c)
		}
	}
	return nonAnsiSequences.String()
}
