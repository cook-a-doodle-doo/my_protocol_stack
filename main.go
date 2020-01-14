package main

import (
	"time"

	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
	"github.com/cook-a-doodle-do/my_protocol_stack/link"
	"github.com/cook-a-doodle-do/my_protocol_stack/link/arp"
	"github.com/cook-a-doodle-do/my_protocol_stack/link/ethernet"
	"github.com/cook-a-doodle-do/my_protocol_stack/network"
	"github.com/cook-a-doodle-do/my_protocol_stack/network/ipv4"
)

func main() {
	//ethernetデヴァイスを作る
	linkdev, err := ethernet.NewDevice()
	if err != nil {
		panic(err)
	}
	//ethernetデヴァイスをlink層のデヴァイスとして登録する
	if err := link.RegistDevice(linkdev); err != nil {
		panic(err)
	}

	//リンク層にARPを登録
	link.RegistProtocol(enums.EtherTypeARP, arp.RxHandler)
	//ネットワーク層にIPv4を登録
	network.RegistProtocol(enums.EtherTypeIPv4, ipv4.CallbackHandler)

	//IPインターフェイスを作製，アドレス追加
	ipv4IF := ipv4.NewInterface()
	ipv4IF.SetIPAddr([]byte{10, 0, 0, 2})
	ipv4IF.SetNetMask([]byte{255, 255, 255, 255})

	//network層のデヴァイスを作る
	netdev, err := network.NewDevice(linkdev)
	if err != nil {
		panic(err)
	}
	netdev.AppendInterface(ipv4IF)

	time.Sleep(3 * time.Minute)
}

//	raw, err := raw_device.New(raw_device.TAP)
//	if err!= nil {
//		return
//	}

//	tap, err := raw.New(raw.TAP)
//	if err != nil {
//		log.Fatalf(err.Error())
//	}
//	defer tap.Close()
//	buf := make([]byte, 1500)
//
//	s, err := tap.Addr()
//	if err != nil {
//		fmt.Println("can't get mac addr")
//	}
//	fmt.Printf("%x:%x:%x:%x:%x:%x\n", s[0], s[1], s[2], s[3], s[4], s[5])
//	fmt.Println("start tap0")
/*
	for i := 0; ; i++ {
		//		bufio.NewScanner(os.Stdin).Scan()
		//		time.Sleep(time.Second)
		n, _ := tap.Read(buf)
		fmt.Printf("\nupdate~~~~~~~~~~~~~~~~~~~~~~~~~~~~%d\n", i)
		fmt.Println(hex.Dump(buf[:n]))
	}
*/
