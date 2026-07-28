package output

import (
	"fmt"
	"os"
)

func Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Warn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Warning: "+format+"\n", args...)
}

func Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
