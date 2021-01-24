package proc

import (
	"container/list"
	"fmt"
	"sort"
	"unsafe"

	"github.com/pkg/errors"
	win "golang.org/x/sys/windows"
)

func (info *SysInfo) Collect() error {
	var err error

	info.procs = list.New()

	info.Hostname, err = win.ComputerName()
	if err != nil {
		errors.Wrap(err, "get host name")
		return err
	}

	if err = info.takeProcessList(); err != nil {
		errors.Wrap(err, "take process list")
		return err
	}

	return nil
}

func (info *SysInfo) takeProcessList() error {
	snapshot, err := win.CreateToolhelp32Snapshot(win.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		errors.Wrap(err, "take process list snapshot")
		return err
	}
	defer win.CloseHandle(snapshot)

	var entry win.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err = win.Process32First(snapshot, &entry); err != nil {
		errors.Wrap(err, "take first process entry from process list")
		return err
	}

	for {
		var pinfo ProcInfo

		if err = pinfo.collectProcessInfo(info.Hostname, &entry); err != nil {
			errors.Wrap(err, "collect process info")
			break
		}

		info.procs.PushBack(pinfo)

		if err := win.Process32Next(snapshot, &entry); err != nil {
			if err == win.ERROR_NO_MORE_FILES {
				break
			}
			errors.Wrap(err, "take next process entry from process list")
			break
		}
	}

	return nil
}

func (info *SysInfo) Free() {
	/*
		for elem := info.procs.Front(); elem.Next() != nil; elem = elem.Next() {
			e := elem.Value.(ProcInfo)
		}
	*/
}

func (info *SysInfo) Show() {

	data := make(map[uint32]ProcInfo)

	for elem := info.procs.Front(); elem.Next() != nil; elem = elem.Next() {
		e := elem.Value.(ProcInfo)
		data[e.PID] = e
	}

	keys := make([]uint32, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, pid := range keys {
		e := data[pid]
		user := ""
		if e.Owner != "" {
			user = "\\\\" + e.Domain + "\\" + e.Owner
		}
		fmt.Printf("[%v] %v \n", pid, e.Name)
		fmt.Printf("\tSession: %v\n", e.SessionId)
		fmt.Printf("\tLUID: %v\n", e.LUID)
		fmt.Printf("\tSID: %v\n", e.SID)
		fmt.Printf("\tUser: %v\n", user)
		fmt.Printf("\tCreation Time: %v\n", e.CPU.creation.String())
		fmt.Printf("\tRunning Time: %v\n", e.CPU.running.String())
	}
}
