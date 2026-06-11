package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}


type GOBDecoder struct {}


func (decoder GOBDecoder) Decode(r io.Reader, rpc *RPC) error {
	
	return gob.NewDecoder(r).Decode(rpc); 
}


type DefaultDecoder struct {}

func (decoder DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	peekBuf := make([]byte, 1);

	if _, err := r.Read(peekBuf); err != nil {
		return err
	}

	// If r is a stream, we are not decoding what is being sent over the network.
	stream := peekBuf[0] == IncomingStream;
	if stream {
		rpc.Stream = true
		return nil
	}

	buf := make([]byte, 1028);

	n, err := r.Read(buf);
	if err != nil {
		return err
	}
	
	rpc.Payload = buf[:n]
	return nil
}