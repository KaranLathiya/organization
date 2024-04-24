package utils

import "time"

func AddMinutesToCurrentUTCTime(addMinutes int) string {
	currentTime := time.Now().UTC().Add(time.Minute * time.Duration(addMinutes))
	outputFormat := "2006-01-02 15:04:05-07:00"
	currentFormattedTime := currentTime.Format(outputFormat)
	return currentFormattedTime
}
