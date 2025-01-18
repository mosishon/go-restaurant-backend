package helpers

import "time"

func InTimeSpan(startDate, endDate time.Time, now time.Time) bool {
	return now.After(startDate) && now.Before(endDate)
}
