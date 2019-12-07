package raw_device

import (
	"os"
)

const (
	cIFF_TUN   = 0x0001
	cIFF_TAP   = 0x0002
	cIFF_NO_PI = 0x1000
)

type TapLinux struct {
	name string
	file *os.File
}

func NewTapLinux() (*TapLinux, error) {
	fd, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
}

func (t *TapLinux) Close() {
}

func (t *TapLinux) Rx() {
}

func (t *TapLinux) Tx() {
}

func (t *TapLinux) Name() {
}

func (t *TapLinux) Addr() {
}
