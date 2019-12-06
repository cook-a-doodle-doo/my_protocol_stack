package raw_device

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

const (
	cIFF_TUN   = 0x0001
	cIFF_TAP   = 0x0002
	cIFF_NO_PI = 0x1000
)

type TapDevLinux struct {
	file *os.File
	name string
}

type tuntapInterface struct {
	name    [0x10]byte
	flags   uint16
	padding [0x28 - 0x10 - 2]byte
}

func TapDevLinuxOpen(name string) (*TapDevLinux, error) {
	f, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	req := tuntapInterface{flags: cIFF_TAP | cIFF_NO_PI}
	copy(req.name[:], name)
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(syscall.TUNSETIFF),
		uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		f.Close()
		return nil, errno
	}
	return &TapDevLinux{name: name, file: f}, nil
}

func (t *TapDevLinux) Name() string {
	return string(t.name)
}

func (t *TapDevLinux) Close() error {
	return t.file.Close()
}

//Receiver X
func (t *TapDevLinux) Rx(c func() (int, int, int)) error {
	return nil
}

//Transmitter X
func (t *TapDevLinux) Tx(p []byte) (int, error) {
	return t.file.Write(p)
}

//Address
func (t *TapDevLinux) Addr() (int, error) {
	return 1, nil
}
