// +build windows

package winapi

const (
	// do not reorder
	ProcessMemoryPriority = 1 + iota
	ProcessMemoryExhaustionInfo
	ProcessAppMemoryInfo
	ProcessInPrivateInfo
	ProcessPowerThrottling
	ProcessReservedValue1
	ProcessTelemetryCoverageInfo
	ProcessProtectionLevelInfo
	ProcessLeapSecondInfo
	ProcessInformationClassMax
)

type AppMemoryInformation struct {
	AvailableCommit        uint64
	PrivateCommitUsage     uint64
	PeakPrivateCommitUsage uint64
	TotalCommitUsage       uint64
}
