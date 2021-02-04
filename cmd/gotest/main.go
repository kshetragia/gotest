package main

import (
	"context"
	"fmt"
	"gotest/info"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	var count int
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

		fmt.Printf("\tNetwork:\n")
		for key, val := range e.NetInfo {
			fmt.Printf("  [%v]: %v\n", key, val.Proto)
			fmt.Printf("    State: %v\n", val.State)
			fmt.Printf("    LocalAddr: %v:%v\n", val.LocalAddr, val.LocalPort)
			fmt.Printf("    RemoteAddr: %v:%v\n", val.RemoteAddr, val.RemotePort)
			fmt.Printf("    In: %v Bytes (%v Bytes/s)\n", val.IOStat.BytesIn, val.IOStat.BandIn)
			fmt.Printf("    Out: %v Bytes (%v Bytes/s)\n", val.IOStat.BytesOut, val.IOStat.BandOut)
		}

		fmt.Println()
		count++
	}
	fmt.Println("Count:", count)
}

func infoHandler(w http.ResponseWriter, req *http.Request) {
	pinfo, err := info.Collect()
	if err != nil {
		showErrors(err)
		return
	}
	json, err := pinfo.Json()
	if err != nil {
		showErrors(err)
		return
	}
	fmt.Fprintf(w, string(json))
	// Show(pinfo)
}

func main() {
	fmt.Println("Trying to raise up HTTP server...")

	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", infoHandler)

	// Raise up HTTP server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("Server error:", err)
		}
	}()

	fmt.Println("Server is running on *:8080")

	// Wait SIGIING
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("Shutting down HTTP server")

	// Try to shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	}
}
