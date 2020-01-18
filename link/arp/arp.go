//Address Resolution Protocol
//上位アドレスをキーとしてHardWareAddressを検索する
package arp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
	"github.com/cook-a-doodle-do/my_protocol_stack/link"
)

const (
	OperationRequest = 1
	OperationReplay  = 2
)

const (
	HeaderSize = 8
)

type HardwareAddr []byte
type ProtocolAddr []byte

func (p ProtocolAddr) Entity() []byte {
	return p
}
func (p ProtocolAddr) Length() uint {
	return uint(len(p))
}

func (h HardwareAddr) Entity() []byte {
	return h
}
func (h HardwareAddr) Length() uint {
	return uint(len(h))
}

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

func RxHandler(dev link.Device, buf []byte, src, dst link.HardwareAddr) {
	fmt.Println("<< arp rx ====================== >>")
	fmt.Println(hex.Dump(buf))
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

	//指定のハードウェアを扱えない→終了
	devs := link.Devices(enums.HardwareType(hdr.HardwareAddrFormat))
	if devs == nil {
		fmt.Printf("can't use %s\n", enums.HardwareType(hdr.HardwareAddrFormat).Name())
		return
	}

	//指定のプロトコルを扱えない→終了
	_, ok := link.Protocols()[enums.EtherType(hdr.ProtocolAddrFormat)]
	if !ok {
		fmt.Printf("can't use %s\n", enums.EtherType(hdr.ProtocolAddrFormat).Name())
		return
	}

	//Marge_flag を false にする
	margeFlag := false

	//もし、プロトコルタイプと送信者プロトコルアドレスが既に ARP テーブ ルに含まれていたら、
	_, ok = table.Get(
		enums.EtherType(hdr.ProtocolAddrFormat), p.SenderPA,
		enums.HardwareType(hdr.HardwareAddrFormat),
	)
	if ok {
		//そのエントリのハードウェアアドレスを新しいアドレスに更新し、
		fmt.Println("arp table updated")
		table.Set(
			enums.EtherType(hdr.ProtocolAddrFormat), p.SenderPA,
			enums.HardwareType(hdr.HardwareAddrFormat), p.SenderHA,
		)
		//Marge_flag を true にする
		margeFlag = true
	}

	//自分がターゲットプロトコルアドレスでない→終了
	//var target HardwareAddr
	for _, d := range link.Devices(enums.HardwareType(hdr.HardwareAddrFormat)) {
		for _, v := range link.Interfaces(d) {
			if bytes.Equal(v.ProtocolAddr().Entity(), p.TargetPA) {
				goto DONE
			}
		}
	}
	fmt.Printf("protocol address %d is not mine.\n", p.TargetPA)
	return
DONE:

	//もし、 Marge_flag が false なら、「プロトコル、送信者プロトコルア ドレス、送信者ハードウェアアドレス」を ARP テーブルに追加する
	if !margeFlag {
		fmt.Println("inserted new addr.")
		table.Insert(
			enums.EtherType(hdr.ProtocolAddrFormat), p.SenderPA,
			enums.HardwareType(hdr.HardwareAddrFormat), p.SenderHA,
		)
	}
	//OP code が Request でない→終了
	if hdr.Operation != OperationRequest {
		fmt.Println("operation is not request.")
		return
	}

	/*
		fmt.Println(p.SenderHA)
		fmt.Println(p.SenderPA)
		fmt.Println(p.TargetHA)
		fmt.Println(p.TargetPA)
	*/

	//送信者欄とターゲット欄を交換し、自ホストのハードウェアアドレスとプロトコルアドレスを送信者欄に記入する
	tmp := p.TargetPA
	p.TargetHA = p.SenderHA
	p.TargetPA = p.SenderPA
	p.SenderHA = dev.Addr().Entity()
	p.SenderPA = tmp

	/*
		fmt.Println(p.SenderHA)
		fmt.Println(p.SenderPA)
		fmt.Println(p.TargetHA)
		fmt.Println(p.TargetPA)
	*/
	//OP code を Reply にする
	hdr.Operation = OperationReplay

	//このパケットが送られてきたホストに対し，作成したARPパケットを送信する
	b, err := Write(hdr, p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("<< arp tx ====================== >>")
	fmt.Println(hex.Dump(b))
	dev.Send(enums.EtherTypeARP, dev.BroadcastAddr(), b)
}

func Write(hdr *header, p *params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, hdr)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, p.SenderHA.Entity())
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, p.SenderPA.Entity())
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, p.TargetHA.Entity())
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, p.TargetPA.Entity())
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

//SendRequest()
