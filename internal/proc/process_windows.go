package proc

import (
	"github.com/pkg/errors"
	win "golang.org/x/sys/windows"
)

type ProcInfo struct {
	Name      string
	Path      string
	PID       uint32
	PPID      uint32
	Owner     string
	Domain    string
	SessionId uint32
	LUID      win.LUID
	SID       string
	CPU       CPUInfo
}

func (p *ProcInfo) collectProcessInfo(hostname string, entry *win.ProcessEntry32) error {
	var err error

	p.Name = win.UTF16ToString(entry.ExeFile[:])
	p.PID = entry.ProcessID
	p.PPID = entry.ParentProcessID

	// Get session ID
	if err = win.ProcessIdToSessionId(entry.ProcessID, &(p.SessionId)); err != nil {
		errors.Wrap(err, "get session ID")
		return nil
	}

	// Get process handler
	const flags = win.STANDARD_RIGHTS_READ | win.PROCESS_QUERY_INFORMATION | win.SYNCHRONIZE
	handler, err := win.OpenProcess(flags, false, entry.ProcessID)
	if err != nil {
		// Seems we do not have permissions to open handler
		return nil
	}
	defer win.CloseHandle(handler)

	var token win.Token
	if err = win.OpenProcessToken(handler, win.TOKEN_QUERY, &token); err != nil {
		errors.Wrap(err, "open token")
		return err
	}
	defer token.Close()

	tokenUser, err := token.GetTokenUser()
	if err != nil {
		errors.Wrap(err, "get token user")
		return err
	}

	sid := tokenUser.User.Sid
	p.SID = sid.String()

	p.Owner, p.Domain, _, err = sid.LookupAccount(hostname)
	if err != nil {
		errors.Wrap(err, "lookup SID process owner")
		return err
	}

	// Get CPU info
	p.CPU.Read(handler)

	return err
}
