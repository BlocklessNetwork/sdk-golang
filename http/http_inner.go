package http

import (
	"reflect"
	"syscall"
	"unsafe"
)

type innerHandle uint32

type Options struct {
	//http method, GET POST etc.
	Method string `json:"method"`
	//connect timeout
	ConnectTimeout int32 `json:"connectTimeout"`
	//read timeout
	ReadTimeout int32 `json:"readTimeout"`
}

type HttpHandle struct {
	// inner handle for http from http driver.
	inner innerHandle
}

//go:wasm-module blockless_http
//export http_open
func http_open(a string, opts string, fd *innerHandle) syscall.Errno

//go:wasm-module blockless_http
//export http_close
func http_close(fd uint32) syscall.Errno

//go:wasm-module blockless_http
//export http_read_header
func http_read_header(fd uint32, header string, buf uintptr, bufLen uint32, num *uint32) syscall.Errno

//go:wasm-module blockless_http
//export http_read_body
func http_read_body(fd uint32, buf uintptr, bufLen uint32, num *uint32) syscall.Errno

//open a url with the options
//if success return the http handle
func HttpOpen(url string, options string) (*HttpHandle, error) {
	var handle innerHandle
	err := http_open(url, options, &handle)
	if err != 0 {
		return nil, Error(err)
	}
	return &HttpHandle{inner: innerHandle(handle)}, nil
}

//http handle close
func (h *HttpHandle) Close() error {
	err := http_close(uint32(h.inner))
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (h *HttpHandle) GetHeader(header string) (string, error) {
	buf := make([]byte, 1024*10)
	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&buf))
	var num uint32 = 0
	err := http_read_header(uint32(h.inner), header, sliceHeader.Data, uint32(len(buf)), &num)
	if err != 0 {
		return "", Error(err)
	}
	return string(buf[:num]), nil
}

func (h *HttpHandle) ReadBody(buf []byte) (string, error) {
	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&buf))
	var num uint32 = 0
	err := http_read_body(uint32(h.inner), sliceHeader.Data, uint32(len(buf)), &num)
	if err != 0 {
		return "", Error(err)
	}
	return string(buf[:num]), nil
}
