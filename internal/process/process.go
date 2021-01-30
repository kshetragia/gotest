package process

import (
	"fmt"

	"golang.org/x/sys/windows"
)

type User struct {
	Name             string
	Domain           string
	AuthenticationId windows.LUID
	SessionId        string
	SID              string
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
	fmt.Printf("\tSID: %v\n", e.User.SID)
	fmt.Printf("\tLUID: %v\n", e.User.AuthenticationId)
}
