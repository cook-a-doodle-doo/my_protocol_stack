package socket

/*
import (
	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
	"github.com/cook-a-doodle-do/my_protocol_stack/link"
	"github.com/cook-a-doodle-do/my_protocol_stack/link/arp"
	"github.com/cook-a-doodle-do/my_protocol_stack/link/ethernet"
	"github.com/cook-a-doodle-do/my_protocol_stack/network"
	"github.com/cook-a-doodle-do/my_protocol_stack/network/icmp"
	"github.com/cook-a-doodle-do/my_protocol_stack/network/ipv4"
)

type ProtocolFamily uint

const (
	INET ProtocolFamily = iota
	INET6
	PACKET
)

type Type uint8

const (
	Stream Type = iota
	Dgram
)

type SocketType struct {
	PF       ProtocolFamily
	Type     Type
	Protocol Protocol
}

const (
	INetStreamIPProtoTCP = SocketType{PF: INET, Type: Stream, Protocol: IPProtoTCP}
	INetDgramIPProtoUDP  = SocketType{PF: INET, Type: Dgram, Protocol: IPProtoUDP}
)

func New(st SocketType) (Socket, error) {}

type Socket struct {
}

func (s *Socket) Bind(pf ProtocolFamily, sa SockAddr) error {
}

func (s *Socket) Connect() {
}

func (s *Socket) Accept() Connection {
}

func (s *Socket) Listen(waitnum int) Connection {
}

type Connection interface {
	io.ReadWriteCloser
}

func init() {
	linkdev, err := ethernet.NewDevice()
	if err != nil {
		panic(err)
	}
	link.AppendDevice(linkdev)
	link.RegistProtocol(enums.EtherTypeARP, arp.CallbackHandler)
	network.RegistProtocol(enums.EtherTypeIPv4, ipv4.CallbackHandler)
	netdev, err := network.NewDevice(linkdev)
	if err != nil {
		panic(err)
	}
	ipv4IF := ipv4.NewInterface()
	ipv4IF.SetIPAddr([4]byte{10, 0, 0, 2})
	ipv4IF.SetNetMask([4]byte{255, 255, 255, 255})
	netdev.AppendInterface(ipv4IF)
	ipv4.RegistProtocol(ipv4.ProtocolTypeICMP, icmp.CallbackHandler)
}
*/

//type AddressFamily uint
//
//const (
//	AFIPv4 AddressFamily = iota
//)
//
//type ProtocolFamily uint
