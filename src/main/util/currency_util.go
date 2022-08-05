package util

import (
	"fmt"
	"strconv"
)

func FormatCurrency(amount float64) string {
	return FormatFloat(amount, 2)
}

func FormatFloat(f float64, precision int) string {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", f)
}
