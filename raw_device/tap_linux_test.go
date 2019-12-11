package raw_device

import (
	"testing"
)

func TestHelloWorld2(t *testing.T) {
	// t.Fatal("not implemented")
	tap, err := NewTapLinux("tap%d")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer tap.Close()
	//	buf := make([]byte, 100)
}
