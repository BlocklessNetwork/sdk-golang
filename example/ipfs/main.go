package main

import (
	"fmt"

	ipfs "github.com/txlabs/blockless-sdk-golang/ipfs"
)

func main() {
	var err error
	var num uint32
	var files []*ipfs.File
	var fstat *ipfs.FileStat

	err = ipfs.IpfsCreateDir("/test/test", true)
	if err != nil {
		panic(err)
	}
	if fstat, err = ipfs.IpfsFileStat("/test/test"); err != nil {
		panic(err)
	}
	fmt.Println("test stat:", fstat)
	err = ipfs.IpfsFileRemove("/test/test", true, true)
	if err != nil {
		panic(err)
	}

	files, err = ipfs.IpfsFileList("/")
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
	var fn = "/test_golang.txt"
	var fpath = ipfs.NewFileWriteOptions(fn)
	if num, err = ipfs.IpfsFileWrite(fpath, []byte("12345678910ABC")); err != nil {
		panic(err)
	}
	fmt.Println("file write", fn, num)
	var buf = make([]byte, 1024)
	if num, err = ipfs.IpfsFileRead(fn, 0, buf); err != nil {
		panic(err)
	}
	fmt.Println("file read", fn, num, "content:", string(buf))
	ipfs.IpfsFileRemove(fn, true, true)
}
