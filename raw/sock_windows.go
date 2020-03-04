package raw

import "fmt"

type Sock struct {
}

func NewSock(name string) (*Sock, error) {
	fmt.Println("New")
	return &Sock{}, nil
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
