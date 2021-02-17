// +build windows

package procnet

import (
	"testing"
)

func TestCollect(t *testing.T) {
	netstat := Init()
	if err := netstat.Collect(); err != nil {
		t.Fatalf("Could not collect net IO stat")
	}
}
