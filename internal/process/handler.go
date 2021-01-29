package process

// Often we need to open process handler and process token to get any info.
// So.. we can manage it simultaneously for our

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

type prochdlr struct {
	count   uint8
	handler windows.Handle
	token   windows.Token
}

// handlerOpen is opening process' handler and token descriptors
// If descriptors already opened it does nothing, just increase some internal counter of function calls
// access - desired access flags for handler
// pid - process ID
func (h *prochdlr) open(pid uint32, access uint32) (err error) {
	h.handler, err = windows.OpenProcess(access, false, pid)
	if err != nil {
		errors.Wrap(err, "open process handler")
		return err
	}

	if err = windows.OpenProcessToken(h.handler, windows.TOKEN_QUERY, &h.token); err != nil {
		errors.Wrap(err, "open process token")
		return err
	}

	return nil
}

// close() close opened process descriptors
func (h *prochdlr) close() {
	// Close descriptors in reverse order
	if h.token != 0 {
		h.token.Close()
		h.token = 0
	}

	if h.handler != 0 {
		windows.CloseHandle(h.handler)
		h.handler = 0
	}
}
