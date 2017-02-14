package util

import (
	"time"

	"github.com/dustin/go-humanize"
)

func HumanizeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	} else {
		return humanize.Time(t)
	}
}
