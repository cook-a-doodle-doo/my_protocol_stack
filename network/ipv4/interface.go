package ipv4

import "github.com/cook-a-doodle-do/my_protocol_stack/enums"
import "github.com/cook-a-doodle-do/my_protocol_stack/network"

type IPAddr []byte

const (
	IPAddrSize uint = 4
)

func (i IPAddr) Entity() []byte {
	b := make([]byte, IPAddrSize)
	copy(b, i[:])
	return b
}

func (i IPAddr) Length() uint {
	return IPAddrSize
}

type Interface struct {
	IPAddr  IPAddr
	NetMask IPAddr
}

func NewInterface() *Interface {
	var i Interface
	i.IPAddr = make([]byte, IPAddrSize)
	i.NetMask = make([]byte, IPAddrSize)
	return &i
}

func (i *Interface) SetIPAddr(ip IPAddr) {
	copy(i.IPAddr, ip)
}

func (i *Interface) SetNetMask(mask IPAddr) {
	copy(i.NetMask, mask)
}

func (i *Interface) ProtocolAddr() network.ProtocolAddr {
	return i.IPAddr
}

func (i *Interface) EtherType() enums.EtherType {
	return enums.EtherTypeIPv4
}
