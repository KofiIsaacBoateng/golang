package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)


func newEncryptionKey() []byte {
	keybuf := make([]byte, 32);
	io.ReadFull(rand.Reader, keybuf);
	
	return keybuf
}


func copyEncrypt(key []byte, src io.Reader, dest io.Writer) (int, error) {
	// generate a cipher of the key 
	block, err := aes.NewCipher(key);
	if err != nil {
		return 0, err;
	}

	// slice of bytes from a random reader
	iv := make([]byte, block.BlockSize()) // 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err;
	}

	// preppend vi to the destination file
	// because we will need vi to decrypt the file
	nn, err := dest.Write(iv);
	if err != nil {
		return 0, err
	}


	//stream point... max-memory size 32 * 1024
	var (
		buf = make([]byte, 32 * 1024);
		stream = cipher.NewCTR(block, iv)
	)

	for {
		n, err := src.Read(buf);
		
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n]);
			nw, err := dest.Write(buf);
			if(err != nil) {
				return 0, err
			}

			nn += nw
		}

		if err == io.EOF {
			break;
		}


		if err != nil {
			return 0, err
		}

	}

	return nn, nil
}

