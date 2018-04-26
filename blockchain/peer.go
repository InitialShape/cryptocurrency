package blockchain

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// TODO: Use log for logging

const (
	CONN_TYPE = "tcp"
)

type Peer struct {
	Port  string
	Host  string
	Store Store
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
		peer := scanner.Text()
		p.RegisterPeer(peer)
	}
	return err
}

func (p *Peer) Start() {
	p.RegisterDefaultPeers()
	go p.Discovery()
	address := fmt.Sprintf(":%s", p.Port)
	listener, err := net.Listen(CONN_TYPE, address)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
	defer listener.Close()
	log.Printf("Peer is listening on %s:%s\n", p.Host, p.Port)

	go p.CheckHeartBeat()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
		}

		go p.Handle(conn)
	}
}

func (p *Peer) Handle(conn net.Conn) {
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading: ", err)
	}

	req := string(buf[:n])
	resp := []byte{}

	log.Printf("Received message from: %s %s\n", conn.RemoteAddr().String(), req)

	if strings.Contains(req, "PING") {
		peers := strings.Split(req, " ")
		if len(peers) > 1 {
			p.RegisterPeer(peers[1])
		}
		resp = p.Pong()
	}
	if strings.Contains(req, "PEERS") {
		peers, err := p.GetPeers()
		if err != nil {
			log.Println("Error getting peers on request: ", err)
		}
		resp = peers
	}
	if strings.Contains(req, "TRANSACTION") {
		// in new function also check for index
		transactionJSON := strings.Split(req, " ")[1]
		var transaction Transaction
		err := json.Unmarshal([]byte(transactionJSON), &transaction)
		if err != nil {
			log.Println("Couldn't read transaction JSON: ", err)
		}
		err = p.Store.AddTransaction(transaction)
		fmt.Println("Added new transaction: ", transactionJSON)
	}
	if strings.Contains(req, "BLOCK") {
		// in new function also check for index
		blockJSON := strings.Split(req, " ")[1]
		var block Block
		err := json.Unmarshal([]byte(blockJSON), &block)
		if err != nil {
			fmt.Println("Couldn't read block JSON: ", err)
		}
		err = p.Store.AddBlock(block)
		log.Println("Added new block: ", blockJSON)
	}
	if strings.Contains(req, "CHAIN") {
		blocks, err := p.GetChain()
		if err != nil {
			log.Fatal("Error getting chain", err)
		}
		resp = blocks
	}

	conn.Write(resp)
	conn.Close()
}

func (p *Peer) GetChain() ([]byte, error) {
	blocks, err := p.Store.GetChain()
	if err != nil {
		log.Println("Error getting chain from store", err)
		return []byte{}, err
	}

	blocksJSON, err := json.Marshal(blocks)
	if err != nil {
		log.Println("Error marshalling blocks to JSON", err)
		return []byte{}, err
	}

	return blocksJSON, err
}

func (p *Peer) GetPeers() ([]byte, error) {
	peers, err := p.Store.GetPeers()
	if err != nil {
		return []byte{}, err
	}
	return []byte(strings.Join(peers, "\n")), err
}

func (p *Peer) Pong() []byte {
	return []byte("PONG")
}

func (p *Peer) RegisterPeer(peer string) []byte {
	self := fmt.Sprintf("%s:%s", p.Port, p.Host)
	if peer != self {
		p.Store.AddPeer(peer)
		fmt.Println("Registered new peer: ", peer)
		return []byte("REGISTERED")
	}
	return []byte("NOT REGISTERED")
}

func (p *Peer) DiscoverPeers(peer string) error {
	fmt.Println("Requesting new peers from: ", peer)
	conn, err := net.Dial(CONN_TYPE, peer)
	if err != nil {
		fmt.Println("Error dialing peer: ", peer, err)
		fmt.Println("Removing peer: ", peer)
		p.Store.DeletePeer(peer)
		return err
	}
	conn.Write([]byte("PEERS"))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("Error reading PEERS response: ", err)
		return err
	}
	respString := string(resp)
	peers := strings.Split(respString, "\n")
	fmt.Println("New peers received: ", peers)
	for _, peer := range peers {
		p.RegisterPeer(peer)
	}

	return err
}

