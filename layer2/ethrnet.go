package layer2

import nd "github.com/cook-a-doodle-do/my_protocol_stack/net_device"

type Ethernet struct {
	NetDev nd.NetDevice
}

/*
48.bit: Ethernet address of destination
48.bit: Ethernet address of sender
16.bit: Protocol type = ether_type$ADDRESS_RESOLUTION
Ethernet packet data:
*/

type EtherHeader struct {
	destinationAddr [0x30]byte
	senderAddr      [0x30]byte
	protocolType    [0x10]byte
}

func NewEthernet(netDev nd.NetDevice) *Ethernet {
	return &Ethernet{NetDev: netDev}
}

func (e *Ethernet) Close() error {
	return nil
}

func (e *Ethernet) Run() error {
	return nil
}

func (e *Ethernet) Stop() error {
	return nil
}

//Transmitter X
func (e *Ethernet) Tx() error {
	return nil
}
