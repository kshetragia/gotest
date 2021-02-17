// +build windows

package winapi

import "syscall"

const (
	errnoERRORIOPENDING = 997
)

var (
	errERRORIOPENDING error = syscall.Errno(errnoERRORIOPENDING)
)

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERRORIOPENDING:
		return errERRORIOPENDING
	}
	return e
}
