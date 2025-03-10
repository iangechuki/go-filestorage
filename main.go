package main

import (
	"bytes"
	"log"
	"time"

	"github.com/IANGECHUKI176/go-filestorage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	TCPTransport := p2p.NewTCPTransport(tcpOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         TCPTransport,
		BootstrapNodes:    nodes,
	}
	s := NewFileServer(fileServerOpts)
	TCPTransport.OnPeer = s.OnPeer
	return s
}
func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(1 * time.Second)

	go s2.Start()
	time.Sleep(1 * time.Second)

	data := bytes.NewReader([]byte("new data to  be stored"))
	s2.StoreData("Key", data)
	select {}
}
