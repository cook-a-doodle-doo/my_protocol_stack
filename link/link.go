//リンク層のプロトコルそのものに関するパッケージ
package link

import (
	"encoding/hex"
	"fmt"
	"io"
)

type HardWareAddr struct {
	Buf []byte
	Len uint
}

type HardWareType int

//ハードウェアタイプをここに追記
const (
	HardwareTypeEthernet = iota
	HardwareTypeLoopBack = iota
)

type RxHandler func(dst, src HardWareAddr)

type Device interface {
	Type() HardWareType
	Name() string
	Addr() HardWareAddr
	BroadcastAddr() HardWareAddr //Broadcast Address
	MTU() uint                   //Maximum Transmission Unit
	HeaderSize() uint
	io.ReadWriteCloser
	/*
		Read([]byte) (int, error)
		Write([]byte) (int, error)
		Close() error
	*/
}

//リンク層の全デバイスがここに入る============================================
var devices map[string]*Device = make(map[string]*Device)

//登録した瞬間に動き始める====================================================
//インターフェイスがないので細かい設定は出来ない
func RegistDevice(d Device) error {
	//rxloop
	go func() {
		for {
			buf := make([]byte, d.MTU()+d.HeaderSize())
			n, err := d.Read(buf)
			fmt.Println("<< ethernet flame ================= >>")
			fmt.Println(hex.Dump(buf[:n]))
			if err != nil {
				panic(err)
			}
		}
	}()
	return nil
}

/*
func (d Device) Regist(up UpperProtocol) {
}
*/

type UpperProtocol interface {
	EtherType()
	RxHandler()
}

//============================================================================
/*ネットワークインターフェイス(デバイスをどう扱うか)
type NetIF struct {
	ups    []*UpperProtocol
	device *Device
}
*/
