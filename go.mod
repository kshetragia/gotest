module bitbucket.org/kshetragia/gotest/src/master

go 1.15

require (
	github.com/pkg/errors v0.9.1
	github.com/pytimer/win-netstat v0.0.0-20180710031115-efa1aff6aafc // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c
	gotest/info v0.0.0-00010101000000-000000000000
	gotest/process v0.0.0 // indirect
	gotest/procnet v0.0.0-00010101000000-000000000000 // indirect
	gotest/users v0.0.0-00010101000000-000000000000 // indirect
	gotest/winapi v0.0.0
	gotest/privilege v0.0.0
)

replace gotest/info => ./internal/info

replace gotest/procnet => ./internal/procnet

replace gotest/process => ./internal/process

replace gotest/users => ./internal/users

replace gotest/winapi => ./internal/winapi

replace gotest/privilege => ./internal/privilege