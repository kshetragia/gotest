package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func main() {
	ostype := runtime.GOOS
	if ostype != "windows" {
		fmt.Println("This OS is not supported")
		os.Exit(-1)
	}
	/*
		info := proc.SysInfo{}

		info.Test()
		err := info.Collect()

		if err, ok := err.(stackTracer); ok {
			for _, f := range err.StackTrace() {
				fmt.Printf("%+s: %d\n", f, f)
			}
		}

		info.Show()
		info.Free()
	*/
}
