// +build windows

package privilege

import (
	"gotest/winapi"
	"unsafe"

	"golang.org/x/sys/windows"
)

func IsEnabled(priv string) bool {

	luid := windows.LUID{}
	winPriv, _ := windows.UTF16PtrFromString(priv)
	if err := windows.LookupPrivilegeValue(nil, winPriv, &luid); err != nil {
		return false
	}

	handler := windows.CurrentProcess()
	defer windows.CloseHandle(handler)

	var token windows.Token
	if err := windows.OpenProcessToken(handler, windows.TOKEN_QUERY, &token); err != nil {
		return false
	}
	defer token.Close()

	var result bool
	set := winapi.PrivilegeSet{
		Control:        winapi.PRIVILEGE_SET_ALL_NECESSARY,
		PrivilegeCount: 1,
	}
	set.Privilege[0].Luid = luid
	set.Privilege[0].Attributes = winapi.SE_PRIVILEGE_ENABLED

	if err := winapi.PrivilegeCheck(token, &set, &result); err != nil {
		return false
	}

	return result
}

// Set function changes provided privelege status
func Set(priv string, doEnable bool) bool {

	luid := windows.LUID{}
	winPriv, _ := windows.UTF16PtrFromString(priv)
	if err := windows.LookupPrivilegeValue(nil, winPriv, &luid); err != nil {
		return false
	}

	handler := windows.CurrentProcess()
	defer windows.CloseHandle(handler)

	var token windows.Token
	if err := windows.OpenProcessToken(handler, windows.TOKEN_QUERY|windows.TOKEN_ADJUST_PRIVILEGES, &token); err != nil {
		return false
	}
	defer token.Close()

	oldState := windows.Tokenprivileges{}
	newState := windows.Tokenprivileges{PrivilegeCount: 1}
	newState.Privileges[0].Luid = luid

	if doEnable {
		newState.Privileges[0].Attributes = winapi.SE_PRIVILEGE_ENABLED
	} else {
		newState.Privileges[0].Attributes = winapi.SE_PRIVILEGE_REMOVED
	}

	newSize := uint32(unsafe.Sizeof(newState))

	var retSize uint32
	if err := windows.AdjustTokenPrivileges(token, false, &newState, newSize, &oldState, &retSize); err != nil {
		return false
	}

	return true
}
