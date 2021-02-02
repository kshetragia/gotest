// +build windows

package process

import (
	"gotest/winapi"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

// Collect gathers all executed processes info and returns process list
// It intentionaly cleanup process list and create new one.
func (plist *List) Collect() error {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return errors.Wrap(err, "take process list snapshot")
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err = windows.Process32First(snapshot, &entry); err != nil {
		return errors.Wrap(err, "take first process entry from process list")
	}

	for {
		var inf Process
		if err := inf.collect(&entry); err != nil {
			// Temporary ignore errors
			// Just show it on display
			// fmt.Printf("collect process info: %v", err)
		} else {
			plist.Add(&inf)
		}

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break
			}
			return errors.Wrap(err, "take next process entry from process list")
		}
	}

	return nil
}

func (p *Process) collect(entry *windows.ProcessEntry32) error {
	// Open process tokens.
	var hdlr prochdlr
	err := hdlr.open(entry.ProcessID, windows.PROCESS_QUERY_INFORMATION)
	if err != nil {
		return errors.Wrapf(err, "[%v(%v)] open process", p.Name, p.PID)
	}
	defer hdlr.close()

	// Collect common always accessible process information
	p.Name = windows.UTF16ToString(entry.ExeFile[:])
	p.PID = entry.ProcessID
	p.PPID = entry.ParentProcessID

	p.Path, err = processPath(entry.ProcessID)
	if err != nil {
		return errors.Wrapf(err, "[%v(%v)] get process execution path", p.Name, p.PID)
	}

	// Getting LUID for logon session user
	var data tokenStatistics
	err = hdlr.getTokenInfo(uint32(syscall.TokenStatistics), &data)
	if err != nil {
		return errors.Wrapf(err, "[%v(%v)] get token statistics", p.Name, p.PID)
	}
	p.UserKey = &data.AuthenticationId

	// Getting CPU usage info
	p.StartTime, p.CPUTime, err = hdlr.cpuInfo()
	if err != nil {
		return errors.Wrapf(err, "[%v(%v)] get CPU usage info", p.Name, p.PID)
	}

	// Getting Memory usage
	p.MemoryInfo, err = hdlr.memInfo()
	if err != nil {
		return errors.Wrapf(err, "[%v(%v)] get memory usage info", p.Name, p.PID)
	}

	return nil
}

// processPath is working with process modules structures list.
// It can be reworked to get all modules information for the specified process.
func processPath(pid uint32) (string, error) {
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
