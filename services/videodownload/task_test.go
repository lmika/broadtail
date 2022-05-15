package videodownload

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormatETA(t *testing.T) {
	scenarios := []struct {
		eta      time.Duration
		expected string
	}{
		{eta: 0, expected: ""},
		{eta: -1, expected: ""},
		{eta: 2 * time.Second, expected: "ETA 0:02"},
		{eta: 18 * time.Second, expected: "ETA 0:18"},
		{eta: 2*time.Minute + 5*time.Second, expected: "ETA 2:05"},
		{eta: 21*time.Minute + 18*time.Second, expected: "ETA 21:18"},
		{eta: 41*time.Minute + 9*time.Second, expected: "ETA 41:09"},
		{eta: 1*time.Hour + 2*time.Minute + 3*time.Second, expected: "ETA 1:02:03"},
		{eta: 12*time.Hour + 34*time.Minute + 56*time.Second, expected: "ETA 12:34:56"},
	}
	for i, scenario := range scenarios {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, scenario.expected, formatETA(scenario.eta))
		})
	}
}
