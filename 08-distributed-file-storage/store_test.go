package main

import (
	"bytes"
	"fmt"
	"io"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "somefilename";
	expectedFilename := "14113855c07f88b94807c82a99b8f0a351ffd044"
	expectedPathname := "14113/855c0/7f88b/94807/c82a9/9b8f0/a351f/fd044"
	pathkey := CASPathTransformFunc(key)

	assert.Equal(t, expectedFilename, pathkey.Filename, fmt.Sprintf("[EXPECTED]: %s [GOT] %s", expectedFilename, pathkey.Filename))
	assert.Equal(t, expectedPathname, pathkey.Filepath, fmt.Sprintf("[EXPECTED]: %s [GOT] %s", expectedPathname, pathkey.Filepath))
}


// func TestStoreDelete(t *testing.T) {
// 	storeOpts := StoreOpts{
// 		PathTransformerFunc: CASPathTransformFunc,
// 		RootDir: "test_root",
// 	}

// 	s := NewStore(storeOpts)

// 	key := "somekeyinourstore";
// 	data := []byte("This is a jpeg file that I am writing to disk.")

// 	// write file to storage
// 	if _, err := s.Write(key, bytes.NewReader(data)); err != nil {
// 		t.Error(err)
// 	}

// 	// delete
// 	if err := s.Delete(key); err != nil {
// 		t.Error(err)
// 	}
// }


func TestStore(t *testing.T) {
	storeOpts := StoreOpts{
		PathTransformerFunc: CASPathTransformFunc,
		RootDir: "test_root",
	}

	s := NewStore(storeOpts)

	key := "somekeyinourstore";
	data := []byte("This is a jpeg file that I am writing to disk.")

	// write file to storage
	if _, err := s.Write(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	// verify if file exists
	if ok := s.Has(key); !ok {
		t.Error("Write stream must have failed. File does not exist!")
	}

	// read file
	r, err := s.Read(key); 
	if err != nil {
		t.Error(err)
	}

	bytes, err := io.ReadAll(r)
	if  err != nil {
		t.Error(err)
	}
	
	assert.Equal(t, string(bytes), string(data), fmt.Sprintf("[EXPECTED]: %s [GOT] %s", string(data), string(bytes)))

	// remove
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

	// verify if file is removed
	if ok := s.Has(key); ok {
		t.Error("Delete must have failed. File still exist!")
	}
}
