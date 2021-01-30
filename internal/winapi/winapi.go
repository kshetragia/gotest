// +build windows

package winapi

var (
	secur32  = windows.NewLazySystemDLL("Secur32.dll")
	advapi32 = windows.NewLazySystemDLL("Advapi32.dll")
)
