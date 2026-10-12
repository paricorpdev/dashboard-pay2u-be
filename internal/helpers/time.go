package helpers

import "time"

func GetCurrentTime() time.Time {
	return time.Now()
}

func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}
