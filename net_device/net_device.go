package net_device

//Rx: Receiver X
//Tx: Transmiter X

type NetDevice interface {
	Name() string
	Close() error
	Rx(func() (int, int, int)) error
	Tx([]byte) (int, error)
	Addr() (int, error)
}

func Open(typeName, name string) (NetDevice, error) {
	switch typeName {
	case "TAP":
		return TapDevLinuxOpen(name)
	default:
		return nil, nil
	}
}
