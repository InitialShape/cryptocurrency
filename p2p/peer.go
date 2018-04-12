package p2p

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"github.com/InitialShape/blockchain/blockchain"
	"io/ioutil"
	"os"
	"bufio"
)

const (
	CONN_TYPE = "tcp"
)

type Peer struct {
	Port string
	Host string
	Store blockchain.Store
}

func (p *Peer) RegisterDefaultPeers() error {
	file, err := os.Open("peers.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		p.RegisterPeer(scanner.Text())
	}
	return err
}


func (p *Peer) Start() {
	p.RegisterDefaultPeers()
	address := fmt.Sprintf("%s:%s", p.Port, p.Host)
	listener, err := net.Listen(CONN_TYPE, address)
	if err != nil {
		log.Fatal("Error listening: ",err)
	}
	defer listener.Close()
	fmt.Printf("Peer is listening on %s:%s\n", p.Port, p.Host)

	go p.CheckHeartBeat()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
		}

		go p.Handle(conn)
	}
}

func (p *Peer) Handle(conn net.Conn) {
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err)
	}

	req := string(buf)
	resp := []byte{}

	fmt.Printf("Received message from: %s %s\n", conn.RemoteAddr().String(), req)

	// Buffer is 1024 big, even after casting to string
	if strings.Contains(req, "PING") {
		resp = p.Pong()
	}
	if strings.Contains(req, "REGISTER") {
		resp = p.RegisterPeer(conn.RemoteAddr().String())
	}

	conn.Write(resp)
	conn.Close()
}

func (p *Peer) Pong() []byte {
	return []byte("PONG")
}

func (p *Peer) RegisterPeer(peer string) []byte {
	p.Store.AddPeer(peer)
	fmt.Println("Registered new peer: ", peer)
	return []byte("REGISTERED")
}

func (p *Peer) Ping(peer string) error {
	conn, err := net.Dial(CONN_TYPE, peer)
	if err != nil {
		fmt.Println("Error dialing peer: ", peer, err)
		fmt.Println("Removing peer: ", peer)
		p.Store.DeletePeer(peer)
		return err
	}

	conn.Write([]byte("PING"))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		p.Store.DeletePeer(peer)
		return err
	}

	if string(resp) != "PONG" {
		p.Store.DeletePeer(peer)
	} else {
		fmt.Println("Received message from: ", conn.RemoteAddr().String(), string(resp))
	}
	return err
}

func (p *Peer) CheckHeartBeat() {
	for range time.Tick(time.Second * 15) {
		peers, err := p.Store.GetPeers()
		if err != nil {
			log.Fatal("Error getting peers: ", err)
		}
		for _, peer := range peers {
			p.Ping(peer)
		}
	}
}
