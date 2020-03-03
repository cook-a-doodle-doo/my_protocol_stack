package network

import (
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
	"github.com/cook-a-doodle-do/my_protocol_stack/link"
)

//network層のデバイス
//リンク層のデバイスにnetwork層独自のデータ(IPアドレスとか)を紐付ける
type Device struct {
	link.Device
	IFs map[enums.EtherType][]Interface
}

//リンクデバイスに紐付いているネットデバイスが欲しい
var devices map[link.Device]*Device = make(map[link.Device]*Device)

//リンクデバイスに1対1で紐づくネットワークデバイスを作製
func NewDevice(link link.Device) (*Device, error) {
	d, ok := devices[link]
	if ok {
		return nil, fmt.Errorf("The link device is already allocated.\n")
	}
	d = &Device{
		Device: link,
		IFs:    make(map[enums.EtherType][]Interface),
	}
	devices[link] = d
	return d, nil
}

type HardwareAddr []byte

func (h HardwareAddr) Entity() []byte {
	b := make([]byte, 6)
	copy(b, h[:])
	return b
}

func (h HardwareAddr) Length() uint {
	return uint(len(h))
}

//networkデバイスに論理インターフェイスを紐付ける
func (d *Device) AppendInterface(i Interface) {
	d.IFs[i.EtherType()] = append(d.IFs[i.EtherType()], i)
	link.AppendInterface(
		d.Device,
		&LinkInterface{
			netIF: i,
		})
}

type LinkInterface struct {
	netIF Interface
}

func (l *LinkInterface) ProtocolAddr() link.ProtocolAddr {
	return l.netIF.ProtocolAddr()
}

func (l *LinkInterface) EtherType() enums.EtherType {
	return l.netIF.EtherType()
}

const (
	RxQueueSize = 10
)

//network層のデバイスを受け取る(同時にインターフェイスも)
type ProtocolRxHandler func(*Device, []byte)

func RegistProtocol(et enums.EtherType, f ProtocolRxHandler) {
	type packet struct {
		link    link.Device
		payload []byte
	}
	rxQueue := make(chan packet, RxQueueSize)

	//リンク層から呼び出しを受けたらキューに追加する
	link.RegistProtocol(
		et,
		func(link link.Device, payload []byte, dst, src link.HardwareAddr) {
			rxQueue <- packet{
				link:    link,
				payload: payload,
			}
		})
	//キューからひたすら読んでハンドラにぶっこむ
	go func() {
		for {
			p := <-rxQueue
			f(devices[p.link], p.payload)
		}
	}()
}

type ProtocolAddr interface {
	Entity() []byte
	Length() uint
}

type ProtocolNum uint8

const (
	ProtocolNumICMP      = 1   //Internet Control Message
	ProtocolNumIGMP      = 2   //Internet Group Management
	ProtocolNumIP        = 4   //IP in IP ( encapsulation )
	ProtocolNumTCP       = 6   //Transmission Control
	ProtocolNumCBT       = 7   //CBT
	ProtocolNumEGP       = 8   //Exterior Gateway Protocol
	ProtocolNumIGP       = 9   //any private interior gateway
	ProtocolNumUDP       = 17  //User Datagram
	ProtocolNumIPv6      = 41  //Ipv6
	ProtocolNumIPv6Route = 43  //Routing Header for IPv6
	ProtocolNumIPv6Frag  = 44  //Fragment Header for IPv6
	ProtocolNumIDRP      = 45  //Inter-Domain Routing Protocol
	ProtocolNumRSVP      = 46  //Reservation Protocol
	ProtocolNumGRE       = 47  //General Routing Encapsulation
	ProtocolNumESP       = 50  //Encap Security Payload
	ProtocolNumAH        = 51  //Authentication Header
	ProtocolNumMOBILE    = 55  //IP Mobility
	ProtocolNumIPv6ICMP  = 58  //ICMP for IPv6
	ProtocolNumIPv6NoNxt = 59  //No Next Header for IPv6
	ProtocolNumIPv6Opts  = 60  //Destination Options for IPv6
	ProtocolNumEIGRP     = 88  //EIGRP
	ProtocolNumOSPF      = 89  //OSPF
	ProtocolNumIPIP      = 94  //IP-within-IP Encapsulation Protocol
	ProtocolNum103PIM    = 103 //Protocol Independent Multicast
	ProtocolNum112VRRP   = 112 //Virtual Router Redundancy Protocol
	ProtocolNum113PGM    = 113 //PGM Reliable Transport Protocol
	ProtocolNum115L2TP   = 115 //Layer Two Tunneling Protocol ProtocolNum
)

//network層の何らかの情報が入る
//IPアドレスとか
type Interface interface {
	ProtocolAddr() ProtocolAddr
	EtherType() enums.EtherType
	Tx(pn ProtocolNum, data []byte, dst ProtocolAddr) error
}
