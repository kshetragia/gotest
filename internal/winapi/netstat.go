// +build windows

package winapi

import (
	"syscall"
	"unsafe"
)

var (
	procGetPerTcpConnectionEStats = iphlpapi.NewProc("GetPerTcpConnectionEStats")
)

// GetPerTcpConnectionEStats function retrieves extended statistics for an IPv4 TCP connection.
// See also: https://docs.microsoft.com/en-us/windows/win32/api/iphlpapi/nf-iphlpapi-getpertcpconnectionestats
func GetPerTcpConnectionEStats(row *MIB_TCPROW, estatsType TCP_ESTATS_TYPE, rw *byte, rwVersion uint64, rwSize uint64, ros *byte,
	rosVersion uint64, rosSize uint64, rod *byte, rodVersion uint64, rodSize uint64) (err error) {

	r1, _, e1 := syscall.Syscall12(procGetProcessMemoryInfo.Addr(), 11, uintptr(unsafe.Pointer(row)), uintptr(estatsType),
		uintptr(unsafe.Pointer(rw)), uintptr(rwVersion), uintptr(rwSize), uintptr(unsafe.Pointer(ros)),
		uintptr(rosVersion), uintptr(rosSize), uintptr(unsafe.Pointer(rod)), uintptr(rodVersion), uintptr(rodSize), 0)

	if r1 == 0 {
		err = errnoErr(e1)
	}

	return
}

// GetPerTcp6ConnectionEStats function retrieves extended statistics for an IPv4 TCP connection.
// See also: https://docs.microsoft.com/en-us/windows/win32/api/iphlpapi/nf-iphlpapi-getpertcp6connectionestats
func GetPerTcp6ConnectionEStats(row *MIB_TCP6ROW, estatsType TCP_ESTATS_TYPE, rw *byte, rwVersion uint64, rwSize uint64, ros *byte,
	rosVersion uint64, rosSize uint64, rod *byte, rodVersion uint64, rodSize uint64) (err error) {

	r1, _, e1 := syscall.Syscall12(procGetProcessMemoryInfo.Addr(), 11, uintptr(unsafe.Pointer(row)), uintptr(estatsType),
		uintptr(unsafe.Pointer(rw)), uintptr(rwVersion), uintptr(rwSize), uintptr(unsafe.Pointer(ros)),
		uintptr(rosVersion), uintptr(rosSize), uintptr(unsafe.Pointer(rod)), uintptr(rodVersion), uintptr(rodSize), 0)

	if r1 == 0 {
		err = errnoErr(e1)
	}

	return
}
