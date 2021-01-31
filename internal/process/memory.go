// +build windows

package process

import (
	"gotest/winapi"
	"unsafe"

	"github.com/pkg/errors"
)

// memInfo returns the total amount of private memory
// that the memory manager has committed for a running process.
func (hdlr *prochdlr) memInfo() (uint64, error) {
	counters := winapi.ProcessMemoryCountersEx{}
	size := uint32(unsafe.Sizeof(counters))
	err := winapi.GetProcessMemoryInfo(hdlr.handler, &counters, size)
	if err != nil {
		return 0, errors.Wrap(err, "get process memory info")
	}

	return uint64(counters.PrivateUsage), nil
}
