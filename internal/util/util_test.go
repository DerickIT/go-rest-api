package util_test

import (
	"testing"
	"time"

	"github.com/derickit/go-rest-api/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestFormatTimeToISO(t *testing.T) {
	got := util.FormatTimeToISO(time.Date(2024, 5, 16, 9, 34, 0, 0, time.UTC))
	want := "2024-05-16T09:34:00Z"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCurrentISOTime(t *testing.T) {
	got := util.CurrentISOTime()
	parsedTime, err := time.Parse(time.RFC3339, got)
	z, offset := parsedTime.Zone()
	if err != nil {
		t.Error("Error parsing time")
	}
	assert.Equal(t, "UTC", z)
	assert.Equal(t, 0, offset)
}

type DevModeTestCase struct {
	input  string
	result bool
}

func TestIsDevMode(t *testing.T) {
	testCases := []DevModeTestCase{
		{input: "local", result: true},
		{input: "dev", result: true},
		{input: "production", result: false},
		{input: "staging", result: false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := util.IsDevMode(tc.input)
			assert.Equal(t, tc.result, got)
		})
	}
}
