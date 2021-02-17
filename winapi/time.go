// +build windows

package winapi

import (
	"time"
)

// WinToUnixTime converts windows time to Unix time format.
// Microsoft's standard time format measured as the number of
// 100-nano-second intervals since 1st January 1601, 00:00:00 UTC.
// But Unix time starts from 1st January 1970 UTC
func WinToUnixTime(nsec uint64) time.Time {
	const winEpoch uint64 = 116444736000000000 // (369 * 365 + 89) * 24 * 3600 * 10000000
	return time.Unix(0, int64((nsec-winEpoch)*100))
}
