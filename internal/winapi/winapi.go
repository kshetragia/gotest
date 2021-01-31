// +build windows

package winapi

import (
	"golang.org/x/sys/windows"
)

var (
	advapi32 = windows.NewLazySystemDLL("Advapi32.dll")
	kernel32 = windows.NewLazySystemDLL("Kernel32.dll")
	psapi    = windows.NewLazySystemDLL("Psapi.dll")
	secur32  = windows.NewLazySystemDLL("Secur32.dll")
)
