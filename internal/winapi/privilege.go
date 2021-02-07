// +build windows

package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	PRIVILEGE_SET_ALL_NECESSARY uint32 = 1
)

const (
	SE_PRIVILEGE_ENABLED_BY_DEFAULT uint32 = 1 << iota
	SE_PRIVILEGE_ENABLED
	SE_PRIVILEGE_REMOVED

//	SE_PRIVILEGE_USED_FOR_ACCESS uint64 = 8 << 8
)

var (
	procPrivilegeCheck = advapi32.NewProc("PrivilegeCheck")
)

// PrivilegeSet is a PRIVILEGE_SET windows structure specifies a set of privileges.
// See also https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-privilege_set
type PrivilegeSet struct {
	PrivilegeCount uint32
	Control        uint32
	Privilege      [1]windows.LUIDAndAttributes // ANYSIZE_ARRAY len defined in winnt.h
}

// PrivilegeCheck function determines whether a specified set of privileges are enabled in an access token.
// See also https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-privilegecheck
func PrivilegeCheck(token windows.Token, set *PrivilegeSet, res *bool) (err error) {
	r1, _, e1 := syscall.Syscall(procPrivilegeCheck.Addr(), 3, uintptr(token), uintptr(unsafe.Pointer(set)), uintptr(unsafe.Pointer(res)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
