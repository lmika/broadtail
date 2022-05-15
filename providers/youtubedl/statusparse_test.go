package youtubedl

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseETA(t *testing.T) {
	scenarios := []struct {
		etaString string
		expected  time.Duration
	}{
		{etaString: "0:15", expected: 15 * time.Second},
		{etaString: "00:32", expected: 32 * time.Second},
		{etaString: "1:49", expected: 1*time.Minute + 49*time.Second},
		{etaString: "23:45", expected: 23*time.Minute + 45*time.Second},
		{etaString: "2:03:05", expected: 2*time.Hour + 3*time.Minute + 5*time.Second},
		{etaString: "4:19:31", expected: 4*time.Hour + 19*time.Minute + 31*time.Second},
		{etaString: "18:00:00", expected: 18 * time.Hour},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.etaString, func(t *testing.T) {
			assert.Equal(t, scenario.expected, parseETA(scenario.etaString))
		})
	}
}
