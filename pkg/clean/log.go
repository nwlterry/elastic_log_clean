package clean

import (
	"fmt"
	"os"
)

// Verbose enables per-file process output on stderr.
var Verbose bool

func Logf(format string, args ...any) {
	if !Verbose {
		return
	}
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
