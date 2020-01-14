package ipv4

import (
	"fmt"

	"github.com/cook-a-doodle-do/my_protocol_stack/network"
)

//ProtocolRxHandler
func CallbackHandler(d *network.Device, payload []byte, src, dst network.HardwareAddr) {
	fmt.Println("IP_CALLBACK_HANDLER")
}
