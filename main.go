package main

import (
	"time"

	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
	"github.com/cook-a-doodle-do/my_protocol_stack/link"
	"github.com/cook-a-doodle-do/my_protocol_stack/link/arp"
	"github.com/cook-a-doodle-do/my_protocol_stack/link/ethernet"
	"github.com/cook-a-doodle-do/my_protocol_stack/network"
	"github.com/cook-a-doodle-do/my_protocol_stack/network/icmp"
	"github.com/cook-a-doodle-do/my_protocol_stack/network/ipv4"
	"github.com/cook-a-doodle-do/my_protocol_stack/raw"
)

func main() {
	//ethernetデバイスを作る
	raw, err := raw.New(raw.TAP)
	if err != nil {
		panic(err)
	}
	linkdev, err := ethernet.NewDevice(raw)
	if err != nil {
		panic(err)
	}
	//ethernetデヴァイスをlink層のデヴァイスとして登録する
	link.AppendDevice(linkdev)

	//リンク層にARPを登録
	link.RegistProtocol(enums.EtherTypeARP, arp.CallbackHandler)
	//ネットワーク層にIPv4を登録
	network.RegistProtocol(enums.EtherTypeIPv4, ipv4.CallbackHandler)

	//network層のデヴァイスを作る
	netdev, err := network.NewDevice(linkdev)
	if err != nil {
		panic(err)
	}

	//IPインターフェイスを作製，アドレス追加
	ipv4IF := ipv4.NewInterface(netdev)
	ipv4IF.SetIPAddr([4]byte{10, 0, 0, 2})
	ipv4IF.SetNetMask([4]byte{255, 255, 255, 255})

	netdev.AppendInterface(ipv4IF)
	ipv4.RegistProtocol(ipv4.ProtocolTypeICMP, icmp.CallbackHandler)

	time.Sleep(3 * time.Minute)
}
