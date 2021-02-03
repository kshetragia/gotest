// +build windows

package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procGetProcessInformation = kernel32.NewProc("GetProcessInformation")
)

// GetProcessInformation retrieves information about the specified process
func GetProcessInformation(handle windows.Handle, infoClass uint32, info *byte, infoLen uint32, retLen *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procGetProcessInformation.Addr(), 5, uintptr(handle), uintptr(infoClass), uintptr(unsafe.Pointer(info)), uintptr(infoLen), uintptr(unsafe.Pointer(retLen)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
