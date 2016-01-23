package artosis

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

type newFunc func() hash.Hash

var hashes map[string]newFunc = map[string]newFunc{
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
}

type Config struct {
	// Include if non-empty will only scan the files/directories specified
	// if empty, defaults to /
	Include map[string]bool
	// Ignored will construct a set of files to ignore within Include
	Ignored map[string]bool
	Hash    string
}

type File struct {
	Stat os.FileInfo
	Path string
	Hash []byte
}

func Scan(c Config) []File {
	h := hashes[c.Hash]()
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		fmt.Printf("%x  %s\n", h.Sum(nil), info.Name())
		h.Reset()
		return nil
	}
	for r, _ := range c.Include {
		fmt.Printf("scanning: %s\n", r)
		filepath.Walk(r, walk)
	}
	return nil
}

func genScan() (filepath.WalkFunc, chan *File) {
	c := make(chan *File)
	ret := func(path string, info os.FileInfo, err error) error {
		f := &File{
			Stat: info,
			Path: path,
		}
		c <- f
		return nil
	}
	return ret, c
}
