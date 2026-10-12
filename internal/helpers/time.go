package helpers

import "time"

func GetCurrentTime() time.Time {
	return time.Now()
}

func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ResolveMonthYear falls back to current month and year if unspecified or invalid.
func ResolveMonthYear(month int, year int) (int, int) {
	now := time.Now()
	if month < 1 || month > 12 {
		month = int(now.Month())
	}
	if year < 2000 || year > 2100 {
		year = now.Year()
	}
	return month, year
}

// GetMonthDateRange calculates start, end, and total days for given year and month.
func GetMonthDateRange(year int, month int, loc *time.Location) (startDate time.Time, endDate time.Time, daysInMonth int) {
	if loc == nil {
		loc = time.Local
	}
	startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	endDate = startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)
	daysInMonth = startDate.AddDate(0, 1, -1).Day()
	return startDate, endDate, daysInMonth
}
