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

	fmt.Println("start tap0")
	for i := 0; ; i++ {
		//		bufio.NewScanner(os.Stdin).Scan()
		//		time.Sleep(time.Second)
		fmt.Println("\nupdate~~~~~~~~~~~~~~~~~~~~~~~~~~~~%d", i)
		n, _ := tap.Read(buf)
		fmt.Println(hex.Dump(buf[:n]))
	}
}
