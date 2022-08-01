package util

import (
	"fmt"
	"math"
	"strings"
)

func ClearConsoleLines(count int) {
	fmt.Print(strings.Repeat("\033[1A\033[K", count))
}

func ProgressBar(completedRatio float64, progressBarSize int) string {
	progressBars := int(math.Round(completedRatio * float64(progressBarSize)))
	completedBars := fmt.Sprintf("%s>", strings.Repeat("=", int(math.Max(float64(progressBars-1), 0))))
	remainingBars := strings.Repeat(" ", progressBarSize-progressBars)
	return fmt.Sprintf("[%s%s]", completedBars, remainingBars)
}
