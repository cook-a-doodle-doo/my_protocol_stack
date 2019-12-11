package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/cook-a-doodle-do/my_protocol_stack/raw_device"
)

func main() {
	tap, err := raw_device.NewTapLinux("tap%d")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer tap.Close()
	buf := make([]byte, 1500)

	s, err := tap.Addr()
	if err != nil {
		fmt.Println("can't get mac addr")
	}
	fmt.Printf("%x:%x:%x:%x:%x:%x\n", s[0], s[1], s[2], s[3], s[4], s[5])
	fmt.Println("start tap0")
	for i := 0; ; i++ {
		//		bufio.NewScanner(os.Stdin).Scan()
		//		time.Sleep(time.Second)
		n, _ := tap.Read(buf)
		fmt.Printf("\nupdate~~~~~~~~~~~~~~~~~~~~~~~~~~~~%d\n", i)
		fmt.Println(hex.Dump(buf[:n]))
	}
}
