package ipv4

import "github.com/cook-a-doodle-do/my_protocol_stack/enums"
import "github.com/cook-a-doodle-do/my_protocol_stack/network"

func (i IPAddr) Entity() []byte {
	var b IPAddr
	copy(b[:], i[:])
	return b[:]
}

func (i IPAddr) Length() uint {
	return IPAddrSize
}

type Interface struct {
	IPAddr  IPAddr
	NetMask IPAddr
	device  *network.Device
}

func NewInterface(dev *network.Device) *Interface {
	var i Interface
	i.device = dev
	return &i
}

func (i *Interface) SetIPAddr(ip IPAddr) {
	copy(i.IPAddr[:], ip[:])
}

func (i *Interface) SetNetMask(mask IPAddr) {
	copy(i.NetMask[:], mask[:])
}

func (i *Interface) ProtocolAddr() network.ProtocolAddr {
	return i.IPAddr
}

func (i *Interface) EtherType() enums.EtherType {
	return enums.EtherTypeIPv4
}

func (i *Interface) Tx(pn network.ProtocolNum, data []byte, dst network.ProtocolAddr) error {
	err := i.device.Tx(enums.EtherTypeIPv4, dst, data)
	return err
}
