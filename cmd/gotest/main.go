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

	pinfo, err := process.Collect()
	for _, info := range *pinfo {
		info.Show()
	}
	if err != nil {
		for e := err; e != nil; e = errors.Unwrap(err) {
			fmt.Println("Error:", e)
		}
		if err, ok := err.(stackTracer); ok {
			for _, f := range err.StackTrace() {
				fmt.Printf("%+s: %d\n", f, f)
			}
		}

	}
}
