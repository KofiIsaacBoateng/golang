package main

import (
	"bytes"
	"distributed-fs/p2p"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
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

type Message struct {
	Payload any
}

type FileMessage struct {
	Key string
	Size int64
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


func (s *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...);


	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	var (
		fileBuf = new(bytes.Buffer)
		tee = io.TeeReader(r, fileBuf)
	)

	n, err := s.store.Write(key, tee);
	if err != nil {
		return err
	}



	msgBuf := new(bytes.Buffer)
	msg := Message{
		Payload: FileMessage{
			Key: key,
			Size: n,
		},
	}

	// encode Message struct into a buffer
	if err := gob.NewEncoder(msgBuf).Encode(&msg); err != nil {
		return err
	}

	for _ , peer := range(s.peers) {
		if err := peer.Send(msgBuf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 2)

	for _, peer := range(s.peers) {
		_, err := io.Copy(peer, fileBuf);
		if err != nil {
			return err
		}
	}

	return nil

	// // 1. store this file to disk
	// buf := new(bytes.Buffer);
	// tee := io.TeeReader(r, buf);

	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }

	// // 2. Broadcast to all known peers on the network

	// msg := MessageData{
	// 	Key: key,
	// 	Data: buf.Bytes(),
	// }

	// return s.broadcast(&Message{
	// 	From: "todo",
	// 	Payload: msg,
	// })
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


func (s *FileServer) MessageLoop() {
	defer func(){
		log.Println("Server closed by user quitch action!")
		s.Transport.Close()
	}()
	
	for {
		select{
		case rpc :=  <- s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				fmt.Println("Payload error: ", err)
				continue
			}

			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println(err)
			}

		case <- s.quitch:
			return
		}
	}
}


func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case FileMessage:
		s.handleFileMessage(from, v)
	}
	return nil
}


func (s *FileServer) handleFileMessage(from string, msg FileMessage) error {
	peer, ok := s.peers[from];
	if(!ok) {
		return fmt.Errorf("Peer %s not found in peer map!\n", from)
	}

		
	n, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size)); 
	if err != nil {
		return err
	}
	fmt.Printf("Written (%d)bytes to disk\n", n)

	peer.(*p2p.TCPPeer).Wg.Done();

	return nil
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

	s.MessageLoop()

	return nil
}



func init () {
	gob.Register(FileMessage{})
}
