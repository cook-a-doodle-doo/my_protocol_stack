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

type DeviceType int

const (
	TAP DeviceType = iota
	SOCK
)

func New(t DeviceType) (Device, error) {
	switch t {
	case TAP:
		raw, err := NewTap("tap%d")
		if err != nil {
			return nil, err
		}
		return raw, nil
	case SOCK:
		sock, err := NewSock("wlp7s0")
		if err != nil {
			return nil, err
		}
		return sock, nil
	}
	return nil, errors.New("unknown device type")
}
