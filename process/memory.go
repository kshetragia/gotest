// +build windows

package process

import (
	"gotest/winapi"
	"unsafe"

	"github.com/pkg/errors"
)

// MemoryInfo collects memory statistics about process
type MemoryInfo struct {
	WorkingSetSize         uintptr `json:"WorkingSetSize"`
	QuotaPagedPoolUsage    uintptr `json:"QuotaPagedPoolUsage"`
	QuotaNonPagedPoolUsage uintptr `json:"QuotaNonPagedPoolUsage"`
	PrivateUsage           uintptr `json:"PrivateUsage"`
}

// memInfo returns the total amount of private memory
// that the memory manager has committed for a running process.
func (hdlr *prochdlr) memInfo() (*MemoryInfo, error) {
	counters := winapi.ProcessMemoryCountersEx{}
	size := uint32(unsafe.Sizeof(counters))
	err := winapi.GetProcessMemoryInfo(hdlr.handler, &counters, size)
	if err != nil {
		return nil, errors.Wrap(err, "get process memory info")
	}

	minfo := MemoryInfo{
		WorkingSetSize:         counters.WorkingSetSize,
		QuotaPagedPoolUsage:    counters.QuotaPagedPoolUsage,
		QuotaNonPagedPoolUsage: counters.QuotaNonPagedPoolUsage,
		PrivateUsage:           counters.PrivateUsage,
	}

	return &minfo, nil
}
