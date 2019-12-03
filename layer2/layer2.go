package layer2

import nd "github.com/cook-a-doodle-do/my_protocol_stack/net_device"

type Layer2 interface {
	Close() error
	Run() error
	Stop() error
	Tx() error
}

func New(netDev nd.NetDevice, pType string) Layer2 {
	switch pType {
	case "ETHERNET":
		return NewEthernet(netDev)
	default:
		return nil
	}
}
