// +build windows

package process

import (
	"fmt"
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
	var err error

	// Open process tokens.
	var hdlr prochdlr
	err = hdlr.open(entry.ProcessID, windows.PROCESS_QUERY_INFORMATION)
	if err != nil {
		return fmt.Errorf("[%v(%v)] open process: %v", p.Name, p.PID, err)
	}

	// Collect common always accessible process information
	p.Name = windows.UTF16ToString(entry.ExeFile[:])
	p.PID = entry.ProcessID
	p.PPID = entry.ParentProcessID

	p.Path, err = processPath(entry.ProcessID)
	if err != nil {
		hdlr.close()
		return fmt.Errorf("[%v(%v)] get process execution path: %v", p.Name, p.PID, err)
	}

	// Getting LUID for logon session user
	var data tokenStatistics
	err = hdlr.getTokenInfo(uint32(syscall.TokenStatistics), &data)
	if err != nil {
		hdlr.close()
		return fmt.Errorf("[%v(%v)] get token statistics: %v", p.Name, p.PID, err)
	}
	p.UserKey = &data.AuthenticationId

	// Getting CPU usage info
	p.StartTime, p.CPUTime, err = hdlr.cpuInfo()
	if err != nil {
		hdlr.close()
		return fmt.Errorf("[%v(%v)] get CPU usage info: %v", p.Name, p.PID, err)
	}

	// Getting Memory usage
	p.MemoryUsage, err = hdlr.memInfo()
	if err != nil {
		hdlr.close()
		return fmt.Errorf("[%v(%v)] get memory usage info: %v", p.Name, p.PID, err)
	}

	// Save data and close descriptors
	hdlr.close()

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
