package main

import (
	"fmt"
	"gotest/info"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func showErrors(err error) {
	for e := err; e != nil; e = errors.Unwrap(err) {
		fmt.Println("Error:", e)
	}
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s: %d\n", f, f)
		}
	}
	return
}

func Show(info *info.FullInfo) {
	for _, e := range *info {
		u := e.User
		user := "\\\\" + u.Domain + "\\" + u.Name

		fmt.Printf("[%v] %v \n", e.PID, e.Name)
		fmt.Printf("\tPath: %v\n", e.Path)
		fmt.Printf("\tExecution time: %v\n", e.StartTime)

		fmt.Printf("\tOwner:\n")
		fmt.Printf("\t  Name: %v\n", user)
		fmt.Printf("\t  SessionId: %v\n", u.SessionID)
		fmt.Printf("\t  Last Login: %v\n", u.LastSuccessLogon)
		fmt.Printf("\t  SID: %v\n", u.SID)

		fmt.Printf("\tCPU:\n")
		fmt.Printf("\t  Kernel: %.2fs\n", e.CPUTime.Kernel)
		fmt.Printf("\t  User: %.2fs\n", e.CPUTime.User)
		fmt.Printf("\t  System: %.2f%%\n", e.CPUTime.System)
		fmt.Printf("\t  Total:  %v\n", e.CPUTime.Total)

		fmt.Printf("\tMemory:\n")
		fmt.Printf("\t  WorkingSetSize: %v\n", e.MemoryInfo.WorkingSetSize)
		fmt.Printf("\t  QuotaPagedPoolUsage: %v\n", e.MemoryInfo.QuotaPagedPoolUsage)
		fmt.Printf("\t  QuotaNonPagedPoolUsage: %v\n", e.MemoryInfo.QuotaNonPagedPoolUsage)
		fmt.Printf("\t  PrivateUsage: %v\n", e.MemoryInfo.PrivateUsage)

		// fmt.Printf("\tLUID: %v\n", e.User.AuthenticationID)
		fmt.Println()
	}

}

func main() {
	pinfo, err := info.Collect()
	if err != nil {
		showErrors(err)
	}
	Show(pinfo)

	json, err := pinfo.Json()
	if err != nil {
		showErrors(err)
	}
	fmt.Println(string(json))
}
