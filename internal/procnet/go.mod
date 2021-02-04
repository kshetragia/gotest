module gotest/net

go 1.15

require (
	github.com/pkg/errors v0.9.1
	github.com/pytimer/win-netstat v0.0.0-20180710031115-efa1aff6aafc
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c
	gotest/winapi v0.0.0
)

replace gotest/winapi => ../winapi