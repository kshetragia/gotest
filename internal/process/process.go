package process

type Process struct {
	Name   string
	Path   string
	User   string
	Domain string
	PID    uint32
	PPID   uint32
	LUID   windows.LUID
	Rusage
}
