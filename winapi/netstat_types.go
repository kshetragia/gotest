// +build windows

package winapi

// TCP_ESTATS_TYPE enumeration defines the type of extended statistics for a TCP connection that is requested or being set.
// https://docs.microsoft.com/en-us/windows/win32/api/tcpestats/ne-tcpestats-tcp_estats_type
type TCP_ESTATS_TYPE int32

const (
	// do not reorder
	TcpConnectionEstatsSynOpts TCP_ESTATS_TYPE = 1 + iota
	TcpConnectionEstatsData
	TcpConnectionEstatsSndCong
	TcpConnectionEstatsPath
	TcpConnectionEstatsSendBuff
	TcpConnectionEstatsRec
	TcpConnectionEstatsObsRec
	TcpConnectionEstatsBandwidth
	TcpConnectionEstatsFineRtt
	TcpConnectionEstatsMaximum
)

// MIB_TCP_STATE The state of the TCP connection.
type MIB_TCP_STATE uint32

const (
	// do not reorder
	MIB_TCP_STATE_CLOSED MIB_TCP_STATE = 1 + iota
	MIB_TCP_STATE_LISTEN
	MIB_TCP_STATE_SYN_SENT
	MIB_TCP_STATE_SYN_RCVD
	MIB_TCP_STATE_ESTAB
	MIB_TCP_STATE_FIN_WAIT1
	MIB_TCP_STATE_FIN_WAIT2
	MIB_TCP_STATE_CLOSE_WAIT
	MIB_TCP_STATE_CLOSING
	MIB_TCP_STATE_LAST_ACK
	MIB_TCP_STATE_TIME_WAIT
	MIB_TCP_STATE_DELETE_TCB
)

// MIB_TCPROW_LH structure contains information that descibes an IPv4 TCP connection.
// https://docs.microsoft.com/en-us/windows/win32/api/tcpmib/ns-tcpmib-mib_tcprow_lh
type MIB_TCPROW struct {
	State        uint32
	DwLocalAddr  uint32
	DwLocalPort  uint32
	DwRemoteAddr uint32
	DwRemotePort uint32
}

// MIB_TCP6ROW structure contains information that describes an IPv6 TCP connection.
// https://docs.microsoft.com/en-us/windows/win32/api/tcpmib/ns-tcpmib-mib_tcp6row
type MIB_TCP6ROW struct {
	State           uint32
	LocalAddr       [16]byte
	DwLocalScopeId  uint32
	DwLocalPort     uint32
	RemoteAddr      [16]byte
	DwRemoteScopeId uint32
	DwRemotePort    uint32
}

type TCP_ESTATS_BANDWIDTH_ROD_v0 struct {
	OutboundBandwidth       uint64
	InboundBandwidth        uint64
	OutboundInstability     uint64
	InboundInstability      uint64
	OutboundBandwidthPeaked byte // BOOLEAN
	InboundBandwidthPeaked  byte // BOOLEAN
}

type TCP_ESTATS_DATA_ROD_v0 struct {
	DataBytesOut      uint64
	DataSegsOut       uint64
	DataBytesIn       uint64
	DataSegsIn        uint64
	SegsOut           uint64
	SegsIn            uint64
	SoftErrors        uint32
	SoftErrorReason   uint32
	SndUna            uint32
	SndNxt            uint32
	SndMax            uint32
	ThruBytesAcked    uint64
	RcvNxt            uint32
	ThruBytesReceived uint64
}
