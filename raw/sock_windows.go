package raw

import "fmt"

type Sock struct {
}

func NewSock() *Sock {
	fmt.Println("New")
	return &Sock{}
}

func (s *Sock) Read(buf []byte) (int, error) {
	fmt.Println("Readed")
	return 0, nil
}
func (s *Sock) Write(buf []byte) (int, error) {
	fmt.Println("Write")
	return 0, nil
}
func (s *Sock) Close() error {
	return nil
}
