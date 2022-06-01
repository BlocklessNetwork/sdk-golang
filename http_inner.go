package blockless

import "syscall"

type Options struct {
	Method         string `json:"method"`
	ConnectTimeout int32  `json:"connectTimeout"`
	ReadTimeout    int32  `json:"readTimeout"`
}

//go:wasm-module blockless_drivers
//export blockless_open
func blockless_open(a string, opts string, fd *int) syscall.Errno

func HttpOpen(url string, options string) (int, error) {
	return 0, nil
}
