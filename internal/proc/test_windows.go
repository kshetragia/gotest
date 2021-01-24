package proc

import (
	"fmt"
	"unsafe"

	win "golang.org/x/sys/windows"
)

func setPrivilege(privilegeName string) {

	// Evaluate privileges to get debug info
	handle, err := win.GetCurrentProcess()
	if err != nil {
		fmt.Println("get process handler")
		return
	}
	defer win.CloseHandle(handle)

	// Getting the process token
	var tk win.Token
	if err = win.OpenProcessToken(handle, win.TOKEN_QUERY|win.TOKEN_ADJUST_PRIVILEGES, &tk); err != nil {
		fmt.Println("open process token")
		return
	}
	defer tk.Close()

	// Lookup LUID for privileges
	var luid win.LUID
	privName := win.StringToUTF16Ptr(privilegeName)
	if err = win.LookupPrivilegeValue(nil, privName, &luid); err != nil {
		fmt.Println("lookup privilege")
		return
	}
	fmt.Println("LUID:", luid)

	// Adjust current priviledges
	var tkp win.Tokenprivileges

	tkp.PrivilegeCount = 1
	tkp.Privileges[0].Luid = luid
	tkp.Privileges[0].Attributes = win.SE_PRIVILEGE_ENABLED
	tkpLen := uint32(unsafe.Sizeof(tkp))

	if err = win.AdjustTokenPrivileges(tk, false, &tkp, tkpLen, nil, nil); err != nil {
		fmt.Println("adjust new privileges")
		return
	}
}

func (info *SysInfo) Test() {
	fmt.Println("Start")

	setPrivilege("SeDebugPrivilege")
	setPrivilege("SeImpersonatePrivilege")

	fmt.Println("Finish")
}
