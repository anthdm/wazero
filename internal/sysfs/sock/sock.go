package sock

import "fmt"

// ProtocolFamily is a socket protocol family.
type ProtocolFamily int32

const (
	UnspecifiedFamily ProtocolFamily = iota
	InetFamily
	Inet6Family
	UnixFamily
)

func (pf ProtocolFamily) String() string {
	switch pf {
	case UnspecifiedFamily:
		return "UnspecifiedFamily"
	case InetFamily:
		return "InetFamily"
	case Inet6Family:
		return "Inet6Family"
	case UnixFamily:
		return "UnixFamily"
	default:
		return fmt.Sprintf("ProtocolFamily(%d)", pf)
	}
}

// Protocol is a socket protocol.
type Protocol int32

const (
	IPProtocol Protocol = iota
	TCPProtocol
	UDPProtocol
)

func (p Protocol) String() string {
	switch p {
	case IPProtocol:
		return "IPProtocol"
	case TCPProtocol:
		return "TCPProtocol"
	case UDPProtocol:
		return "UDPProtocol"
	default:
		return fmt.Sprintf("Protocol(%d)", p)
	}
}

// SocketType is a type of socket.
type SocketType int32

const (
	AnySocket SocketType = iota
	DatagramSocket
	StreamSocket
)

func (st SocketType) String() string {
	switch st {
	case AnySocket:
		return "AnySocket"
	case DatagramSocket:
		return "DatagramSocket"
	case StreamSocket:
		return "StreamSocket"
	default:
		return fmt.Sprintf("SocketType(%d)", st)
	}
}
