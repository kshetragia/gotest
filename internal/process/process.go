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
	Name string
	Path string
	PID  uint32
	PPID uint32
	User
	Rusage
}

func (info *Info) Show() {
	e := info
	user := "\\\\" + e.User.Domain + "\\" + e.User.Name
	fmt.Printf("[%v] %v \n", e.PID, e.Name)
	fmt.Printf("\tUser: %v\n", user)
	fmt.Printf("\tSessionId: %v\n", e.User.SessionID)
	fmt.Printf("\tLast Successful Login: %v\n", e.User.LastSuccessLogon.String())
	fmt.Printf("\tSID: %v\n", e.User.SID)
	fmt.Printf("\tLUID: %v\n", e.User.AuthenticationID)
}
