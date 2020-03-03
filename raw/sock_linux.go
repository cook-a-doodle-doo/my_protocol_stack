package raw

import (
	"fmt"
	"syscall"
	"unsafe"
)

type Sock struct {
	fd   int
	name string
}

func NewSock(name string) (*Sock, error) {
	s, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		return nil, err
	}
	index, err := getIFIndex(name)
	if err != nil {
		syscall.Close(s)
		return nil, fmt.Errorf("Can't get IF Index\n")
	}
	addr := &syscall.SockaddrLinklayer{
		Protocol: syscall.ETH_P_ALL,
		Ifindex:  int(index),
	}

	if err = syscall.Bind(s, addr); err != nil {
		syscall.Close(s)
		return nil, fmt.Errorf("Bind failure\n")
	}

	if err := setDevFlags(name, syscall.IFF_PROMISC); err != nil {
		syscall.Close(s)
		return nil, err
	}
	return &Sock{fd: s, name: name}, nil
}

func (s *Sock) Name() string {
	return s.name
}

func (s *Sock) Addr() ([]byte, error) {
	var req ifreq_sockaddr
	copy(req.name[:syscall.IFNAMSIZ-1], s.name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.fd),
		syscall.SIOCGIFHWADDR, uintptr(unsafe.Pointer(&req))); errno != 0 {
		return nil, errno
	}
	return req.addr.addr[:], nil
}

func (s *Sock) Read(buf []byte) (int, error) {
	return syscall.Read(s.fd, buf)
}
func (s *Sock) Write(buf []byte) (int, error) {
	return syscall.Write(s.fd, buf)
}
func (s *Sock) Close() error {
	return syscall.Close(s.fd)
}

func getIFIndex(name string) (int32, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [16]byte
		index int32
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(soc),
		syscall.SIOCGIFINDEX,
		uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return 0, errno
	}
	return ifreq.index, nil
}
