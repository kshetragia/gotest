module bitbucket.org/kshetragia/gotest/src/master

go 1.15

require (
	github.com/pkg/errors v0.9.1
	gotest/info v0.0.0-00010101000000-000000000000
	gotest/process v0.0.0 // indirect
	gotest/users v0.0.0-00010101000000-000000000000 // indirect
)

replace gotest/process => ./internal/process

replace gotest/winapi => ./internal/winapi

replace gotest/users => ./internal/users

replace gotest/info => ./internal/info
