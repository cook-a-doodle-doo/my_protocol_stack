//リンク層のプロトコルそのものに関するパッケージ
package link

import (
	"encoding/hex"
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
	io.ReadWriteCloser
	/*
		Read([]byte) (int, error)
		Write([]byte) (int, error)
		Close() error
	*/
}

//リンク層の全デバイスがここに入る============================================
var (
	devices        map[string]Device          = make(map[string]Device)
	upperProtocols map[string][]UpperProtocol = make(map[string][]UpperProtocol)
)

//登録した瞬間に動き始める====================================================
//インターフェイスがないので細かい設定は出来ない
func RegistDevice(d Device) error {
	//rxloop
	fmt.Println(d.Name())
	go func() {
		for {
			buf := make([]byte, d.MTU()+d.HeaderSize())
			n, err := d.Read(buf)
			if err != nil {
				panic(err)
			}
			fmt.Println("<< ethernet flame ================= >>")
			fmt.Println(hex.Dump(buf[:n]))
			if err != nil {
				panic(err)
			}
			d.RxHandler(buf, rxHandler)
		}
	}()
	return nil
}

type RxHandler func(HardwareAddr, HardwareAddr, UpperProtocolType, []byte)

func rxHandler(dst, src HardwareAddr, upt UpperProtocolType, payload []byte) {
	s := src.Entity()
	d := dst.Entity()
	fmt.Printf("src: %x:%x:%x:%x:%x:%x\n", s[0], s[1], s[2], s[3], s[4], s[5])
	fmt.Printf("dst: %x:%x:%x:%x:%x:%x\n", d[0], d[1], d[2], d[3], d[4], d[5])
	fmt.Printf("type:%s\n", upt.Name())
	for key, protocols := range upperProtocols {
		for _, protocol := range protocols {
			protocol.RxHandler(devices[key], payload)
		}
	}
}

func Devices() map[string]Device {
	return devices
}

func RegistUpperProtocol(d Device, up UpperProtocol) {
	upperProtocols[d.Name()] = append(upperProtocols[d.Name()], up)
	up.RegistLinkDevice(d)
}

type UpperProtocol interface {
	RegistLinkDevice(Device)
	RxHandler(Device, []byte)
	Type() UpperProtocolType
}

type UpperProtocolType uint

const (
	UpperProtocolType_IPv4  = iota
	UpperProtocolType_IPv6  = iota
	UpperProtocolType_ARP   = iota
	UpperProtocolType_RARP  = iota
	UpperProtocolType_UnDef = iota
)

func (u UpperProtocolType) Name() string {
	switch u {
	case UpperProtocolType_IPv4:
		return "IPv4"
	case UpperProtocolType_IPv6:
		return "IPv6"
	case UpperProtocolType_ARP:
		return "ARP"
	case UpperProtocolType_RARP:
		return "RARP"
	}
	return "UnDef"
}

//============================================================================
/*ネットワークインターフェイス(デバイスをどう扱うか)
type NetIF struct {
	ups    []*UpperProtocol
	device *Device
}
*/

/*
func (u UpperProtocolType) EtherType() string {
	switch u {
	case UpperProtocolType_IPv4:
		return 0x0800
	case UpperProtocolType_IPv6:
		return 0x86dd
	case UpperProtocolType_ARP:
		return 0x0806
	case UpperProtocolType_RARP:
		return 0x8635
	default:
		return 0x0000
	}
}

func (u UpperProtocolType) EtherType() [6]byte {
	switch u {
	case UpperProtocolType_ICMP:
		return "ICMP"
	case UpperProtocolType_IPv4:
		return "IPv4"
	case UpperProtocolType_IPv6:
		return "IPv6"
	case UpperProtocolType_ARP:
		return "ARP"
	case UpperProtocolType_RARP:
		return "RAPR"
	default:
		return "undefined"
	}
}

*/
