package process

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

type CPUTime struct {
	Kernel float64
	User   float64
	System float64
	Total  string
}

// cpuInfo is using windows GetProcessTimes() function to get CPU process time
// See also: https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getprocesstimes
func (hdlr *prochdlr) cpuInfo() (string, CPUTime, error) {
	var createTime, exitTime, kernelTime, userTime windows.Filetime
	var cpu CPUTime

	if err := windows.GetProcessTimes(hdlr.handler, &createTime, &exitTime, &kernelTime, &userTime); err != nil {
		return "", cpu, errors.Wrap(err, "get process times")
	}

	now := time.Now()
	create := time.Unix(0, createTime.Nanoseconds())

	// Total CPU time
	total := now.Sub(create)

	// The Filetime structure has two uint32 parts of uint64 time number ticked by 100ns intervals
	// So we should do the follow to get real time from time chunks:
	//     time(seconds) = (high << 32 + low) * 1e-9 * 100
	cpu.User = float64(uint64(userTime.HighDateTime)<<32+uint64(userTime.LowDateTime)) * 1e-7
	cpu.Kernel = float64(uint64(kernelTime.HighDateTime)<<32+uint64(kernelTime.LowDateTime)) * 1e-7
	cpu.System = (cpu.User + cpu.Kernel) * 100 / total.Seconds()
	cpu.Total = total.Round(time.Millisecond).String()

	return create.String(), cpu, nil
}
