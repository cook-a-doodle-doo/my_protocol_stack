//Address Resolution Protocol
//上位アドレスをキーとしてHardWareAddressを検索する
package arp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/link"
)

const (
	Operation_REQUEST = 1
	Operation_REPLAY  = 2
)

const (
	Format_Ethernet = 1
)

const (
	HeaderSize = 8
)

type HardwareAddr []byte
type ProtocolAddr []byte

func (h HardwareAddr) Entity() []byte {
	return h
}
func (h HardwareAddr) Length() uint {
	return uint(len(h))
}

var arpTable map[*ProtocolAddr]*HardwareAddr = make(map[*ProtocolAddr]*HardwareAddr)

type header struct {
	HardwareAddrFormat uint16
	ProtocolAddrFormat uint16
	HardwareAddrLength uint8
	ProtocolAddrLength uint8
	Operation          uint16
}

type params struct {
	SenderHA HardwareAddr
	SenderPA ProtocolAddr
	TargetHA HardwareAddr
	TargetPA ProtocolAddr
}

func parseHeader(buf []byte) (*header, error) {
	var hdr header
	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	return &hdr, nil
}

func parseParams(payload []byte, hlen, plen uint8) (*params, error) {
	p := params{
		SenderHA: make([]byte, hlen),
		SenderPA: make([]byte, plen),
		TargetHA: make([]byte, hlen),
		TargetPA: make([]byte, plen),
	}
	reader := bytes.NewReader(payload)
	err := binary.Read(reader, binary.BigEndian, &p.SenderHA)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &p.SenderPA)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &p.TargetHA)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &p.TargetPA)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func RxHandler(dev link.Device, buf []byte) {
	fmt.Println("<< arp ================= >>")
	fmt.Println(hex.Dump(buf[HeaderSize:]))
	hdr, err := parseHeader(buf)
	if err != nil {
		fmt.Println("invalid header!")
		return
	}
	p, err := parseParams(buf[HeaderSize:], hdr.HardwareAddrLength, hdr.ProtocolAddrLength)
	if err != nil {
		fmt.Println("invalid payload!")
		return
	}
	fmt.Println(hdr.Operation)
	fmt.Println(hdr.HardwareAddrFormat)
	fmt.Println(hdr.ProtocolAddrFormat)
	fmt.Println(p.SenderHA)
	fmt.Println(p.SenderPA)
	fmt.Println(p.TargetHA)
	fmt.Println(p.TargetPA)
	fmt.Println("hoge")
	//指定のハードウェアを扱えない→終了
	if !link.HaveHardware(link.HardwareType(hdr.HardwareAddrFormat)) {
		return
	}
	//指定のプロトコルを扱えない→終了
	//	t := ethernet.ProtocolType2EtherType(ethernet.EtherType(hdr.ProtocolAddrFormat))
	//	if !link.HaveProtocol() {
	//		return
	//	}
	//Marge_flag を false にする
	//もし、プロトコルタイプと送信者プロトコルアドレスが既に ARP テーブ ルに含まれていたら、
	//
	//    そのエントリのハードウェアアドレスを新しいアドレスに更新し、
	//    Marge_flag を true にする
	//
	//自分がターゲットプロトコルアドレスでない→終了
	//もし、 Marge_flag が false なら、「プロトコル、送信者プロトコルア ドレス、送信者ハードウェアアドレス」を ARP テーブルに追加する
	//OP code が Request でない→終了
	if hdr.Operation != Operation_REQUEST {
		return
	}
	//送信者欄とターゲット欄を交換し、自ホストのハードウェアアドレスとプ ロトコルアドレスを送信者欄に記入する
	//OP code を Reply にする
	//このパケットが送られてきたホストに対して、作成した ARP パケットを 送信する

	//	dev.Send(link.ProtocolType_ARP, p.SenderHA, buf)
}
