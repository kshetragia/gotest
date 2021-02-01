package main

import (
	"fmt"
	"gotest/info"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func main() {
	pinfo, err := info.Collect()
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
	pinfo.Show()
}
