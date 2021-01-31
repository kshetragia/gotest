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
	AuthenticationID windows.LUID
	SessionID        uint32
	SID              string
	LastSuccessLogon time.Time
}

type Info struct {
	Name      string
	Path      string
	PID       uint32
	PPID      uint32
	StartTime *time.Time
	Running   *time.Duration
	User
	Rusage
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
	fmt.Printf("\tLUID: %v\n", e.User.AuthenticationID)
	fmt.Println()
}
