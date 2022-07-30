package util

import "fmt"

func Pluralise(str string, i int) string {
	sMaybe := ""
	if i > 1 {
		sMaybe = "s"
	}
	return fmt.Sprintf("%s%s", str, sMaybe)
}
