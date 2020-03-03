package ipv4

//ルーティング，フラグメントについてのサポートの無い通信を提供
//フラグメンテーションメカニズムは提供する
//「パケット紛失?? 順番?? 知らんがな」 という感じ
//目的のIPアドレス届けるように最大限努力しますという感じ (出来るとは言っていない)

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
	"github.com/cook-a-doodle-do/my_protocol_stack/network"
)

type IPAddr [4]byte

const (
	IPAddrSize uint = 4
)

type header struct {
	VersionAndHeaderLength uint8
	ServiceType            uint8
	TotalLength            uint16
	Identification         uint16
	FlagsAndFragmentOffset uint16
	TimeToLive             uint8
	ProtocolType           ProtocolType
	HeaderChecksum         uint16
	Source                 IPAddr
	Destination            IPAddr
}

func (h *header) Version() uint {
	return uint(h.VersionAndHeaderLength >> 4)
}

//byte
func (h *header) HeaderLength() uint {
	return uint(h.VersionAndHeaderLength&0x0f) * 4
}

func (h *header) Flags() uint16 {
	return uint16(h.FlagsAndFragmentOffset)
}

func (h *header) GetFlagmentOffset() uint16 {
	return uint16(h.FlagsAndFragmentOffset&OFFMASK) << 3
}

func (h *header) SetFlagmentOffset(val uint16) {
	h.FlagsAndFragmentOffset = (h.FlagsAndFragmentOffset & ^OFFMASK) | (val >> 3)
}

func flagsDump(flag uint16) string {
	//TODO ゴミ
	str := ""
	if 0 != Reserved&flag {
		str += "Reserved"
	}
	str += " | "
	if 0 != Dont&flag {
		str += "Dont"
	}
	str += " | "
	if 0 != More&flag {
		str += "More"
	}
	return str
}

func (h *header) Println() {
	fmt.Printf("version        :%d\n", h.Version())
	fmt.Printf("header len     :%d\n", h.HeaderLength())
	fmt.Printf("service type   :%b\n", h.ServiceType)
	fmt.Printf("total len      :%d\n", h.TotalLength)
	fmt.Printf("identification :%d\n", h.Identification)
	fmt.Printf("flags          :%s\n", flagsDump(h.Flags()))
	fmt.Printf("time to live   :%d\n", h.TimeToLive)
	fmt.Printf("protocolNum    :%s\n", h.ProtocolType.Name())
	fmt.Printf("HeaderChecksum :%d\n", h.HeaderChecksum)
	fmt.Println(h.Source)
	fmt.Println(h.Destination)
}

func CalcChacksum(buf []byte) (uint16, error) {
	if len(buf) < 20 {
		return 0, fmt.Errorf("invarid length")
	}
	//オプションを除き全部加算
	var acc uint32
	for i := 0; i < 20; i += 2 {
		acc += uint32(buf[i]) * 0x0100
		acc += uint32(buf[i+1])
	}
	//チェックサムの影響を除く
	acc -= uint32(buf[10]) * 0x0100
	acc -= uint32(buf[11])
	//桁上がりを1桁目に加算し16bit化
	var v uint32 = acc / 0x10000
	acc &= 0x0000ffff
	acc += v + (acc+v)/0x10000
	acc &= 0x0000ffff
	return uint16(acc) ^ 0xffff, nil
}

const (
	EtherType enums.EtherType = enums.EtherTypeIPv4
)

var (
	protocols map[ProtocolType]RxHandler = make(map[ProtocolType]RxHandler)
)

// payload
type RxHandler func(network.Interface, IPAddr, IPAddr, []byte)

func RegistProtocol(num ProtocolType, rx RxHandler) {
	protocols[num] = rx
}

//Callback~ 外部から呼び出してもらうための関数
//networkパッケージから呼び出してもらう
func CallbackHandler(d *network.Device, payload []byte) {
	fmt.Println("<< ipv4 rx ===================== >>")
	//そもそもpayloadに中身はあるの???
	if len(payload) < 20 {
		fmt.Println("too short payload.")
		return
	}
	//ヘッダ・分けま・しょうね
	hdr, err := parseHeader(payload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	hdr.Println()

	//チェックサムを調べる
	cksum, err := CalcChacksum(payload)
	if err != nil {
		return
	}
	if cksum != hdr.HeaderChecksum {
		fmt.Println("chacksum did not match.")
		return
	}

	//バージョンが合っているか
	if hdr.Version() != 4 {
		fmt.Println("this IPdatagram is not ipv4.")
		return
	}

	//オタクくんさぁ…まだ生きてる??
	if hdr.TimeToLive < 0 {
		fmt.Println("this IPdatagram has timed out.")
		return
	}

	rx, ok := protocols[hdr.ProtocolType]
	if !ok {
		fmt.Printf("%s is not registerd\n", hdr.ProtocolType.Name())
		return
	}
	//インターフェイスが目的のアドレスを持っていればハンドラを呼び出し
	for _, v := range d.IFs[EtherType] {
		//TODO Broadcast
		if v.ProtocolAddr() != hdr.Destination {
			break
		}
		fmt.Println("")
		rx(v, hdr.Source, hdr.Destination, payload[hdr.HeaderLength():])
	}
}

func parseHeader(buf []byte) (*header, error) {
	var hdr header
	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	return &hdr, nil
}

const (
	Reserved uint16 = 0x8000 //0が格納される
	Dont     uint16 = 0x4000
	More     uint16 = 0x2000
)
const (
	OFFMASK uint16 = 0x1fff
)

func send(src, dst IPAddr,
	prot ProtocolType,
	tos uint8,
	ttl uint8,
	buf []byte,
	Id uint8,
	DF uint8,
	opt uint8) error {
	//	hdr := header{}

	fmt.Println("<< ipv4 tx ===================== >>")
	return nil
}

//       src =送信元アドレス
//       dst =宛先アドレス
//       prot =プロトコル
//       TOS =サービスの種類
//       TTL =存続時間
//       BufPTR =バッファーポインター
//       len =バッファの長さ
//       Id =識別子
//       DF =フラグメント化しない
//       opt =オプションデータ
//      結果=応答
//         OK =データグラムは正常に送信されました
//        エラー=引数のエラーまたはローカルネットワークエラー

//===========================================================================
//上位プロトコル番号
type ProtocolType uint8

const (
	ProtocolTypeICMP ProtocolType = 1
	ProtocolTypeIGMP ProtocolType = 2
	ProtocolTypeIP   ProtocolType = 4
	ProtocolTypeTCP  ProtocolType = 6
	ProtocolTypeEGP  ProtocolType = 8
	ProtocolTypeUDP  ProtocolType = 17
	ProtocolTypeIPv6 ProtocolType = 41
	ProtocolTypeRSVP ProtocolType = 46
	ProtocolTypeOSPF ProtocolType = 89
)

func (p ProtocolType) Name() string {
	switch p {
	case ProtocolTypeICMP:
		return "ICMP"
	case ProtocolTypeIGMP:
		return "IGMP"
	case ProtocolTypeIP:
		return "IP"
	case ProtocolTypeTCP:
		return "TCP"
	case ProtocolTypeEGP:
		return "EGP"
	case ProtocolTypeUDP:
		return "UDP"
	case ProtocolTypeIPv6:
		return "IPv6"
	case ProtocolTypeRSVP:
		return "RSVP"
	case ProtocolTypeOSPF:
		return "OSPF"
	}
	return "UnKnown"
}
