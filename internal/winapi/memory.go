// +build windows

package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procGetProcessMemoryInfo = psapi.NewProc("GetProcessMemoryInfo")
)

// ProcessMemoryCountersEx is a copy of PROCESS_MEMORY_COUNTERS windows structure
// See also: https://docs.microsoft.com/en-us/windows/win32/api/psapi/ns-psapi-process_memory_counters
type ProcessMemoryCountersEx struct {
	Cb                         uint32  // DWORD
	PageFaultCount             uint32  // DWORD
	PeakWorkingSetSize         uintptr // SIZE_T
	WorkingSetSize             uintptr // SIZE_T
	QuotaPeakPagedPoolUsage    uintptr // SIZE_T
	QuotaPagedPoolUsage        uintptr // SIZE_T
	QuotaPeakNonPagedPoolUsage uintptr // SIZE_T
	QuotaNonPagedPoolUsage     uintptr // SIZE_T
	PagefileUsage              uintptr // SIZE_T
	PeakPagefileUsage          uintptr // SIZE_T
	PrivateUsage               uintptr // SIZE_T
}

// GetProcessMemoryInfo retrieves information about the memory usage of the specified process.
// See also: https://docs.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getprocessmemoryinfo
func GetProcessMemoryInfo(handle windows.Handle, counters *ProcessMemoryCountersEx, size uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetProcessMemoryInfo.Addr(), 3, uintptr(handle), uintptr(unsafe.Pointer(counters)), uintptr(size))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
