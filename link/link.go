//リンク層のプロトコルそのものに関するパッケージ
package link

import (
	"fmt"
	"io"
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
	Send(ProtocolType, HardwareAddr, []byte) error
	io.ReadWriteCloser
}

//リンク層の全デバイスがここに入る============================================
var (
	devices        map[string]Device                = make(map[string]Device)
	upperProtocols map[ProtocolType]ProtocolHandler = make(map[ProtocolType]ProtocolHandler)
)

func HaveHardware(d HardwareType) bool {
	return true
}

func HaveProtocol(p ProtocolType) bool {
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

type RxHandler func(Device, HardwareAddr, HardwareAddr, ProtocolType, []byte)

func rxHandler(dev Device, dst, src HardwareAddr, upt ProtocolType, payload []byte) {
	s := src.Entity()
	d := dst.Entity()
	fmt.Printf("src: %x:%x:%x:%x:%x:%x\n", s[0], s[1], s[2], s[3], s[4], s[5])
	fmt.Printf("dst: %x:%x:%x:%x:%x:%x\n", d[0], d[1], d[2], d[3], d[4], d[5])
	fmt.Printf("type:%s\n\n", upt.Name())
	rx, ok := upperProtocols[upt]
	if !ok {
		fmt.Println("protocol", upt.Name(), "is not implmented!")
	}
	rx(dev, payload)
}

func Devices() map[string]Device {
	return devices
}

func RegistProtocol(upt ProtocolType, up ProtocolHandler) {
	upperProtocols[upt] = up
}

type ProtocolHandler func(Device, []byte)

type ProtocolType uint

const (
	ProtocolType_IPv4  = iota
	ProtocolType_IPv6  = iota
	ProtocolType_ARP   = iota
	ProtocolType_RARP  = iota
	ProtocolType_UnDef = iota
)

func (p ProtocolType) Name() string {
	switch p {
	case ProtocolType_IPv4:
		return "IPv4"
	case ProtocolType_IPv6:
		return "IPv6"
	case ProtocolType_ARP:
		return "ARP"
	case ProtocolType_RARP:
		return "RARP"
	}
	return "UnDef"
}

//============================================================================
/*ネットワークインターフェイス(デバイスをどう扱うか)
type NetIF struct {
	ups    []*Protocol
	device *Device
}
*/

/*
func (u ProtocolType) EtherType() string {
	switch u {
	case ProtocolType_IPv4:
		return 0x0800
	case ProtocolType_IPv6:
		return 0x86dd
	case ProtocolType_ARP:
		return 0x0806
	case ProtocolType_RARP:
		return 0x8635
	default:
		return 0x0000
	}
}

func (u ProtocolType) EtherType() [6]byte {
	switch u {
	case ProtocolType_ICMP:
		return "ICMP"
	case ProtocolType_IPv4:
		return "IPv4"
	case ProtocolType_IPv6:
		return "IPv6"
	case ProtocolType_ARP:
		return "ARP"
	case ProtocolType_RARP:
		return "RAPR"
	default:
		return "undefined"
	}
}

*/
