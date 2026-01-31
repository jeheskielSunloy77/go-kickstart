package ui

import (
	"fmt"
	"os"
)

func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
