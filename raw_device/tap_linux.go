package raw_device

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	CLONE_DEVICE = "/dev/net/tun"
)

type TapLinux struct {
	io.ReadWriteCloser
	name string
}

type ttIF struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func NewTapLinux(name string) (*TapLinux, error) {
	f, err := os.OpenFile(CLONE_DEVICE, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	var req ttIF
	copy(req.Name[:], name)
	req.Flags = syscall.IFF_NO_PI | syscall.IFF_TAP

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(syscall.TUNSETIFF),
		uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		var msg string
		//TODO 適切なエラーメッセージ
		switch errno {
		case 1:
			msg = fmt.Sprintf("ioctl: requires root privileges")
		default:
			msg = fmt.Sprintf("ioctl: %d", errno)
		}
		f.Close()
		return nil, errors.New(msg)
	}
	n := strings.Trim(string(req.Name[:]), "\x00")
	return &TapLinux{name: n, ReadWriteCloser: f}, nil
}

func (t *TapLinux) Name() string {
	return t.name
}

func (t *TapLinux) Addr() {
}
