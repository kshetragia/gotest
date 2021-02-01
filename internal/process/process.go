// +build windows

package process

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows"
)

type User struct {
	Name             string
	Domain           string
	SID              string
	SessionID        uint32
	LastSuccessLogon time.Time
	AuthenticationID windows.LUID
}

type CPUTime struct {
	Kernel float64
	User   float64
	System float64
	Total  string
}

type Info struct {
	Name        string
	Path        string
	PID         uint32
	PPID        uint32
	StartTime   string
	MemoryUsage uint64
	CPUTime
	User
}

func (info *Info) Show() {
	e := info
	user := "\\\\" + e.User.Domain + "\\" + e.User.Name
	fmt.Printf("[%v] %v \n", e.PID, e.Name)
	fmt.Printf("\tPath: %v\n", e.Path)
	fmt.Printf("\tExecution time: %v\n", e.StartTime)
	fmt.Printf("\tUser: %v\n", user)
	fmt.Printf("\t  SessionId: %v\n", e.User.SessionID)
	fmt.Printf("\t  Last Login: %v\n", e.User.LastSuccessLogon.String())
	fmt.Printf("\t  SID: %v\n", e.User.SID)
	fmt.Printf("\tCPU:\n")
	fmt.Printf("\t  Kernel: %.2fs\n", e.CPUTime.Kernel)
	fmt.Printf("\t  User: %.2fs\n", e.CPUTime.User)
	fmt.Printf("\t  System: %.2f%%\n", e.CPUTime.System)
	fmt.Printf("\t  Total:  %v\n", e.CPUTime.Total)
	fmt.Printf("\tRAM: %v Mb [%v]\n", float64(e.MemoryUsage)/1024/1024, e.MemoryUsage)
	fmt.Printf("\tLUID: %v\n", e.User.AuthenticationID)
	fmt.Println()
}
