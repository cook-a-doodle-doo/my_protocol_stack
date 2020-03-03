// +build linux

package raw

import (
	"io"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/signalsciences/ipv4"
)

const (
	CLONE_DEVICE = "/dev/net/tun"
)

type TapLinux struct {
	io.ReadWriteCloser
	name string
}

type ifreq_flags struct {
	name  [syscall.IFNAMSIZ]byte
	flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func NewTapLinux(name string) (*TapLinux, error) {
	f, err := os.OpenFile(CLONE_DEVICE, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	var req ifreq_flags
	copy(req.name[:syscall.IFNAMSIZ-1], name)
	req.flags = syscall.IFF_NO_PI | syscall.IFF_TAP

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(),
		uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		f.Close()
		return nil, errno
	}
	n := strings.Trim(string(req.name[:]), "\x00")
	if errno := setDevFlags(n, syscall.IFF_UP|syscall.IFF_RUNNING); errno != nil {
		return nil, errno
	}
	if err := setIpAddr(n, "1.0.0.10"); err != nil {
		return nil, err
	}
	if err := setNetMask(n, "0.255.255.255"); err != nil {
		return nil, err
	}
	return &TapLinux{name: n, ReadWriteCloser: f}, nil
}

//func setDevFlags(name string, mask )

func setDevFlags(name string, flags uint16) error {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(s)
	var req ifreq_flags
	copy(req.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s),
		uintptr(syscall.SIOCGIFFLAGS), uintptr(unsafe.Pointer(&req))); errno != 0 {
		return errno
	}
	req.flags |= flags
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s),
		uintptr(syscall.SIOCSIFFLAGS), uintptr(unsafe.Pointer(&req))); errno != 0 {
		return errno
	}
	return nil
}

type in_addr struct {
	s_addr uint32
}

type sockaddr_in struct {
	sin_family int16
	sin_port   uint16
	sin_addr   in_addr
	sin_zero   [8]byte
}

type ifreq_addr struct {
	name [syscall.IFNAMSIZ]byte
	addr sockaddr_in
	pad  [8]byte
}

func setIpAddr(name, addr string) error {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(s)

	var soc sockaddr_in
	soc.sin_family = syscall.AF_INET
	a, err := ipv4.FromDots(addr)
	if err != nil {
		return err
	}
	soc.sin_addr.s_addr = a

	var req ifreq_addr
	req.addr = soc
	// /* IPアドレスを変更するインターフェースを指定 */
	copy(req.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s),
		syscall.SIOCSIFADDR, uintptr(unsafe.Pointer(&req))); errno != 0 {
		return errno
	}
	return nil
}

func setNetMask(name, addr string) error {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(s)

	var soc sockaddr_in
	soc.sin_family = syscall.AF_INET
	a, err := ipv4.FromDots(addr)
	if err != nil {
		return err
	}
	soc.sin_addr.s_addr = a

	var req ifreq_addr
	req.addr = soc
	// /* IPアドレスを変更するインターフェースを指定 */
	copy(req.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s),
		syscall.SIOCSIFNETMASK, uintptr(unsafe.Pointer(&req))); errno != 0 {
		return errno
	}
	return nil
}

func (t *TapLinux) Name() string {
	return t.name
}

type sockaddr struct {
	family uint16
	addr   [16]byte
}

type ifreq_sockaddr struct {
	name [syscall.IFNAMSIZ]byte
	addr sockaddr
	pad  [8]byte
}

func (t *TapLinux) Addr() ([]byte, error) {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(s)
	var req ifreq_sockaddr
	copy(req.name[:syscall.IFNAMSIZ-1], t.name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s),
		syscall.SIOCGIFHWADDR, uintptr(unsafe.Pointer(&req))); errno != 0 {
		return nil, errno
	}
	return req.addr.addr[:], nil
}
