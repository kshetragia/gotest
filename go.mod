module bitbucket.org/kshetragia/gotest/src/master

go 1.15

require (
	github.com/pkg/errors v0.9.1
	github.com/pytimer/win-netstat v0.0.0-20180710031115-efa1aff6aafc // indirect
	gotest/info v0.0.0-00010101000000-000000000000
	gotest/process v0.0.0 // indirect
	gotest/procnet v0.0.0-00010101000000-000000000000 // indirect
	gotest/users v0.0.0-00010101000000-000000000000 // indirect
)

replace gotest/info => ./internal/info

replace gotest/procnet => ./internal/procnet

replace gotest/process => ./internal/process

replace gotest/users => ./internal/users

replace gotest/winapi => ./internal/winapi
