package process

import (
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

func Collect() (*[]Info, error) {
	var pinfo []Info
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		errors.Wrap(err, "take process list snapshot")
		return nil, err
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err = windows.Process32First(snapshot, &entry); err != nil {
		errors.Wrap(err, "take first process entry from process list")
		return nil, err
	}

	const flags = windows.STANDARD_RIGHTS_READ | windows.PROCESS_QUERY_INFORMATION | windows.SYNCHRONIZE
	var hdlr prochdlr
	for {
		var data Info
		err = hdlr.open(entry.ProcessID, flags)
		if err == nil {
			data.Name = windows.UTF16ToString(entry.ExeFile[:])
			data.PID = entry.ProcessID
			data.PPID = entry.ParentProcessID

			hdlr.close()
		}

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break
			}
			errors.Wrap(err, "take next process entry from process list")
			break
		}

		data.Show()
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
