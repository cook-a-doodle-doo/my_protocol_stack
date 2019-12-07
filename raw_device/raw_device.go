package raw_device

//Rx: Receiver X
//Tx: Transmiter X

type RawDevice interface {
	Close() error
	Rx(func() (int, int, int)) error
	Tx([]byte) (int, error)
	Name() string
	Addr() (int, error)
}
