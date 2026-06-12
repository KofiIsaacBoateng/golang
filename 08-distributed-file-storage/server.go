package main

import (
	"bytes"
	"distributed-fs/p2p"
	"encoding/binary"
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

type StoreFileMessage struct {
	Key string
	Size int64
}

type GetFileMessage struct {
	Key string
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		PathTransformerFunc: opts.PathTransformerFunc,
		RootDir: opts.StorageRoot,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch: make(chan struct{}),
		peers: make(map[string]p2p.Peer),
	}
}


func (s *FileServer) stream(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...);


	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)

	// encode Message struct to a buffer
	if err := gob.NewEncoder(buf).Encode(&msg); err != nil {
		return err
	}

	for _ , peer := range(s.peers) {
		peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}


	return nil
}


func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Has(key) {
		log.Printf("[%s]: Serving file [%s] from local disk...\n", s.Transport.Addr(), key)
		_, r, err := s.store.Read(key)

		return r, err
	}


	log.Printf("[%s]: File of key path [%s] not found locally, fetching from network...\n", s.Transport.Addr(), key)

	msg := Message{
		Payload: GetFileMessage{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	time.Sleep(500 * time.Millisecond)


	for _, peer := range s.peers {
		var fileSize int64;
		binary.Read(peer, binary.LittleEndian, &fileSize)
		n, err := s.store.Write(key, io.LimitReader(peer, fileSize));
		if err != nil {
			return nil, err
		}
		
		log.Printf("[%s]: Received (%d) bytes over the network from %s: \n", s.Transport.Addr(), n, peer.RemoteAddr())
		peer.CloseStream()
	}

	_, r, err := s.store.Read(key)

	return r, err
}


func (s *FileServer) Store(key string, r io.Reader) error {
	var (
		fileBuf = new(bytes.Buffer)
		tee = io.TeeReader(r, fileBuf)
	)

	n, err := s.store.Write(key, tee);
	if err != nil {
		return err
	}

	log.Printf("Received and Saved (%d)bytes to disk\n", n)

	msg := Message{
		Payload: StoreFileMessage{
			Key: key,
			Size: n,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)


	// TODO: multiwriter here to peers
	for _, peer := range(s.peers) {
		peer.Send([]byte{p2p.IncomingStream})
		_, err := io.Copy(peer, fileBuf);
		if err != nil {
			return err
		}
	}

	return nil
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
		log.Println("Server closed by error or user quitch action!")
		s.Transport.Close()
	}()
	
	for {
		select{
		case rpc :=  <- s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				fmt.Println("Decode error(consume): ", err)
				continue
			}
			fmt.Println(msg)

			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("Handle message error: ", err)
			}

		case <- s.quitch:
			return
		}
	}
}


func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case StoreFileMessage:
		return s.handleStoreFileMessage(from, v)

	case GetFileMessage:
		return s.handleGetFileMessage(from, v)
	}
	return nil
}


func (s *FileServer) handleGetFileMessage(from string, msg GetFileMessage) error {
	if !s.store.Has(msg.Key) {
		return fmt.Errorf("[%s]: Need to serve file of key path [%s], but it was not found on disk!\n",s.Transport.Addr(), msg.Key)
	}

	log.Printf("[%s]: Serving file over the wire\n", s.Transport.Addr())

	size, r, err := s.store.Read(msg.Key); 
	if err != nil {
		return err
	}

	if rc, ok := r.(io.ReadCloser); ok {
		log.Println("Closing reader!")
		defer rc.Close();
	}

	peer, ok := s.peers[from];
	if !ok{
		return fmt.Errorf("Peer not found!")
	}

	// First IncomingStream then file size then stream
	peer.Send([]byte{p2p.IncomingStream})
	var fileSize int64 = size
	binary.Write(peer, binary.LittleEndian, fileSize);
	n, err := io.Copy(peer, r);
	if err != nil {
		return err
	}

	log.Printf("[%s]: Written (%d) bytes over wire to (%s)\n", s.Transport.Addr(), n, from)
	return nil
}


func (s *FileServer) handleStoreFileMessage(from string, msg StoreFileMessage) error {
	peer, ok := s.peers[from];
	if(!ok) {
		return fmt.Errorf("Peer %s not found in peer map!\n", from)
	}

		
	n, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size)); 
	if err != nil {
		return err
	}
	log.Printf("Written (%d)bytes to disk\n", n)

	peer.CloseStream()

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
	gob.Register(StoreFileMessage{})
	gob.Register(GetFileMessage{})
}
