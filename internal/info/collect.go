package info

import (
	"gotest/process"
	"gotest/users"

	"github.com/pkg/errors"
)

func Collect() (*Info, error) {
	var info Info

	procs := &process.List{}
	procs.Init()
	if err := procs.Collect(); err != nil {
		return nil, errors.Wrap(err, "collect process info")
	}

	users := users.New()
	for ok := procs.First(); ok != false; ok = procs.Next() {
		elem := procs.Read()
		if err := users.Add(elem.UserKey); err != nil {
			return nil, errors.Wrapf(err, "collect users info for User: %v", elem.Name)
		}
	}
	info.Procs = procs
	info.Users = &users

	return &info, nil
}
