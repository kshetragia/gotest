package process

import (
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

func Collect() (*[]Info, error) {
	var pinfo []Info
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, errors.Wrap(err, "take process list snapshot")
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err = windows.Process32First(snapshot, &entry); err != nil {
		return nil, errors.Wrap(err, "take first process entry from process list")
	}

	const flags = windows.PROCESS_QUERY_INFORMATION
	var hdlr prochdlr
	for {
		var inf Info
		err = hdlr.open(entry.ProcessID, flags)
		if err == nil {
			inf.Name = windows.UTF16ToString(entry.ExeFile[:])
			inf.PID = entry.ProcessID
			inf.PPID = entry.ParentProcessID

			var data TOKEN_STATISTICS
			err = hdlr.getTokenInfo(uint32(syscall.TokenStatistics), &data)
			if err != nil {
				hdlr.close()
				return nil, errors.Wrap(err, "get token statistics")
			}

			inf.User.AuthenticationId = data.AuthenticationId

			hdlr.close()
			inf.Show()
		}

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break
			}
			return nil, errors.Wrap(err, "take next process entry from process list")
		}
	}

	return &pinfo, nil
}

/*
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
*/
