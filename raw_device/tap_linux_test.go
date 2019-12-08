package raw_device

import (
	"testing"
	"time"
)

func TestHelloWorld2(t *testing.T) {
	// t.Fatal("not implemented")
	tap, err := NewTapLinux("tap%d")
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(time.Second * 30)
	tap.Close()
}
