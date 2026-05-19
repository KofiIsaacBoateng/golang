package main

import (
	"bytes"
	"distributed-fs/p2p"
	"encoding/gob"
	"fmt"
	"io"
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

type Payload struct {
	Key string
	Data []byte
}


func (s *FileServer) broadcast(p *Payload) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...);


	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. store this file to disk
	buf := new(bytes.Buffer);
	tee := io.TeeReader(r, buf);

	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	// 2. Broadcast to all known peers on the network

	p := &Payload{
		Key: key,
		Data: buf.Bytes(),
	}

	return s.broadcast(p)
}

func (s *FileServer) Stop() {
	close(s.quitch);
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.peers[peer.RemoteAddr().String()] = peer

	fmt.Println("connected to peer network: ", peer.RemoteAddr())

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
			var p Payload
			r := bytes.NewReader(msg.Payload)
			decoder := gob.NewDecoder(r)
			if err := decoder.Decode(&p); err != nil {
				fmt.Println("Payload error: ", err)
			}
			fmt.Printf("Message received from peer: %+v\n", p)
		case <- s.quitch:
			return
		}
	}
}


func (s *FileServer) BootStrapNetwork() error {
	for _, addr := range s.BootStrapNodes {
		go func(addr string) {
			fmt.Println("Attempting to connect with remote peer: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				fmt.Printf("Dial error on addr: %s [%+v]\n", addr, err)
			}
		}(addr)
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
