package logger

import "fmt"

var Enabled = false

func Log(message string, args ...any) {
	if Enabled {
		fmt.Printf(message+"\n", args...)
	}
}
