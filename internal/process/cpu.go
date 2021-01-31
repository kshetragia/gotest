package process

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

func (hdlr *prochdlr) cpuInfo() (*time.Time, *time.Duration, float64, float64, error) {
	var createTime, exitTime, kernelTime, userTime windows.Filetime

	if err := windows.GetProcessTimes(hdlr.handler, &createTime, &exitTime, &kernelTime, &userTime); err != nil {
		return nil, nil, 0, 0, errors.Wrap(err, "get process times")
	}

	now := time.Now()
	create := time.Unix(0, createTime.Nanoseconds())
	running := now.Sub(create)

	user := float64(userTime.HighDateTime)*429.4967296 + float64(userTime.LowDateTime)*1e-7
	kernel := float64(kernelTime.HighDateTime)*429.4967296 + float64(kernelTime.LowDateTime)*1e-7

	return &create, &running, kernel, user, nil
}
