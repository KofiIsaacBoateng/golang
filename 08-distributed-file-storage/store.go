package main

import (
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
	RootDir string
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


func (s *Store) OpenFileForWrite(key string, r io.Reader) (*os.File, error) {
	pathkey := s.PathTransformerFunc(key);

	fileDir := s.RootDir + "/" + pathkey.Filepath
	filepath := fileDir + "/" + pathkey.Filename

	// create directory
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return nil, err;
	}

	// create file
	return os.Create(filepath);
}


func (s *Store) WriteDecrypt(EncKey []byte, key string, r io.Reader) (int64, error) {
	f, err := s.OpenFileForWrite(key, r);
	if err != nil {
		return 0, err
	}

	n, err := copyDecrypt(EncKey, r, f);
	if err != nil {
		return 0, err
	}

	return int64(n), nil
}

func (s *Store) WriteStream(key string, r io.Reader) (int64, error) {
	f, err := s.OpenFileForWrite(key, r);
	if(err != nil) {
		return 0, err
	}

	defer f.Close()

	return io.Copy(f, r);
}

func (s *Store) Read(key string) (int64, io.Reader, error) {
	return s.ReadStream(key);
	
}


func (s *Store) ReadStream(key string) (int64, io.ReadCloser, error) {
	pathkey := s.PathTransformerFunc(key)

	filepath := s.RootDir + "/" + pathkey.Filepath + "/" + pathkey.Filename;
	f, err := os.Open(filepath)
	if err != nil {
		return 0, nil, err
	}

	fs, err := os.Stat(filepath);
	if err != nil {
		return 0, nil, err
	}

	return fs.Size(), f, nil
}


func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformerFunc(key)

	filepath := s.RootDir + "/" + pathkey.Filepath + "/" + pathkey.Filename;
	_, err := os.Stat(filepath);

	return !errors.Is(err, os.ErrNotExist)
}


func (s *Store) Delete(key string) error {
	pathkey := s.PathTransformerFunc(key);

	if err := os.RemoveAll(pathkey.RootDir()); err != nil {
		return err
	}

	fmt.Printf("deleted [%s] from disk", pathkey.Filename)

	return os.RemoveAll(s.RootDir + "/" + pathkey.RootDir())
}