func (p *Peer) SendTransaction(peer string, transaction Transaction) error {
	conn, err := net.Dial(CONN_TYPE, peer)
	if err != nil {
		fmt.Println("Error dialing peer on sending transaction: ", err)
		fmt.Println("Deleting peer: ", peer)
		p.Store.DeletePeer(peer)
		return err
	}

	header := []byte("TRANSACTION ")
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		log.Fatal("Error marshalling transaction: ", err)
	}
	message := append(header, transactionJSON...)

	conn.Write(message)
	return err
}

func (p *Peer) DownloadChain(peer string) ([]Block, error) {
	var chain []Block
	conn, err := net.Dial(CONN_TYPE, peer)
	if err != nil {
		log.Println("Error dialing peer, deleting peer: ", peer)
		p.Store.DeletePeer(peer)
		return chain, err
	}

	msg := []byte("CHAIN")
	conn.Write(msg)

	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Println("Error reading response, deleting peer: ", peer)
		p.Store.DeletePeer(peer)
		return chain, err
	}
	err = json.Unmarshal(resp, &chain)
	if err != nil {
		log.Println("Couldn't read chain JSON: ", err)
	}

	return chain, err
}

func (p *Peer) Ping(peer string) error {
	conn, err := net.Dial(CONN_TYPE, peer)
	if err != nil {
		log.Println("Error dialing peer, Deleting peer: ", peer)
		p.Store.DeletePeer(peer)
		return err
	}
	msg := fmt.Sprintf("PING %s:%s", p.Port, p.Host)

	conn.Write([]byte(msg))
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

func (p *Peer) Download() ([][]Block, error) {
	var chains [][]Block
	peers, err := p.Store.GetPeers()
	if err != nil {
		log.Fatal("Error getting peers: ", err)
		return chains, err
	}

	// TODO: Change this to a smaller set of peers when the network scales
	//		 to more nodes than just a hand full.
	for _, peer := range peers {
		chain, err := p.DownloadChain(peer)
		if err != nil {
			log.Println("Couldn't download chain from", peer)
		} else {
			chains = append(chains, chain)
		}
	}
	return chains, err
}

func (p *Peer) Discovery() error {
	for range time.Tick(time.Second * 15) {
		fmt.Println("Peer discovery initialized")
		peers, err := p.Store.GetPeers()
		if err != nil {
			log.Fatal("Error getting peers: ", err)
			return err
		}
		for _, peer := range peers {
			p.DiscoverPeers(peer)
		}
	}
	return errors.New("Cannot be reached")
}

func (p *Peer) SendPayload(peer string, header []byte, payload []byte) error {
	conn, err := net.Dial(CONN_TYPE, peer)
	if err != nil {
		fmt.Println("Error dialing peer on sending payload: ", err)
		fmt.Println("Deleting peer: ", peer)
		p.Store.DeletePeer(peer)
		return err
	}

	message := append(header, payload...)
	conn.Write(message)
	return err
}

func (p *Peer) GossipTransaction(transaction Transaction) {
	peers, err := p.Store.GetPeers()
	if err != nil {
		fmt.Println("Error getting peers: ", err)
	}
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		log.Fatal("Error marshalling transaction: ", err)
	}

	for _, peer := range peers {
		p.SendPayload(peer, []byte("TRANSACTION "), transactionJSON)
	}
}

func (p *Peer) GossipBlock(block Block) {
	peers, err := p.Store.GetPeers()
	if err != nil {
		fmt.Println("Error getting peers: ", err)
	}
	blockJSON, err := json.Marshal(block)
	if err != nil {
		log.Fatal("Error marshalling transaction: ", err)
	}

	for _, peer := range peers {
		p.SendPayload(peer, []byte("BLOCK "), blockJSON)
	}
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
