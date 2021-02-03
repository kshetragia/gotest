package process

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gotest/winapi"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

func (hdlr *prochdlr) getTokenInfo(class uint32, data interface{}) error {
	var size, retsize uint32

	windows.GetTokenInformation(hdlr.token, class, nil, 0, &size)

	buf := bytes.NewBuffer(make([]byte, size))
	err := windows.GetTokenInformation(hdlr.token, class, &buf.Bytes()[0], uint32(buf.Len()), &retsize)
	if err != nil {
		return errors.Wrap(err, "GetTokenInformation failed")
	}
	if size != retsize {
		err = fmt.Errorf("size mismatch (%v <=> %v)", size, retsize)
		return errors.Wrap(err, "could not read token information")
	}

	binary.Read(buf, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err:", err)
		return errors.Wrap(err, "read binary data")
	}

	return nil
}

func (hdlr *prochdlr) getProcessInfo(class uint32, data interface{}) error {
	var size, retsize uint32

	winapi.GetProcessInformation(hdlr.handler, class, nil, 0, &size)

	buf := bytes.NewBuffer(make([]byte, size))
	err := winapi.GetProcessInformation(hdlr.handler, class, &buf.Bytes()[0], uint32(buf.Len()), &retsize)
	if err != nil {
		return errors.Wrap(err, "GetProcessInformation failed")
	}
	if size != retsize {
		err = fmt.Errorf("size mismatch (%v <=> %v)", size, retsize)
		return errors.Wrap(err, "could not read process information")
	}

	binary.Read(buf, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err:", err)
		return errors.Wrap(err, "read binary data")
	}

	return nil
}
