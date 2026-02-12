package ui

import (
	"fmt"
	"os"
)

func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "❌ launch failed: %v\n", err)
}
