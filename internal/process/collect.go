// +build windows

package process

import (
	"gotest/winapi"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

func collectClose(hdlr *prochdlr, err error, errstr string) (*[]Info, error) {
	if hdlr != nil {
		hdlr.close()
	}
	return nil, errors.Wrap(err, errstr)
}

func Collect() (*[]Info, error) {
	var pinfo []Info
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return collectClose(nil, err, "take process list snapshot")
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
			// Collect common always accessible process information
			inf.Name = windows.UTF16ToString(entry.ExeFile[:])
			inf.PID = entry.ProcessID
			inf.PPID = entry.ParentProcessID

			inf.Path, err = getProcessPath(entry.ProcessID)
			if err != nil {
				return collectClose(&hdlr, err, "get process execution path")
			}

			// Getting LUID for logon session user
			var data tokenStatistics
			err = hdlr.getTokenInfo(uint32(syscall.TokenStatistics), &data)
			if err != nil {
				return collectClose(&hdlr, err, "get token statistics")
			}
			inf.User.AuthenticationID = data.AuthenticationId

			// Getting owner's Name, Domain and SID
			tUser, err := hdlr.token.GetTokenUser()
			if err != nil {
				return collectClose(&hdlr, err, "get token user")
			}
			SID := tUser.User.Sid

			inf.User.SID = SID.String()
			inf.User.Name, inf.User.Domain, _, err = SID.LookupAccount("")
			if err != nil {
				return collectClose(&hdlr, err, "lookup user Name and Domain name by SID")
			}

			// Getting LSA Logon info
			var sessionData *winapi.SecurityLogonSessionData
			err = winapi.LsaGetLogonSessionData(&inf.User.AuthenticationID, &sessionData)
			if err != nil {
				return collectClose(&hdlr, err, "get logon session data")
			}
			inf.User.LastSuccessLogon, _ = winapi.WinToUnixTime(sessionData.LogonTime)
			inf.User.SessionID = sessionData.Session

			// Getting CPU usage info
			inf.StartTime, inf.Running, err = hdlr.cpuInfo()
			if err != nil {
				return collectClose(&hdlr, err, "get CPU usage info")
			}

			// Save data and close descriptors
			hdlr.close()
			pinfo = append(pinfo, inf)
		}

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break
			}
			return collectClose(nil, err, "take next process entry from process list")
		}
	}

	return &pinfo, nil
}

// getProcessPath is working with process modules structures list.
// It can be reworked to get all modules information for the specified process.
func getProcessPath(pid uint32) (string, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, pid)
	if err != nil {
		return "", errors.Wrap(err, "modules list snapshot")
	}
	defer windows.CloseHandle(snap)

	entry := winapi.ModuleEntry32{}
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err = winapi.Module32First(snap, &entry); err != nil {
		return "", errors.Wrap(err, "first process module entry")
	}

	return windows.UTF16ToString(entry.ExePath[:]), nil
}
