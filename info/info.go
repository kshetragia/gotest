package info

import (
	"gotest/process"
	"gotest/procnet"
	"gotest/users"
)

// FullInfo holds full unsorted process list Info
type FullInfo []*Info

// Info has complete data about a single process
type Info struct {
	Name      string `json:"Name"`
	Path      string `json:"Path"`
	PID       uint32 `json:"PID"`
	PPID      uint32 `json:"PPID"`
	StartTime string `json:"StartTime"`

	User       *users.User         `json:"Owner"`
	CPUTime    *process.CPUTime    `json:"CPU"`
	MemoryInfo *process.MemoryInfo `json:"Memory"`
	NetInfo    []*procnet.NetInfo  `json:"Net"`
}
