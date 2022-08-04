package util

import (
	"fmt"
	"math"
	"time"
)

func pluralise(str string, i int64) string {
	sMaybe := "s"
	if i == 1 {
		sMaybe = ""
	}
	return fmt.Sprintf("%s%s", str, sMaybe)
}

func PrettyDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		seconds := int64(duration.Seconds())
		return fmt.Sprintf("%d %s", seconds, pluralise("second", seconds))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		minutes := int64(duration.Minutes())
		seconds := int64(remainingSeconds)
		return fmt.Sprintf("%d %s %d %s", minutes, pluralise("minute", minutes), seconds, pluralise("second", seconds))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		hours := int64(duration.Hours())
		minutes := int64(remainingMinutes)
		seconds := int64(remainingSeconds)
		return fmt.Sprintf("%d %s %d %s %d %s",
			hours, pluralise("hour", hours), minutes, pluralise("minute", minutes), seconds, pluralise("second", seconds))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	days := int64(duration.Hours() / 24)
	hours := int64(remainingHours)
	minutes := int64(remainingMinutes)
	seconds := int64(remainingSeconds)
	return fmt.Sprintf("%d %s %d %s %d %s %d %s",
		days, pluralise("day", days), hours, pluralise("hour", hours),
		minutes, pluralise("minute", minutes), seconds, pluralise("second", seconds))
}
