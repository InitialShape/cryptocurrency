package p2p

import (
	"fmt"
	"log"
	"net"
	"bytes"
)

const (
	CONN_TYPE = "tcp"
)

type Peer struct {
	Port string
	Host string
}

func (p *Peer) Start() {
	address := fmt.Sprintf("%s:%s", p.Port, p.Host)
	listener, err := net.Listen(CONN_TYPE, address)
	if err != nil {
		log.Fatal("Error listening: ",err)
	}
	defer listener.Close()
	fmt.Printf("Peer is listening on %s:%s", p.Port, p.Host)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
		}

		go p.handle(conn)
	}
}

func (p *Peer) handle(conn net.Conn) {
	var buf bytes.Buffer

	_, err := conn.Read(buf.Bytes())
	if err != nil {
		fmt.Println("Error reading: ", err)
	}

	conn.Write([]byte("Message received\n"))
	conn.Close()
}
