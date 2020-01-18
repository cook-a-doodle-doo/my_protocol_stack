//リンク層のプロトコルそのものに関するパッケージ
package link

import (
	"fmt"
	"io"

	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
)

type HardwareAddr interface {
	Entity() []byte
	Length() uint
}

type ProtocolAddr interface {
	Entity() []byte
	Length() uint
}

type Device interface {
	Type() enums.HardwareType
	Name() string
	Addr() HardwareAddr
	BroadcastAddr() HardwareAddr //Broadcast Address
	MTU() uint                   //Maximum Transmission Unit
	HeaderSize() uint
	RxHandler([]byte, RxHandler)
	Send(enums.EtherType, HardwareAddr, []byte) error
	io.ReadWriteCloser
}

var (
	//リンク層の全デバイスがここに入る============================================
	devices map[enums.HardwareType][]Device = make(map[enums.HardwareType][]Device)
	//プロトコルを登録
	protocols map[enums.EtherType]ProtocolHandler = make(map[enums.EtherType]ProtocolHandler)
	//デバイスにインターフェイスを登録
	interfaces map[Device][]Interface = make(map[Device][]Interface)
)

type Interface interface {
	ProtocolAddr() ProtocolAddr
	EtherType() enums.EtherType
}

func AppendInterface(d Device, i Interface) {
	interfaces[d] = append(interfaces[d], i)
}

//登録した瞬間に動き始める====================================================
func AppendDevice(d Device) {
	devices[d.Type()] = append(devices[d.Type()], d)
	go func() {
		for {
			buf := make([]byte, d.MTU()+d.HeaderSize())
			n, err := d.Read(buf)
			if err != nil {
				panic(err)
			}
			d.RxHandler(buf[:n], rxHandler)
		}
	}()
}

func Interfaces(d Device) []Interface {
	i, _ := interfaces[d]
	return i
}

type RxHandler func(Device, HardwareAddr, HardwareAddr, enums.EtherType, []byte)

func rxHandler(dev Device, dst, src HardwareAddr, upt enums.EtherType, payload []byte) {
	s := src.Entity()
	d := dst.Entity()
	fmt.Printf("src: %x:%x:%x:%x:%x:%x\n", s[0], s[1], s[2], s[3], s[4], s[5])
	fmt.Printf("dst: %x:%x:%x:%x:%x:%x\n", d[0], d[1], d[2], d[3], d[4], d[5])
	fmt.Printf("type:%s\n\n", upt.Name())
	rx, ok := protocols[upt]
	if !ok {
		fmt.Println("protocol", upt.Name(), "is not implmented!")
	}
	rx(dev, payload, dst, src)
}

func Devices(ht enums.HardwareType) []Device {
	return devices[ht]
}

func Protocols() map[enums.EtherType]ProtocolHandler {
	return protocols
}

func RegistProtocol(upt enums.EtherType, up ProtocolHandler) {
	protocols[upt] = up
}

//device, payload, dst, src
type ProtocolHandler func(Device, []byte, HardwareAddr, HardwareAddr)
