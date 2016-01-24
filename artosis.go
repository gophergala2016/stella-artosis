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
	"runtime"
	"strings"
	"sync"
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

func Scan(conf Config) []*File {
	ret := []*File{}
	var wg sync.WaitGroup
	fn, c := genScan(conf)

	o := make(chan *File)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		h := hashes[conf.Hash]()
		go func() {
			for f := range c {
				stella(f, h, o)
			}
		}()
	}
	for r, _ := range conf.Include {
		fmt.Printf("scanning: %s\n", r)
		wg.Add(1)
		go func() {
			filepath.Walk(r, fn)
			wg.Done()
		}()
	}
	go func() {
		for f := range o {
			fmt.Println(f.Stat.Name())
			ret = append(ret, f)
		}
	}()
	wg.Wait()
	close(c)
	return ret
}

func skip(s string, ign map[string]bool) bool {
	for i := range ign {
		if strings.HasPrefix(s, i) {
			return true
		}
	}
	return false
}

func stella(file *File, h hash.Hash, out chan *File) error {
	f, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	fmt.Printf("%x  %s\n", h.Sum(nil), file.Stat.Name())
	h.Reset()
	return nil
}

func genScan(conf Config) (filepath.WalkFunc, chan *File) {
	c := make(chan *File)
	ret := func(path string, info os.FileInfo, err error) error {
		if skip(path, conf.Ignored) {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		f := &File{
			Stat: info,
			Path: path,
		}
		c <- f
		return nil
	}
	return ret, c
}
