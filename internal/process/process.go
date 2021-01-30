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
	SID              windows.SID
}

type Info struct {
	Name string
	Path string
	PID  uint32
	PPID uint32
	User
	Rusage
}

func Test() error {
	return nil
}

func (info *Info) Show() {
	e := info
	user := "\\\\" + e.User.Domain + "\\" + e.User.Name
	fmt.Printf("[%v] %v \n", e.PID, e.Name)
	fmt.Printf("\tUser: %v\n", user)
	fmt.Printf("\tLUID: %v\n", e.User.AuthenticationId)
}
