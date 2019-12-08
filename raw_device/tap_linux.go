package raw_device

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
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

type ttIF struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func NewTapLinux(name string) (*TapLinux, error) {
	f, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	var req ttIF
	copy(req.Name[:], name)
	req.Flags = cIFF_NO_PI | cIFF_TAP

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(syscall.TUNSETIFF),
		uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		var msg string
		switch errno {
		case 1:
			msg = fmt.Sprintf("ioctl: requires root privileges")
		default:
			msg = fmt.Sprintf("ioctl: %d", errno)
		}
		return nil, errors.New(msg)
	}
	n := strings.Trim(string(req.Name[:]), "\x00")
	return &TapLinux{name: n, file: f}, nil
}

func (t *TapLinux) Close() error {
	return t.file.Close()
}

func (t *TapLinux) Rx() {
}

func (t *TapLinux) Tx() {
}

func (t *TapLinux) Name() {
}

func (t *TapLinux) Addr() {
}
