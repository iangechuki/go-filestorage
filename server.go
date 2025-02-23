package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/IANGECHUKI176/go-filestorage/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}
type FileServer struct {
	FileServerOpts

	peerLock sync.RWMutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		PathTransformFunc: opts.PathTransformFunc,
		Root:              opts.StorageRoot,
	}
	return &FileServer{
		FileServerOpts: opts,
		peers:          make(map[string]p2p.Peer),
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
	}
}

type Message struct {
	From    string
	Payload any
}
type MessageStoreFile struct {
	Key  string
	Size int64
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}
	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
func (s *FileServer) StoreData(key string, r io.Reader) error {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)
	size, err := s.store.Write(key, tee)
	if err != nil {
		return err
	}
	log.Printf("wrote %d bytes to disk\n", size)

	msg := &Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	if err := s.broadcast(msg); err != nil {
		return err
	}
	time.Sleep(time.Second * 3)
	// TODO: use multiwriter
	for _, peer := range s.peers {
		n, err := io.Copy(peer, buf)
		if err != nil {
			return err
		}
		log.Printf("wrote %d bytes to %s\n", n, peer.RemoteAddr().String())
	}
	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}
func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("peer connected: %s\n", p.RemoteAddr().String())
	return nil

}
func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Fatal(err)
			}
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println(err)
				return
			}

			// if err := s.handleMessage(&m); err != nil {
			// 	log.Println(err)
			// }
			// fmt.Printf("Message %+v\n", msg)
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, v)

	}
	return nil
}
func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer not found")
	}
	if _, err := s.store.Write(msg.Key, io.LimitReader(peer, int64(msg.Size))); err != nil {
		return err
	}
	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}
func (s *FileServer) bootstrapNetwork() error {
	log.Println("bootstraping the network")
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("attempting to connect with remote: ", addr)
			err := s.Transport.Dial(addr)
			if err != nil {
				log.Println("dial error: ", err)
			}

		}(addr)
	}
	return nil
}
func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	s.loop()
	return nil
}

func (s *FileServer) Store(key string, r io.Reader) (int64, error) {
	return s.store.Write(key, r)
}
func init() {
	gob.Register(MessageStoreFile{})
}
