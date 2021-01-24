package proc

import (
	"time"

	"github.com/pkg/errors"
	win "golang.org/x/sys/windows"
)

type CPUInfo struct {
	creation time.Time
	running  time.Duration
	exit     int64
	kernel   int64
	user     int64
}

func (cpu *CPUInfo) Read(handler win.Handle) error {
	var createTime, exitTime, kernelTime, userTime win.Filetime

	if err := win.GetProcessTimes(handler, &createTime, &exitTime, &kernelTime, &userTime); err != nil {
		errors.Wrap(err, "get CPU usage")
		return err
	}

	now := time.Now()
	cpu.creation = time.Unix(0, createTime.Nanoseconds())
	cpu.running = now.Sub(cpu.creation)

	/*
		cpu.creation = time.Unix(creation, 0).String()
		cpu.running = time.Unix(running, 0).String()
	*/

	/*
		cpu.exit = time.Unix(0, exitTime.Nanoseconds())
		cpu.kernel = time.Unix(0, kernelTime.Nanoseconds())
		cpu.user = time.Unix(0, userTime.Nanoseconds())
	*/

	return nil
}
