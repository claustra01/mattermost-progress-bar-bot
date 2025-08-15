package date

import (
	"time"
)

const (
	// "YYYY/MM/DD"
	StartDate = "2025/06/02"
	EndDate   = "2026/02/28"
)

func StrToDate(str string) time.Time {
	t, err := time.Parse("2006/01/02", str)
	if err != nil {
		panic("Invalid date format")
	}
	return t
}

func GetProgress(now time.Time) float64 {
	startDate := StrToDate(StartDate)
	endDate := StrToDate(EndDate)
	progress := float64(now.Sub(startDate)) / float64(endDate.Sub(startDate))
	return progress
}

func GetRemainingDays(now time.Time) int {
	endDate := StrToDate(EndDate)
	remainingDays := int(endDate.Sub(now).Hours() / 24)
	return remainingDays
}
