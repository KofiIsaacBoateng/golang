package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncrypt(t *testing.T) {
	src := bytes.NewReader([]byte("foo not bar!"))
	dest := new(bytes.Buffer)
	key := newEncryptionKey();

	_, err := copyEncrypt(key, src, dest)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(dest.Bytes())
}