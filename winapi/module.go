// +build windows

package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procModule32FirstW = kernel32.NewProc("Module32FirstW")
	procModule32NextW  = kernel32.NewProc("Module32NextW")
)

const (
	// MaxModuleName32 is a buffer size for module name string
	MaxModuleName32 = 255
)

// ModuleEntry32 is analgue of MODULEENTRY32 of Windows API
// See: https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/ns-tlhelp32-moduleentry32
type ModuleEntry32 struct {
	Size          uint32                      // dwSize
	ModuleID      uint32                      // th32ModuleID
	ProcessID     uint32                      // th32ProcessID
	GlblcntUsage  uint32                      // GlblcntUsage
	ProccntUsage  uint32                      // ProccntUsage
	modBaseAddr   *uint8                      // *modBaseAddr
	ModBaseSize   uint32                      // modBaseSize
	ModuleHandler windows.Handle              // hModule
	ModuleName    [MaxModuleName32 + 1]uint16 // szModule
	ExePath       [windows.MAX_PATH]uint16    // szExePath
}

// Module32First retrieves information about the first module associated with a process.
func Module32First(snapshot windows.Handle, lpme *ModuleEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procModule32FirstW.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(lpme)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

// Module32Next retrieves information about the next module associated with a process.
func Module32Next(snapshot windows.Handle, lpme *ModuleEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procModule32NextW.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(lpme)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
