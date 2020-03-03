package icmp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/network"
	"github.com/cook-a-doodle-do/my_protocol_stack/network/ipv4"
)

type format struct {
	Type     Type
	Code     uint8
	Chacksum uint16
}

type info struct {
	format
	iface   network.Interface
	src     ipv4.IPAddr
	dst     ipv4.IPAddr
	payload []byte
}

func CallbackHandler(iface network.Interface, src, dst ipv4.IPAddr, payload []byte) {
	fmt.Println("<< icmp rx ===================== >>")
	fmt.Println(hex.Dump(payload))
	var f format
	reader := bytes.NewReader(payload)
	if err := binary.Read(reader, binary.BigEndian, &f); err != nil {
		return
	}
	cs, err := CalcChacksum(payload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if cs != f.Chacksum {
		fmt.Println("invarid chacksum")
		return
	}

	f.Type.FunCall(info{
		format:  f,
		iface:   iface,
		dst:     dst,
		src:     src,
		payload: payload[:4],
	})

	fmt.Println("<< icmp tx ===================== >>")
}

type Type uint8

const (
	TypeEchoReply              Type = 0
	TypeDestinationUnreachable Type = 3
	TypeSourceQuenchMessage    Type = 4
	TypeRedirect               Type = 5
	TypeEchoRequest            Type = 8
	TypeRouterAvertisement     Type = 9
	TypeRouterSolicitation     Type = 10
	TypeTimeExceeded           Type = 11
	TypeParameterProblem       Type = 12
	TypeTimestamp              Type = 13
	TypeTimestampReqpay        Type = 14
	TypeInformationRequest     Type = 15
	TypeInformationReply       Type = 16
	TypeAddressMaskRequest     Type = 17
	TypeAddressMaskReqly       Type = 18
	TypeTranceroute            Type = 30
)

func (t Type) Name() string {
	msg, ok := messages[t]
	if !ok {
		return "Unknown Type"
	}
	return msg.name
}

func (t Type) FunCall(info info) error {
	msg, ok := messages[t]
	if !ok {
		return fmt.Errorf("Unknown Type")
	}
	fmt.Println(msg.name)
	return msg.callback(info)
}

type message struct {
	name     string
	callback func(info) error
}

var messages = map[Type]message{
	TypeEchoReply: message{
		name: "Echo Reply",
		callback: func(info info) error {
			return nil
		}},
	TypeDestinationUnreachable: message{
		name: "Destination Unreachable",
		callback: func(info info) error {
			return nil
		}},
	TypeEchoRequest: message{
		name: "Echo Request",
		callback: func(info info) error {
			fmt.Println(info)
			//	iface := info.iface
			return nil
		}},
}

func CalcChacksum(buf []byte) (uint16, error) {
	if len(buf) < 4 {
		return 0, fmt.Errorf("invarid length")
	}
	//オプションを除き全部加算
	var acc uint32
	for i := 0; i < len(buf); i += 2 {
		acc += uint32(buf[i]) * 0x0100
		if i > len(buf)-1 {
			break
		}
		acc += uint32(buf[i+1])
	}
	//チェックサムの影響を除く
	acc -= uint32(buf[2]) * 0x0100
	acc -= uint32(buf[3])
	//桁上がりを1桁目に加算し16bit化
	var v uint32 = acc / 0x10000
	acc &= 0x0000ffff
	acc += v + (acc+v)/0x10000
	acc &= 0x0000ffff
	return uint16(acc) ^ 0xffff, nil
}

/*
	switch t {
	case TypeEchoReply:
		return "Echo Reply"
	case TypeDestinationUnreachable:
		return "Destination Unreachable"
	case SourceQuenchMessage:
		return "SourceQuenchMessage"
	case TypeRedirect:
		return "Redirect"
	case TypeEchoRequest:
		return "Echo Request"
	case TypeRouterAvertisement:
		return "TypeRouterAvertisement"
	case TypeRouterSolicitation:
		return "TypeRouterSolicitation"
	case TypeTimeExceeded:
		return "TypeTimeExceeded"
	case TypeParameterProblem:
		return "TypeParameterProblem"
	case TypeTimestamp:
		return "TypeTimestamp"
	case TypeTimestampReqpay:
		return "TypeTimestampReqpay"
	case TypeInformationRequest:
		return "TypeInformationRequest"
	case TypeInformationReply:
		return "TypeInformationReply"
	case TypeAddressMaskRequest:
		return "TypeAddressMaskRequest"
	case TypeAddressMaskReqly:
		return "TypeAddressMaskReqly"
	case TypeTranceroute:
		return "TypeTranceroute"
	}
	return "Un Known"
}
*/
