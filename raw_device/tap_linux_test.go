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
	defer tap.Close()
	buf := make([]byte, 100)
	t.Log("\nstart\n")

	for true {
		time.Sleep(time.Second * 1)
		t.Log("hoge\n")
		t.Log(string(buf))
		break
	}
}
