package ethernet

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
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

type header struct {
	Destination MacAddr
	Source      MacAddr
	EtherType   enums.EtherType
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
	copy(d.addr[:], addr)
	return &d, nil
}

func (d *Device) Type() enums.HardwareType {
	return enums.HardwareTypeEthernet
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
	fmt.Println("<< ethernet tx ================= >>")
	fmt.Println(hex.Dump(buf))
	n, err := d.raw.Write(buf)
	return n, err
}

func (d *Device) RxHandler(flame []byte, f link.RxHandler) {
	fmt.Println("<< ethernet rx ================= >>")
	fmt.Println(hex.Dump(flame))
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
	f(d, dst, src, hdr.EtherType, flame[HeaderSize:])
}

func (d *Device) Send(et enums.EtherType, hrd link.HardwareAddr, buf []byte) error {
	h := header{
		EtherType: et,
	}
	//TODO 送信サイズのチェック
	copy(h.Destination[:], hrd.Entity())
	copy(h.Source[:], d.Addr().Entity())

	b := new(bytes.Buffer)
	err := binary.Write(b, binary.BigEndian, h)
	if err != nil {
		return err
	}
	_, err = b.Write(buf)
	if err != nil {
		return err
	}

	_, err = d.Write(b.Bytes())
	if err != nil {
		return err
	}
	return nil
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
