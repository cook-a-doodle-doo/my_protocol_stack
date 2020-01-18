package arp

import (
	"bytes"
	"reflect"
	"time"

	"github.com/cook-a-doodle-do/my_protocol_stack/enums"
)

//実は中身はウンチで出来ているというね…
type Table []*info //map[enums.HardwareType][]info

//ゴミのようなネーミング
type info struct {
	ht   enums.HardwareType
	ha   HardwareAddr
	pt   enums.EtherType
	pa   ProtocolAddr
	last time.Time
}

//実体を作る
var table Table //= make(map[enums.HardwareType][]info)

func (t Table) Insert(pt enums.EtherType, pa ProtocolAddr, ht enums.HardwareType, ha HardwareAddr) {
	table = append(table, &info{
		ht:   ht,
		ha:   ha,
		pt:   pt,
		pa:   pa,
		last: time.Now(),
	})
}

func (t Table) Set(pt enums.EtherType, pa ProtocolAddr, ht enums.HardwareType, ha HardwareAddr) {
	for _, v := range t {
		if pt == v.pt &&
			ht == v.ht &&
			reflect.DeepEqual(v, pa) {
			v.ha = ha
			v.last = time.Now()
			return
		}
	}
}

//プロトコルタイプ，プロトコルアドレス，ハードウェアタイプの組をキーとしてハードウェアアドレスを引く
func (t Table) Get(pt enums.EtherType, pa ProtocolAddr, ht enums.HardwareType) (HardwareAddr, bool) {
	for _, v := range t {
		if pt == v.pt &&
			ht == v.ht &&
			bytes.Equal(v.pa, pa) {
			return v.ha, true
		}
	}
	return nil, false
}
