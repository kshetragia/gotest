// +build windows

package process

import (
	"golang.org/x/sys/windows"
)

// Process takes information about one executed process
type Process struct {
	CPUTime // struct of cpu usage statistics

	Name        string
	Path        string
	PID         uint32
	PPID        uint32
	StartTime   string
	UserKey     *windows.LUID
	MemoryUsage uint64
}
