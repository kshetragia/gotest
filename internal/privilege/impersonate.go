// +build windows

package privilege

import (
	"fmt"
	"gotest/winapi"
	"os"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

// Impersonate takes pid of process, steals process rights and use them for the current thread
func Impersonate(pid uint32) (err error) {

	handler, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		fmt.Println(err, "open process handler")
		return
	}

	var token windows.Token
	if err = windows.OpenProcessToken(handler, windows.TOKEN_IMPERSONATE|windows.TOKEN_DUPLICATE, &token); err != nil {
		fmt.Println(err, "open process token")
		return
	}

	// Steal priveleges with SecurityIdentification key (2)
	// https://docs.microsoft.com/en-us/windows/win32/api/winnt/ne-winnt-security_impersonation_level
	var newToken windows.Token
	if err = winapi.DuplicateToken(token, 2, &newToken); err != nil {
		fmt.Println(err, "impersonate process token")
		return
	}

	if err = windows.SetThreadToken(nil, newToken); err != nil {
		fmt.Println(err, "Set thread token")
		return
	}

	return
}

// RevertToSelf flushes impersonated access rights to default user's ones
func RevertToSelf(doPanic bool) (err error) {
	if err := windows.RevertToSelf(); err != nil {
		errors.Wrap(err, "Revert Access Rights")
		if doPanic {
			panic(err)
		}
	}
	return
}

// IsAdmin returns true if program executed was run as Administrator
func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

// PidByName returns process ID by process name
func PidByName(needName string) (pid uint32, err error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err = windows.Process32First(snapshot, &entry); err != nil {
		return
	}

	for {
		name := windows.UTF16ToString(entry.ExeFile[:])
		if name == needName {
			pid = entry.ProcessID
			return
		}
		if err = windows.Process32Next(snapshot, &entry); err != nil {
			return
		}
	}

	return
}
