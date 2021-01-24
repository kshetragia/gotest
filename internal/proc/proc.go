package proc

import "container/list"

type SysInfo struct {
	Hostname string
	procs    *list.List
}

type SystemInfo interface {
	Collect() error
}
