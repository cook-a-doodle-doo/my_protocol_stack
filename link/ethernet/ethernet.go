package ethernet

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/link"
	"github.com/cook-a-doodle-do/my_protocol_stack/raw"
)

const (
	MTU        = 1500
	HeaderSize = 6 + 6 + 2
	MacAddrLen = 6
)

type MacAddr [6]byte

func (m MacAddr) Entity() []byte {
	b := make([]byte, 6)
	copy(b, m[:])
	return b
}

func (m MacAddr) Length() uint {
	return MacAddrLen
}

var (
	BroadcastAddr = MacAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

type Device struct {
	raw  raw.Device
	addr MacAddr
}

type EtherType uint16

type header struct {
	Destination MacAddr
	Source      MacAddr
	EtherType   EtherType
}

func init() {
	dev, err := NewDevice()
	if err != nil {
		panic(err)
	}
	if err := link.RegistDevice(dev); err != nil {
		panic(err)
	}
}

func NewDevice() (*Device, error) {
	rd, err := raw.New(raw.TAP)
	if err != nil {
		return nil, err
	}
	addr, err := rd.Addr()
	if err != nil {
		return nil, err
	}
	var d Device
	d.raw = rd
	copy(addr, d.addr[:])
	return &d, nil
}

func (d *Device) Type() link.HardwareType {
	return link.HardwareTypeEthernet
}

func (d *Device) Name() string {
	return d.raw.Name()
}

func (d *Device) Addr() link.HardwareAddr {
	return &d.addr
}

func (d *Device) BroadcastAddr() link.HardwareAddr {
	return BroadcastAddr
}

func (d *Device) MTU() uint {
	return MTU
}

func (d *Device) HeaderSize() uint {
	return HeaderSize
}

func (d *Device) Read(buf []byte) (int, error) {
	n, err := d.raw.Read(buf)
	return n, err
}

func (d *Device) Write(buf []byte) (int, error) {
	n, err := d.Write(buf)
	return n, err
}

func (d *Device) RxHandler(flame []byte, f link.RxHandler) {
	fmt.Println("<< ethernet flame ================= >>")
	hdr, err := parseHeader(flame)
	if err != nil {
		return
	}
	var (
		dst link.HardwareAddr = hdr.Destination
		src link.HardwareAddr = hdr.Source
	)
	//TODO 理解が怪しい
	if dst != d.addr && dst != BroadcastAddr && (dst.Entity()[0]&0x01 != 0) {
		dst := dst.Entity()
		src := src.Entity()
		fmt.Println("flame is destroyed! because address is not mine.")
		fmt.Printf("dst: %x:%x:%x:%x:%x:%x\n", dst[0], dst[1], dst[2], dst[3], dst[4], dst[5])
		fmt.Printf("src: %x:%x:%x:%x:%x:%x\n\n", src[0], src[1], src[2], src[3], src[4], src[5])
		return
	}
	fmt.Println(hex.Dump(flame[HeaderSize:]))
	f(d, dst, src, hdr.EtherType.ProtocolType(), flame[HeaderSize:])
}

func (d *Device) Send(t link.ProtocolType, hrd link.HardwareAddr, buf []byte) error {
	h := header{
		EtherType: ProtocolType2EtherType(t),
	}
	copy(h.Destination[:], hrd.Entity())
	copy(h.Source[:], d.Addr().Entity())
	//	Destination MacAddr
	//	Source      MacAddr
	//	EtherType   EtherType
	return nil
}

func ProtocolType2EtherType(p link.ProtocolType) EtherType {
	switch p {
	case link.ProtocolType_IPv4:
		return 0x0800
	case link.ProtocolType_ARP:
		return 0x0806
	case link.ProtocolType_RARP:
		return 0x8635
	case link.ProtocolType_IPv6:
		return 0x86dd
	default:
		return 0x0000 //link.ProtocolType_UnDef
	}
}

func (e EtherType) ProtocolType() link.ProtocolType {
	switch e {
	case 0x0800:
		return link.ProtocolType_IPv4
	case 0x0806:
		return link.ProtocolType_ARP
	case 0x8635:
		return link.ProtocolType_RARP
	case 0x86dd:
		return link.ProtocolType_IPv6
	default:
		return link.ProtocolType_UnDef
	}
}

func (d *Device) Close() error {
	return d.Close()
}

func parseHeader(buf []byte) (*header, error) {
	var hdr header
	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	return &hdr, nil
}
