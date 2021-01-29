package main

// +build windows

import (
	"fmt"
	"gotest/process"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func main() {

	_, err := process.Collect()

	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s: %d\n", f, f)
		}
	}
}
