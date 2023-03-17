package after

import (
	"fmt"
	"os"
)

func ErrorExit(code int, message string) {
	_, err := fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
	if err != nil {
		panic(err)
	}
	defer os.Exit(code)
}
