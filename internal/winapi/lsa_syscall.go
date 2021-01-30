// +build windows

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procLsaGetLogonSessionData = secur32.NewProc("LsaGetLogonSessionData")
	procLsaNtStatusToWinError  = advapi32.NewProc("LsaNtStatusToWinError")
)

// LsaGetLogonSessionData function retrieves information about
// a specified logon session. To retrieve information about a logon session,
// the caller must be the owner of the session or a local system administrator.
func LsaGetLogonSessionData(luid *windows.LUID, sessionData **SecurityLogonSessionData) error {
	r0, _, _ := syscall.Syscall(procLsaGetLogonSessionData.Addr(), 2, uintptr(unsafe.Pointer(luid)), uintptr(unsafe.Pointer(sessionData)), 0)
	return LsaNtStatusToWinError(r0)
}

// LsaNtStatusToWinError
func LsaNtStatusToWinError(status uintptr) error {
	r0, _, errno := syscall.Syscall(procLsaNtStatusToWinError.Addr(), 1, status, 0, 0)
	switch errno {
	case windows.ERROR_SUCCESS:
		if r0 == 0 {
			return nil
		}
	case windows.ERROR_MR_MID_NOT_FOUND:
		return fmt.Errorf("Unknown LSA NTSTATUS code %x", status)
	}
	return syscall.Errno(r0)
}
