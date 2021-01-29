package process

import (
	"fmt"
)

type Info struct {
	Name   string
	Path   string
	User   string
	Domain string
	PID    uint32
	PPID   uint32
	LUID   LUID
	Rusage
}

func Test() error {
	return nil
}

func (info *Info) Show() {
	e := info
	user := ""
	if e.User != "" {
		user = "\\\\" + e.Domain + "\\" + e.User
	}
	fmt.Printf("[%v] %v \n", e.PID, e.Name)
	fmt.Printf("\tLUID: %v\n", e.LUID)
	fmt.Printf("\tUser: %v\n", user)
}
