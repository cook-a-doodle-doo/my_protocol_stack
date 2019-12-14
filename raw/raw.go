package raw

import (
	"errors"
	"io"
)

//Rx: Receiver X
//Tx: Transmiter X

type Device interface {
	io.ReadWriteCloser
	Name() string
	Addr() ([]byte, error)
}

const (
	TAP = iota
)

func New(t int) (Device, error) {
	switch t {
	case TAP:
		raw, err := NewTapLinux("tap%d")
		if err != nil {
			return nil, err
		}
		return raw, nil
	}
	return nil, errors.New("unknown device type")
}
