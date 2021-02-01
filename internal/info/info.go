package info

import (
	"gotest/process"
	"gotest/users"
)

// Info consists all data about working proceses
type Info struct {
	Procs *process.List
	Users *users.Users
}
