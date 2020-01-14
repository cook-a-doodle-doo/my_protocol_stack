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
}

const (
	RxQueueSize = 10
)

//network層のデバイスを受け取る(同時にインターフェイスも)
type ProtocolRxHandler func(*Device, []byte, HardwareAddr, HardwareAddr)

func RegistProtocol(et enums.EtherType, f ProtocolRxHandler) {
	type packet struct {
		link    link.Device
		payload []byte
		dst     link.HardwareAddr
		src     link.HardwareAddr
	}
	rxQueue := make(chan packet, RxQueueSize)

	link.RegistProtocol(
		et,
		//リンク層から呼び出しを受けたらキューに追加する
		func(link link.Device, payload []byte, dst, src link.HardwareAddr) {
			rxQueue <- packet{
				link:    link,
				payload: payload,
				dst:     dst,
				src:     src,
			}
		})
	//キューからひたすら読んでハンドラにぶっこむ
	go func() {
		for {
			p := <-rxQueue
			fmt.Println(p)
			d := p.dst.Entity()
			s := p.src.Entity()
			f(devices[p.link], p.payload, d, s)
		}
	}()
}

type ProtocolAddr interface {
	Entity() []byte
	Length() uint
}

//network層の何らかの情報が入る
//IPアドレスとかね
type Interface interface {
	ProtocolAddr() ProtocolAddr
	EtherType() enums.EtherType
}
