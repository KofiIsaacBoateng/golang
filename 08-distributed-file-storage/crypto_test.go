package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncrypt(t *testing.T) {
	payload := []byte("foo not bar!")
	src := bytes.NewReader(payload)
	dest := new(bytes.Buffer)
	key := newEncryptionKey();

	if _, err := copyEncrypt(key, src, dest); err != nil {
		t.Error(err)
	}

	fmt.Println("Encrypt:", dest.String())

	out := new(bytes.Buffer);
	n, err := copyDecrypt(key, dest, out); 
	if err != nil {
		t.Error(err)
	}

	if n != len(payload) + 16 {
		t.Fail();
	}

	fmt.Println("Decrypt:", out.String())

	if out.String() != string(payload) {
		t.Error("Decryption failed!")
	}

}