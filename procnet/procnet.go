// +build windows

package procnet

// I don't like use doubtful repositories but i have not time to realize it by myself
// https://github.com/sherifeldeeb/win-netstat
import (
	"fmt"
	"gotest/winapi"
	"net"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
	netwin "github.com/pytimer/win-netstat"
	"golang.org/x/sys/windows"
)

// NetIOStat consist net IO statistics
type NetIOStat struct {
	BytesIn  uint64
	BytesOut uint64
	BandIn   uint64 // Bandwidth (bytes in sec)
	BandOut  uint64 // Bandwidth (bytes in sec)
}

// NetInfo collects net IO statistics in bytes by process
type NetInfo struct {
	PID        uint32
	Proto      string
	LocalAddr  string
	LocalPort  uint16
	RemoteAddr string
	RemotePort uint16
	State      string

	IOStat NetIOStat `json:"IOStat"`
}

// NetStat aggregates net IO statistics for the process list
type NetStat map[uint32][]*NetInfo

// Init creates new NetStat pool
func Init() NetStat {
	return NetStat(make(map[uint32][]*NetInfo))
}

// Collect gathers information about network IO stat
func (ns *NetStat) Collect() error {
	if err := ns.NetStatTCPProto(syscall.AF_INET); err != nil {
		return errors.Wrap(err, "get IOStat for IPv4")
	}
	if err := ns.NetStatTCPProto(syscall.AF_INET6); err != nil {
		return errors.Wrap(err, "get IOStat for IPv6")
	}
	return nil
}

func (ns *NetStat) NetStatTCPProto(proto uint32) error {
	if proto != syscall.AF_INET && proto != syscall.AF_INET6 {
		return errors.New("unsupported Net protocol")
	}

	var tSize uint32
	var buf []byte
	var err error

	var ptable unsafe.Pointer

	// Execute multiple times to get and enlarge receiving buffer size
	for {
		if len(buf) > 0 {
			ptable = (unsafe.Pointer(&buf[0]))
		}
		// See also: https://docs.microsoft.com/en-us/windows/win32/api/iphlpapi/nf-iphlpapi-getextendedtcptable
		err = netwin.GetExtendedTcpTable(uintptr(ptable), &tSize, false, proto, netwin.TCP_TABLE_OWNER_PID_ALL, 0)
		if err != nil {
			if err != windows.ERROR_INSUFFICIENT_BUFFER {
				return errors.Wrap(err, "get TCP IPv4 connections table")
			}
			buf = make([]byte, tSize)
			continue
		}
		break
	}

	// Parse protocol specific data types
	if proto == syscall.AF_INET {
		ns.getTCP4Stat(&buf)
	} else {
		ns.getTCP6Stat(&buf)
	}
	return nil
}

func (ns *NetStat) getTCP4Stat(pbuf *[]byte) {

	nss := *ns
	buf := *pbuf
	ptable := (*netwin.MIB_TCPTABLE_OWNER_PID)(unsafe.Pointer(&buf[0]))

	index := int(unsafe.Sizeof(ptable.DwNumEntries))
	step := int(unsafe.Sizeof(ptable.Table))
	for i := 0; i < int(ptable.DwNumEntries); i++ {
		mibs := (*netwin.MIB_TCPROW_OWNER_PID)(unsafe.Pointer(&buf[index]))
		pid := mibs.DwOwningPid

		ni := NetInfo{
			PID:        mibs.DwOwningPid,
			Proto:      "IPv4",
			LocalAddr:  parseIPv4(mibs.DwLocalAddr),
			LocalPort:  decodePort(mibs.DwLocalPort),
			RemoteAddr: parseIPv4(mibs.DwRemoteAddr),
			RemotePort: decodePort(mibs.DwLocalPort),
			State:      netwin.TCPStatuses[netwin.MIB_TCP_STATE(mibs.DwState)],
		}

		if mibs.DwState == uint32(winapi.MIB_TCP_STATE_ESTAB) {
			row := winapi.MIB_TCPROW{
				DwLocalAddr:  mibs.DwLocalAddr,
				DwLocalPort:  mibs.DwLocalPort,
				DwRemoteAddr: mibs.DwRemoteAddr,
				DwRemotePort: mibs.DwRemotePort,
				State:        mibs.DwState,
			}

			// Get TCP counters
			// Dynamic data has uncontrolled buffer size. We will take it enough to get all required data.
			// https://docs.microsoft.com/ru-ru/windows/win32/api/iphlpapi/nf-iphlpapi-getpertcpconnectionestats?redirectedfrom=MSDN
			query := winapi.TcpConnectionEstatsData
			rodSize := uint64(unsafe.Sizeof(winapi.TCP_ESTATS_DATA_ROD_v0{}) * 10)
			bufData := make([]byte, rodSize)
			winapi.GetPerTcpConnectionEStats(&row, query, nil, 0, 0, nil, 0, 0, &bufData[0], 0, rodSize)
			rodData := (*winapi.TCP_ESTATS_DATA_ROD_v0)(unsafe.Pointer(&bufData))

			// Get bandwidth data
			// Dynamic data has uncontrolled buffer size. We will take it enought to get all required data.
			query = winapi.TcpConnectionEstatsBandwidth
			rodSize = uint64(unsafe.Sizeof(winapi.TCP_ESTATS_BANDWIDTH_ROD_v0{}) * 10)
			bufBand := make([]byte, rodSize)
			winapi.GetPerTcpConnectionEStats(&row, query, nil, 0, 0, nil, 0, 0, &bufBand[0], 0, rodSize)
			rodBand := (*winapi.TCP_ESTATS_BANDWIDTH_ROD_v0)(unsafe.Pointer(&bufBand))

			ni.IOStat = NetIOStat{
				BytesIn:  rodData.DataBytesIn,
				BytesOut: rodData.DataBytesOut,
				BandIn:   rodBand.InboundBandwidth,
				BandOut:  rodBand.OutboundBandwidth,
			}
		}

		nss[pid] = append(nss[pid], &ni)
		index += step
	}
}

