package main

import (
	"fmt"
	"os"
)

func exit(code int, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(code)
}

func main() {
	args := ParseArgs()
	_ = args

	exit(0, nil)
}
