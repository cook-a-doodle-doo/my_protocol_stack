package ethernet

import (
	"github.com/cook-a-doodle-do/my_protocol_stack/link"
	"github.com/cook-a-doodle-do/my_protocol_stack/raw"
)

const (
	MTU        = 1500
	HeaderSize = 6 + 6 + 2
)

var (
	BroadcastAddr = link.HardWareAddr{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 6}
)

type Device struct {
	raw  raw.Device
	addr link.HardWareAddr
}

type Header struct {
	destination link.HardWareAddr
	sender      link.HardWareAddr
	etherType   uint16
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
	var d *Device
	d.raw = rd
	d.addr.Len = 6
	copy(addr, d.addr.Buf[:])
	return d, nil
}

func (d *Device) Type() link.HardWareType {
	return link.HardwareTypeEthernet
}

func (d *Device) Name() string {
	return d.Name()
}

func (d *Device) Addr() link.HardWareAddr {
	return d.addr
}

func (d *Device) BroadcastAddr() link.HardWareAddr {
	return BroadcastAddr
}

func (d *Device) MTU() uint {
	return MTU
}

func (d *Device) HeaderSize() uint {
	return HeaderSize
}

func (d *Device) Read(buf []byte) (int, error) {
	n, err := d.Read(buf)
	return n, err
}

func (d *Device) Write(buf []byte) (int, error) {
	n, err := d.Write(buf)
	return n, err
}

func (d *Device) Close() error {
	return d.Close()
}
