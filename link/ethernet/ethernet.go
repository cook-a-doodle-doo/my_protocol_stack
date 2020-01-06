package ethernet

import (
	"bytes"
	"encoding/binary"
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
		fmt.Println("flame is destroyed! because address is not mine.\n")
		return
	}
	f(dst, src, hdr.EtherType.UpperProtocolType(), flame[HeaderSize:])
}

func (e EtherType) UpperProtocolType() link.UpperProtocolType {
	switch e {
	case 0x0800:
		return link.UpperProtocolType_IPv4
	case 0x0806:
		return link.UpperProtocolType_ARP
	case 0x8635:
		return link.UpperProtocolType_RARP
	case 0x86dd:
		return link.UpperProtocolType_IPv6
	default:
		return link.UpperProtocolType_UnDef
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
