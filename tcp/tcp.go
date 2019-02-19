package tcp

import (
	"fmt"
	"net"
)

type Tcp struct {
	Ip      string
	Portnum string
	conn    net.Conn
}

//Write tcp写入
func (t *Tcp) Write(b []byte) (n int, err error) {
	// t.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 10))
	return t.conn.Write(b)
}

//Read tcp读
func (t *Tcp) Read(b []byte) (n int, err error) {
	fmt.Println("tttttttt")
	return t.conn.Read(b)
}

type Nats struct {
	Ip      string
	Portnum string
	conn    net.Conn
}

//Write tcp写入
func (t *Nats) Write(b []byte) (n int, err error) {
	// t.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 10))
	return t.conn.Write(b)
}

//Read tcp读
func (t *Nats) Read(b []byte) (n int, err error) {
	fmt.Println("tttttttt")
	return t.conn.Read(b)
}
