// +build windows

package winapi

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Get Logon session data
// See also:
//     https://docs.microsoft.com/en-us/windows/win32/api/ntsecapi/ns-ntsecapi-security_logon_session_data

// LsaUnicodeString structure is used by various
// Local Security Authority (LSA) functions to specify a Unicode string.
type LsaUnicodeString struct {
	Length        uint16
	MaximumLength uint16
	Buffer        uintptr
}

// LsaLastInterLogonInfo structure contains information about a logon session.
type LsaLastInterLogonInfo struct {
	LastSuccessfulLogin                        uint64
	LastFailedLogin                            uint64
	FailedAttemptCountSinceLastSuccessfulLogin uint32
}

// SecurityLogonSessionData structure contains information about a logon session.
type SecurityLogonSessionData struct {
	Size                  uint32
	LogonID               windows.LUID
	UserName              LsaUnicodeString
	LogonDomain           LsaUnicodeString
	AuthenticationPackage LsaUnicodeString
	LogonType             uint32
	Session               uint32
	Sid                   *windows.SID
	LogonTime             uint64
	LogonServer           LsaUnicodeString
	DNSDomainName         LsaUnicodeString
	Upn                   LsaUnicodeString
	UserFlags             uint32
	LastLogonInfo         LsaLastInterLogonInfo
	LogonScript           LsaUnicodeString
	ProfilePath           LsaUnicodeString
	HomeDirectory         LsaUnicodeString
	HomeDirectoryDrive    LsaUnicodeString
	LogoffTime            uint64
	KickOffTime           uint64
	PasswordLastSet       uint64
	PasswordCanChange     uint64
	PasswordMustChange    uint64
}

// LsaUnicodeToString converts UCS-16 into UTF-8 string
func LsaUnicodeToString(str LsaUnicodeString) string {
	if str.Buffer == 0 || str.Length == 0 {
		return ""
	}
	return windows.UTF16ToString((*[4096]uint16)(unsafe.Pointer(str.Buffer))[:str.Length])
}
