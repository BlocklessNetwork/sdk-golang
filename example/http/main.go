package main

import (
	"fmt"

	http "github.com/txlabs/blockless-sdk-golang/http"
)

func main() {
	var handle *http.HttpHandle
	var err error
	if handle, err = http.HttpRequest("https://www.163.com", http.NewDefaultHttpOptions("GET")); err != nil {
		panic(err)
	}
	defer handle.Close()
	var bs []byte
	if bs, err = handle.ReadAll(); err != nil {
		panic(err)
	}
	fmt.Println(string(bs))

}
