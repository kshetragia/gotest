module bitbucket.org/kshetragia/gotest/src/master

go 1.15

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	gotest/process v0.0.0
)

replace gotest/proc => ./internal/proc

replace gotest/process => ./internal/process
