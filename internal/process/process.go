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

type CPU struct {
	Kernel float64
	User   float64
}

type Info struct {
	Name        string
	Path        string
	PID         uint32
	PPID        uint32
	MemoryUsage uint64
	StartTime   *time.Time
	Running     *time.Duration
	User
	CPU
}

func (info *Info) Show() {
	e := info
	user := "\\\\" + e.User.Domain + "\\" + e.User.Name
	fmt.Printf("[%v] %v \n", e.PID, e.Name)
	fmt.Printf("\tPath: %v\n", e.Path)
	fmt.Printf("\tExecution time: %v (%v running)\n", e.StartTime.String(), e.Running.String())
	fmt.Printf("\tUser: %v\n", user)
	fmt.Printf("\t  SessionId: %v\n", e.User.SessionID)
	fmt.Printf("\t  Last Login: %v\n", e.User.LastSuccessLogon.String())
	fmt.Printf("\t  SID: %v\n", e.User.SID)
	fmt.Printf("\tCPU:\n")
	fmt.Printf("\t  Kernel: %v\n", e.CPU.Kernel)
	fmt.Printf("\t  User: %v\n", e.CPU.User)
	fmt.Printf("\tRAM: %v Mb [%v]\n", float64(e.MemoryUsage)/1024/1024, e.MemoryUsage)
	fmt.Printf("\tLUID: %v\n", e.User.AuthenticationID)
	fmt.Println()
}
