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

func (s *Sock) Addr() ([]byte, error) {
	return []byte{10, 0, 0, 1}, nil
}

func (s *Sock) Name() string {
	return s.name
}
