// Copyright (c) 2020 Carl Pettersson
//
// This code has been taken from
//    http://github.com/carlpet/winlsa
//
// Big thanks to Carl Pettersson for the helpful example of win API interaction

package proc

import (
	"fmt"
	"syscall"
	"unsafe"

	win "golang.org/x/sys/windows"
)

var (
	secur32                       = win.NewLazySystemDLL("Secur32.dll")
	advapi32                      = win.NewLazySystemDLL("Advapi32.dll")
	procLsaEnumerateLogonSessions = secur32.NewProc("LsaEnumerateLogonSessions")
	procLsaGetLogonSessionData    = secur32.NewProc("LsaGetLogonSessionData")
	procLsaFreeReturnBuffer       = secur32.NewProc("LsaFreeReturnBuffer")
	procLsaNtStatusToWinError     = advapi32.NewProc("LsaNtStatusToWinError")
)

type LSA_LAST_INTER_LOGON_INFO struct {
	LastSuccessfulLogon                        uint64
	LastFailedLogon                            uint64
	FailedAttemptCountSinceLastSuccessfulLogon uint32
}

type SECURITY_LOGON_SESSION_DATA struct {
	Size                  uint32
	LogonId               win.LUID
	UserName              LSA_UNICODE_STRING
	LogonDomain           LSA_UNICODE_STRING
	AuthenticationPackage LSA_UNICODE_STRING
	LogonType             uint32
	Session               uint32
	Sid                   *win.SID
	LogonTime             uint64
	LogonServer           LSA_UNICODE_STRING
	DnsDomainName         LSA_UNICODE_STRING
	Upn                   LSA_UNICODE_STRING
	UserFlags             uint32
	LastLogonInfo         LSA_LAST_INTER_LOGON_INFO
	LogonScript           LSA_UNICODE_STRING
	ProfilePath           LSA_UNICODE_STRING
	HomeDirectory         LSA_UNICODE_STRING
	HomeDirectoryDrive    LSA_UNICODE_STRING
	LogoffTime            uint64
	KickOffTime           uint64
	PasswordLastSet       uint64
	PasswordCanChange     uint64
	PasswordMustChange    uint64
}

type LSA_UNICODE_STRING struct {
	Length        uint16
	MaximumLength uint16
	Buffer        uintptr
}

func LsaEnumerateLogonSessions(sessionCount *uint32, sessions *uintptr) error {
	r0, _, _ := syscall.Syscall(procLsaEnumerateLogonSessions.Addr(), 2, uintptr(unsafe.Pointer(sessionCount)), uintptr(unsafe.Pointer(sessions)), 0)
	return LsaNtStatusToWinError(r0)
}

func LsaGetLogonSessionData(luid *win.LUID, sessionData **SECURITY_LOGON_SESSION_DATA) error {
	r0, _, _ := syscall.Syscall(procLsaGetLogonSessionData.Addr(), 2, uintptr(unsafe.Pointer(luid)), uintptr(unsafe.Pointer(sessionData)), 0)
	return LsaNtStatusToWinError(r0)
}

func LsaFreeReturnBuffer(buffer uintptr) error {
	r0, _, _ := syscall.Syscall(procLsaFreeReturnBuffer.Addr(), 1, buffer, 0, 0)
	return LsaNtStatusToWinError(r0)
}

func LsaNtStatusToWinError(ntstatus uintptr) error {
	r0, _, errno := syscall.Syscall(procLsaNtStatusToWinError.Addr(), 1, ntstatus, 0, 0)
	switch errno {
	case win.ERROR_SUCCESS:
		if r0 == 0 {
			return nil
		}
	case win.ERROR_MR_MID_NOT_FOUND:
		return fmt.Errorf("Unknown LSA NTSTATUS code %x", ntstatus)
	}
	return syscall.Errno(r0)
}