func (ns *NetStat) getTCP6Stat(pbuf *[]byte) {
	nss := *ns
	buf := *pbuf
	ptable := (*netwin.MIB_TCPTABLE_OWNER_PID)(unsafe.Pointer(&buf[0]))

	index := int(unsafe.Sizeof(ptable.DwNumEntries))
	step := int(unsafe.Sizeof(ptable.Table))
	for i := 0; i < int(ptable.DwNumEntries); i++ {
		mibs := (*netwin.MIB_TCP6ROW_OWNER_PID)(unsafe.Pointer(&buf[index]))
		pid := mibs.DwOwningPid

		ni := NetInfo{
			PID:        mibs.DwOwningPid,
			Proto:      "IPv6",
			LocalAddr:  parseIPv6(mibs.UcLocalAddr),
			LocalPort:  decodePort(mibs.DwLocalPort),
			RemoteAddr: parseIPv6(mibs.UcRemoteAddr),
			RemotePort: decodePort(mibs.DwLocalPort),
			State:      netwin.TCPStatuses[netwin.MIB_TCP_STATE(mibs.DwState)],
		}

		if mibs.DwState == uint32(winapi.MIB_TCP_STATE_ESTAB) {
			row := winapi.MIB_TCP6ROW{
				LocalAddr:    mibs.UcLocalAddr,
				DwLocalPort:  mibs.DwLocalPort,
				RemoteAddr:   mibs.UcRemoteAddr,
				DwRemotePort: mibs.DwRemotePort,
				State:        mibs.DwState,
			}

			// Get TCP counters
			// Dynamic data has uncontrolled buffer size. We will take it enough to get all required data.
			// https://docs.microsoft.com/ru-ru/windows/win32/api/iphlpapi/nf-iphlpapi-getpertcpconnectionestats?redirectedfrom=MSDN
			query := winapi.TcpConnectionEstatsData
			rodSize := uint64(unsafe.Sizeof(winapi.TCP_ESTATS_DATA_ROD_v0{}) * 100)
			bufData := make([]byte, rodSize)
			winapi.GetPerTcp6ConnectionEStats(&row, query, nil, 0, 0, nil, 0, 0, &bufData[0], 0, rodSize)
			rodData := (*winapi.TCP_ESTATS_DATA_ROD_v0)(unsafe.Pointer(&bufData))

			// Get bandwidth data
			// Dynamic data has uncontrolled buffer size. We will take it enought to get all required data.
			query = winapi.TcpConnectionEstatsBandwidth
			rodSize = uint64(unsafe.Sizeof(winapi.TCP_ESTATS_BANDWIDTH_ROD_v0{}) * 100)
			bufBand := make([]byte, rodSize)
			winapi.GetPerTcp6ConnectionEStats(&row, query, nil, 0, 0, nil, 0, 0, &bufBand[0], 0, rodSize)
			rodBand := (*winapi.TCP_ESTATS_BANDWIDTH_ROD_v0)(unsafe.Pointer(&bufBand))

			ni.IOStat = NetIOStat{
				BytesIn:  rodData.DataBytesIn,
				BytesOut: rodData.DataBytesOut,
				BandIn:   rodBand.InboundBandwidth,
				BandOut:  rodBand.OutboundBandwidth,
			}
		}

		nss[pid] = append(nss[pid], &ni)
		index += step
	}
}

func decodePort(port uint32) uint16 {
	return syscall.Ntohs(uint16(port))
}

func parseIPv4(addr uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", addr&255, addr>>8&255, addr>>16&255, addr>>24&255)
}

func parseIPv6(addr [16]byte) string {
	var ret [16]byte
	for i := 0; i < 16; i++ {
		ret[i] = uint8(addr[i])
	}

	// convert []byte to net.IP
	ip := net.IP(ret[:])
	return ip.String()
}
