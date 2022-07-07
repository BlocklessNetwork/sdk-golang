package ipfs

import (
	"fmt"
	"io"
	"reflect"
	"syscall"
	"unsafe"
)

type innerHandle uint32
type StatusCode int32

//go:wasm-module blockless_ipfs
//export ipfs_command
func ipfs_command(opts string, fd *innerHandle, code *StatusCode) syscall.Errno

//go:wasm-module blockless_ipfs
//export ipfs_close
func ipfs_close(fd innerHandle) syscall.Errno

//go:wasm-module blockless_ipfs
//export ipfs_read
func ipfs_read(fd innerHandle, buf uintptr, bufLen uint32, num *uint32) syscall.Errno

//go:wasm-module blockless_ipfs
//export ipfs_write
func ipfs_write(fd innerHandle, buf uintptr, bufLen uint32, num *uint32) syscall.Errno

type commandRs struct {
	code   StatusCode
	handle innerHandle
}

type CommanResult struct {
	statusCode uint32
	respBody   []byte
}

func NewCommanResult(statusCode uint32, respBody []byte) *CommanResult {
	return &CommanResult{
		statusCode: statusCode,
		respBody:   respBody,
	}
}

func ipfsCommandResult(opts *IpfsOptions) (*CommanResult, error) {
	var err error
	var bs []byte
	var crs *commandRs
	if crs, err = ipfsCommand(opts); err != nil {
		return nil, err
	}
	defer ipfs_close(crs.handle)
	if bs, err = readBodyAll(crs.handle); err != nil {
		return nil, err
	}
	return NewCommanResult(uint32(crs.code), bs), nil
}

func ipfsCommand(opts *IpfsOptions) (*commandRs, error) {
	var handle innerHandle
	var code StatusCode
	var optsJson = opts.ToJsonString()
	fmt.Println(optsJson)
	var errno syscall.Errno
	if errno = ipfs_command(optsJson, &handle, &code); errno != 0 {
		return nil, Error(errno)
	}
	return &commandRs{
		code:   code,
		handle: handle,
	}, nil
}

func readBody(h innerHandle, buf []byte) (uint32, error) {
	if cap(buf) == 0 {
		return 0, BUFFER_TOO_SMALL
	}
	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&buf))
	var num uint32 = 0
	err := ipfs_read(h, sliceHeader.Data, uint32(sliceHeader.Cap), &num)
	if err != 0 {
		return num, Error(err)
	}
	return num, nil
}

func readBodyAll(h innerHandle) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := readBody(h, b[len(b):cap(b)])
		//End
		if n == 0 && err == nil {
			return b, nil
		}
		b = b[:len(b)+int(n)]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}
