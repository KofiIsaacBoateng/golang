package main

import (
	"distributed-fs/p2p"
	"fmt"
	"log"
	"sync"
)

type FileServerOpts struct {
	StorageRoot   string
	Transport     p2p.Transport
	PathTransformerFunc PathTransformerFunc
	BootStrapNodes []string
}

type FileServer struct {
	FileServerOpts

	mu sync.Mutex
	peers map[string]p2p.Peer

	store *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		PathTransformerFunc: opts.PathTransformerFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch: make(chan struct{}),
		peers: make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Stop() {
	close(s.quitch);
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("connected to peer network: ", peer.RemoteAddr())
	s.peers[peer.RemoteAddr().String()] = peer

	return nil
}


func (s *FileServer) Loop() {
	defer func(){
		log.Println("Server closed by user quitch action!")
		s.Transport.Close()
	}()
	for {
		select{
		case msg :=  <- s.Transport.Consume():
			fmt.Printf("Message received from peer: %+v\n", msg)
		case <- s.quitch:
			return
		}
	}
}


func (s *FileServer) BootStrapNetwork() error {
	for _, addr := range s.BootStrapNodes {
		go func(){
			fmt.Println("Attempting to connect with remote peer: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				fmt.Printf("Dial error on addr: %s [%+v]\n", addr, err)
			}
		}()
	}

	return nil
}


func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootStrapNodes) != 0{
		s.BootStrapNetwork()
	}

	s.Loop()

	return nil
}
