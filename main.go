package main

import (
	"fmt"
	"log"

	"github.com/cook-a-doodle-do/my_protocol_stack/layer2"
	netdev "github.com/cook-a-doodle-do/my_protocol_stack/net_device"
)

func main() {
	tap, err := netdev.Open("TAP", "mytap0")
	if err != nil {
		log.Fatal(err.Error())
	}

	eth := layer2.New(tap, "ETHERNET")
	fmt.Printf("%+v", eth)
	//
	for {
		tap.Tx([]byte("hogera"))
	}

}
