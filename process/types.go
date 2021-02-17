package process

import "golang.org/x/sys/windows"

type tokenStatistics struct {
	TokenId            windows.LUID
	AuthenticationId   windows.LUID
	ExpirationTime     uint64
	TokenType          uint32
	ImpersonationLevel uint32
	DynamicCharged     uint32
	DynamicAvailable   uint32
	GroupCount         uint32
	PrivilegeCount     uint32
	ModifiedId         windows.LUID
}
