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

type HardwareType uint

//ハードウェアタイプをここに追記
const (
	HardwareTypeEthernet = iota
	HardwareTypeLoopBack = iota
)

type Device interface {
	Type() HardwareType
	Name() string
	Addr() HardwareAddr
	BroadcastAddr() HardwareAddr //Broadcast Address
	MTU() uint                   //Maximum Transmission Unit
	HeaderSize() uint
	RxHandler([]byte, RxHandler)
	Send(enums.EtherType, HardwareAddr, []byte) error
	io.ReadWriteCloser
}

//リンク層の全デバイスがここに入る============================================
var (
	devices        map[string]Device                   = make(map[string]Device)
	upperProtocols map[enums.EtherType]ProtocolHandler = make(map[enums.EtherType]ProtocolHandler)
)

func HaveHardware(d HardwareType) bool {
	return true
}

func HaveProtocol(p enums.EtherType) bool {
	_, ok := upperProtocols[p]
	return ok
}

//登録した瞬間に動き始める====================================================
//インターフェイスがないので細かい設定は出来ない
func RegistDevice(d Device) error {
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
	return nil
}

type RxHandler func(Device, HardwareAddr, HardwareAddr, enums.EtherType, []byte)

func rxHandler(dev Device, dst, src HardwareAddr, upt enums.EtherType, payload []byte) {
	s := src.Entity()
	d := dst.Entity()
	fmt.Printf("src: %x:%x:%x:%x:%x:%x\n", s[0], s[1], s[2], s[3], s[4], s[5])
	fmt.Printf("dst: %x:%x:%x:%x:%x:%x\n", d[0], d[1], d[2], d[3], d[4], d[5])
	fmt.Printf("type:%s\n\n", upt.Name())
	rx, ok := upperProtocols[upt]
	if !ok {
		fmt.Println("protocol", upt.Name(), "is not implmented!")
	}
	rx(dev, payload, dst, src)
}

func Devices() map[string]Device {
	return devices
}

func RegistProtocol(upt enums.EtherType, up ProtocolHandler) {
	upperProtocols[upt] = up
}

//device, payload, dst, src
type ProtocolHandler func(Device, []byte, HardwareAddr, HardwareAddr)
