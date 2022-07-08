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

	//create dir
	err = ipfs.IpfsCreateDir("/test/test", true)
	if err != nil {
		panic(err)
	}
	//file stat example
	if fstat, err = ipfs.IpfsFileStat("/test/test"); err != nil {
		panic(err)
	}
	fmt.Println("test stat:", fstat)
	//file remove example
	err = ipfs.IpfsFileRemove("/test/test", true, true)
	if err != nil {
		panic(err)
	}

	//file list example
	files, err = ipfs.IpfsFileList("/")
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
	var fn = "/test_golang.txt"
	var fpath = ipfs.NewFileWriteOptions(fn)
	//file write example
	if num, err = ipfs.IpfsFileWrite(fpath, []byte("12345678910ABC")); err != nil {
		panic(err)
	}
	fmt.Println("file write", fn, num)
	var buf = make([]byte, 1024)
	//file read example
	if num, err = ipfs.IpfsFileRead(fn, 0, buf); err != nil {
		panic(err)
	}
	fmt.Println("file read", fn, num, "content:", string(buf))
	ipfs.IpfsFileRemove(fn, true, true)
}
