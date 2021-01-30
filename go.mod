module bitbucket.org/kshetragia/gotest/src/master

go 1.15

require (
	github.com/pkg/errors v0.9.1
	gotest/process v0.0.0
	gotest/winapi v0.0.0-00010101000000-000000000000 // indirect
)

replace gotest/process => ./internal/process

replace gotest/winapi => ./internal/winapi
