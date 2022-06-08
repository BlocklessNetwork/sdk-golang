package http

import (
	"fmt"
	"io"
	"reflect"
	"syscall"
	"unsafe"
)

type innerHandle uint32
type StatusCode int32

type HttpOptions struct {
	//http method, GET POST etc.
	Method string `json:"method"`
	//connect timeout, unit is second.
	ConnectTimeout int32 `json:"connectTimeout"`
	//read timeout, unit is second.
	ReadTimeout int32 `json:"readTimeout"`
	//request Body
	Body string `json:"body"`
}

func NewDefaultHttpOptions(method string) HttpOptions {
	return HttpOptions{
		Method:         method,
		ConnectTimeout: 30,
		ReadTimeout:    30,
	}
}

type HttpHandle struct {
	// inner handle for http from http driver.
	inner innerHandle
	code  StatusCode
}

//go:wasm-module blockless_http
//export http_req
func http_req(a string, opts string, fd *innerHandle, code *StatusCode) syscall.Errno

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
func HttpRequest(url string, options HttpOptions) (*HttpHandle, error) {
	var handle innerHandle
	//format the options to json format, the json string will parse the "".
	//TODO.
	var opts = fmt.Sprintf(`{"method":"%s", "connectTimeout":%d, "readTimeout":%d, "body":"%s"}`,
		options.Method,
		options.ConnectTimeout,
		options.ReadTimeout,
		options.Body,
	)
	var code StatusCode
	err := http_req(url, opts, &handle, &code)
	if err != 0 {
		return nil, Error(err)
	}
	return &HttpHandle{inner: innerHandle(handle), code: code}, nil
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
	err := http_read_header(uint32(h.inner), header, sliceHeader.Data, uint32(cap(buf)), &num)
	if err != 0 {
		return "", Error(err)
	}
	return string(buf[:num]), nil
}

func (h *HttpHandle) StatusCode() StatusCode {
	return h.code
}

func (h *HttpHandle) ReadBody(buf []byte) (uint32, error) {

	if cap(buf) == 0 {
		return 0, BUFFER_TOO_SMALL
	}
	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&buf))
	var num uint32 = 0
	err := http_read_body(uint32(h.inner), sliceHeader.Data, uint32(sliceHeader.Cap), &num)
	if err != 0 {
		return num, Error(err)
	}
	return num, nil
}

func (h *HttpHandle) ReadBodyAll() ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := h.ReadBody(b[len(b):cap(b)])
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
