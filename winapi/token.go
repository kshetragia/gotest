// +build windows

package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procDuplicateToken = advapi32.NewProc("DuplicateToken")
)

// DuplicateToken function creates a new access token that duplicates one already in existence.
// https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-duplicatetoken
func DuplicateToken(token windows.Token, securityLevel uint32, newToken *windows.Token) (err error) {
	r1, _, e1 := syscall.Syscall(procDuplicateToken.Addr(), 3, uintptr(token), uintptr(securityLevel), uintptr(unsafe.Pointer(newToken)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
