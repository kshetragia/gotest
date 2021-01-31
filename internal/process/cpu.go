package process

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

func (hdlr *prochdlr) cpuInfo() (*time.Time, *time.Duration, error) {
	var createTime, exitTime, kernelTime, userTime windows.Filetime

	if err := windows.GetProcessTimes(hdlr.handler, &createTime, &exitTime, &kernelTime, &userTime); err != nil {
		return nil, nil, errors.Wrap(err, "get process times")
	}

	now := time.Now()
	create := time.Unix(0, createTime.Nanoseconds())
	running := now.Sub(create)

	/*
		cpu.exit = time.Unix(0, exitTime.Nanoseconds())
		cpu.kernel = time.Unix(0, kernelTime.Nanoseconds())
		cpu.user = time.Unix(0, userTime.Nanoseconds())
	*/

	return &create, &running, nil
}
