package jobs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateHistory_PushAndList(t *testing.T) {

	t.Run("not yet cycling", func(t *testing.T) {
		uh := newUpdateHistory(5)

		uh.push(Update{Status: "Line 1"})
		uh.push(Update{Status: "Line 2"})
		uh.push(Update{Status: "Line 3"})

		assert.Equal(t, []Update{
			{Status: "Line 1"},
			{Status: "Line 2"},
			{Status: "Line 3"},
		}, uh.list())
	})

	t.Run("single cycle", func(t *testing.T) {
		uh := newUpdateHistory(5)

		for i := 1; i <= 8; i++ {
			uh.push(Update{Status: fmt.Sprintf("Line %d", i)})
		}

		assert.Equal(t, []Update{
			{Status: "Line 4"},
			{Status: "Line 5"},
			{Status: "Line 6"},
			{Status: "Line 7"},
			{Status: "Line 8"},
		}, uh.list())
	})

	t.Run("multiple cycles", func(t *testing.T) {
		uh := newUpdateHistory(7)

		for i := 1; i <= 49; i++ {
			uh.push(Update{Status: fmt.Sprintf("Line %d", i)})
		}

		assert.Equal(t, []Update{
			{Status: "Line 43"},
			{Status: "Line 44"},
			{Status: "Line 45"},
			{Status: "Line 46"},
			{Status: "Line 47"},
			{Status: "Line 48"},
			{Status: "Line 49"},
		}, uh.list())
	})
}
