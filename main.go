package main

import (
	"fmt"
	"time"

	"github.com/cook-a-doodle-do/my_protocol_stack/link"
	_ "github.com/cook-a-doodle-do/my_protocol_stack/link/ethernet"
)

func main() {
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
	fmt.Println(link.Devices())
	time.Sleep(2 * time.Minute)
}
