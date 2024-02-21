package utils

import "time"

func GetUTCDate(input time.Time) time.Time {
	input = input.UTC()

	// set the time to 00:00:00
	dateOnly := time.Date(input.Year(), input.Month(), input.Day(), 0, 0, 0, 0, time.UTC)

	return dateOnly
}
