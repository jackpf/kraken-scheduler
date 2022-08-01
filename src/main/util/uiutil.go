package util

import (
	"fmt"
	"math"
	"strings"
)

func PadLine(str string, length int) string {
	return fmt.Sprintf("%s%s", str, strings.Repeat(" ", length-len(str)))
}

func ClearConsoleLines(count int) {
	fmt.Print(strings.Repeat("\033[1A\033[K", count))
}

func ProgressBar(completedRatio float64, progressBarSize int) string {
	progressBars := int(math.Round(completedRatio * float64(progressBarSize)))
	completedBars := fmt.Sprintf("%s>", strings.Repeat("=", int(math.Max(float64(progressBars-1), 0))))
	minus := 0
	if progressBars == 0 {
		minus = 1
	}
	remainingBars := strings.Repeat(" ", progressBarSize-progressBars-minus)
	return fmt.Sprintf("[%s%s]", completedBars, remainingBars)
}
