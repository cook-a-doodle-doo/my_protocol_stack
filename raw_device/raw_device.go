package raw_device

import "io"

//Rx: Receiver X
//Tx: Transmiter X

type RawDevice interface {
	io.ReadWriteCloser
	Name() string
	Addr() ([]byte, error)
}

const (
	TAP = iota
)

func NewRawDevice(t Type) {
	switch t {
	case TAP:
		NewTapLinux("tap%d")
	}
}
