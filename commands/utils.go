package commands

import (
	"time"
)

func FormatTime(t time.Time) string {
	return t.Local().Format(time.RFC3339)
}
