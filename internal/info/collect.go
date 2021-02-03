package info

import (
	"gotest/process"
	"gotest/users"

	"github.com/pkg/errors"
)

// Collect is gathering information about working processes into a single representation suitable
// to deliver in various formats like JSON, YAML, XML, Plaintext, etc.
func Collect() (*FullInfo, error) {
	procs := &process.List{}
	procs.Init()
	if err := procs.Collect(); err != nil {
		return nil, errors.Wrap(err, "collect process info")
	}

	var full FullInfo
	users := users.Init()
	for ok := procs.First(); ok != false; ok = procs.Next() {
		elem := procs.Read()
		if err := users.Add(elem.UserKey); err != nil {
			return nil, errors.Wrapf(err, "collect users info for User: %v", elem.Name)
		}

		user, err := users.Get(elem.UserKey)
		if err != nil {
			return nil, errors.Wrapf(err, "get user data by luid: (%v)", elem.Name)
		}

		info := Info{
			Name:      elem.Name,
			Path:      elem.Path,
			PID:       elem.PID,
			PPID:      elem.PPID,
			StartTime: elem.StartTime,

			User:       user,
			CPUTime:    elem.CPUTime,
			MemoryInfo: elem.MemoryInfo,
		}

		full = append(full, &info)
	}

	return &full, nil
}
