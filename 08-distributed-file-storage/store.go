package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type PathKey struct {
	Filepath string
	Filename string
}

func (pk *PathKey) RootDir() string {
	return strings.Split(pk.Filepath, "/")[0]
}

type PathTransformerFunc func(string) PathKey

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key));
	hashStr := hex.EncodeToString(hash[:]);


	blockSize := 5
	sliceLen := len(hashStr) / blockSize;
	path := make([]string, sliceLen);
	for i := range sliceLen {
		path[i] = hashStr[i * 5: (i + 1) * 5]
	}

	return PathKey{
		Filepath: strings.Join(path, "/"),
		Filename: hashStr,
	}
}

func DefaultPathTransformerFunc(key string) PathKey {
	return PathKey{
		Filepath: "somefilepath",
		Filename: "somefilename",
	}
}

type StoreOpts struct {
	PathTransformerFunc PathTransformerFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.WriteStream(key, r)
}

func (s *Store) WriteStream(key string, r io.Reader) (int64, error) {
	pathkey := s.PathTransformerFunc(key);

	filepath := pathkey.Filepath + "/" + pathkey.Filename

	// create directory
	if err := os.MkdirAll(pathkey.Filepath, os.ModePerm); err != nil {
		return 0, err;
	}


	// create file
	f, err := os.Create(filepath);
	if(err != nil) {
		return 0, err
	}

	defer f.Close()

	n, err := io.Copy(f, r);
	if err!=nil {
		return 0, err
	}

	return n, nil
}

func (s *Store) Read (key string) (io.Reader, error) {
	f, err := s.ReadStream(key);
	if err != nil {
		return nil, err
	}

	defer f.Close();

	buf := new(bytes.Buffer);
	_, err = io.Copy(buf, f);
	if err != nil {
		return nil, err
	}

	fmt.Printf("Read (%d)bytes: [%s]\n", len(buf.Bytes()), buf.String())

	return buf, nil
}


func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	pathkey := s.PathTransformerFunc(key)

	filepath := pathkey.Filepath + "/" + pathkey.Filename;
	return os.Open(filepath)
}


func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformerFunc(key)

	filepath := pathkey.Filepath + "/" + pathkey.Filename;
	_, err := os.Stat(filepath);

	return !errors.Is(err, os.ErrNotExist)
}


func (s *Store) Delete(key string) error {
	pathkey := s.PathTransformerFunc(key);

	if err := os.RemoveAll(pathkey.RootDir()); err != nil {
		return err
	}

	fmt.Printf("deleted [%s] from disk", pathkey.Filename)

	return os.RemoveAll(pathkey.RootDir())
}