package http

import "syscall"

type innerHandle int

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

//go:wasm-module blockless_drivers
//export blockless_open
func blockless_open(a string, opts string, fd *int) syscall.Errno

//open a url with the options
//if success return the http Object
func HttpOpen(url string, options string) (*HttpHandle, error) {
	var handle int
	err := blockless_open(url, options, &handle)
	if err != 0 {
		return nil, Error(err)
	}
	return &HttpHandle{inner: innerHandle(handle)}, nil
}
