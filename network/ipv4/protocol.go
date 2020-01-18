package ipv4

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/network"
)

type header struct {
	VersionAndHeaderLength uint8
	ServiceType            uint8
	TotalLength            uint16
	Identification         uint16
	FlagsAndFragmentOffset uint16
	TimeToLive             uint8
	ProtocolNum            ProtocolNum
	HeaderChecksum         uint16
	SourceIPAddr           [4]byte
	DestinationIPAddr      [4]byte
}

func (h *header) Version() uint {
	return uint(h.VersionAndHeaderLength >> 4)
}

func (h *header) HeaderLength() uint {
	return uint(h.VersionAndHeaderLength&0x0f) * 4
}

func (h *header) Flags() uint {
	return uint(h.FlagsAndFragmentOffset >> 13)
}

func (h *header) GetFlagmentOffset() uint {
	return uint(h.FlagsAndFragmentOffset&OFFMASK) << 3
}

func (h *header) SetFlagmentOffset(val uint16) {
	h.FlagsAndFragmentOffset = (h.FlagsAndFragmentOffset & ^OFFMASK) | (val >> 3)
}

//ProtocolRxHandler
func CallbackHandler(d *network.Device, payload []byte) {
	fmt.Println("<< ipv4 rx ===================== >>")
	//	fmt.Println(hex.Dump(payload))
	hdr, err := parseHeader(payload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(hdr)
	fmt.Printf("version        :%d\n", hdr.Version())
	fmt.Printf("header len     :%d\n", hdr.HeaderLength())
	fmt.Printf("service type   :%b\n", hdr.ServiceType)
	fmt.Printf("total len      :%d\n", hdr.TotalLength)
	fmt.Printf("identification :%d\n", hdr.Identification)
	fmt.Printf("flags          :%b\n", hdr.Flags())
	fmt.Printf("time to live   :%d\n", hdr.TimeToLive)
	fmt.Printf("protocolNum    :%s\n", hdr.ProtocolNum.Name())
	fmt.Printf("HeaderChecksum :%d\n", hdr.HeaderChecksum)
	fmt.Println(hdr.SourceIPAddr)
	fmt.Println(hdr.DestinationIPAddr)

	fmt.Println("<< ipv4 tx ===================== >>")
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

//===========================================================================
//上位プロトコル番号
type ProtocolNum uint8

const (
	ProtocolICMP ProtocolNum = 1
	ProtocolIGMP ProtocolNum = 2
	ProtocolIP   ProtocolNum = 4
	ProtocolTCP  ProtocolNum = 6
	ProtocolEGP  ProtocolNum = 8
	ProtocolUDP  ProtocolNum = 17
	ProtocolIPv6 ProtocolNum = 41
	ProtocolRSVP ProtocolNum = 46
	ProtocolOSPF ProtocolNum = 89
)

func (p ProtocolNum) Name() string {
	switch p {
	case ProtocolICMP:
		return "ICMP"
	case ProtocolIGMP:
		return "IGMP"
	case ProtocolIP:
		return "IP"
	case ProtocolTCP:
		return "TCP"
	case ProtocolEGP:
		return "EGP"
	case ProtocolUDP:
		return "UDP"
	case ProtocolIPv6:
		return "IPv6"
	case ProtocolRSVP:
		return "RSVP"
	case ProtocolOSPF:
		return "OSPF"
	}
	return "UnKnown"
}
