package main

import (
	"fmt"

	http "github.com/txlabs/blockless-sdk-golang/http"
)

func main() {
	var handle *http.HttpHandle
	var err error
	if handle, err = http.HttpRequest("https://blockless-website.vercel.app/", http.NewDefaultHttpOptions("GET")); err != nil {
		panic(err)
	}
	defer handle.Close()
	if handle.StatusCode() == 200 {
		var head string

		if head, err = handle.GetHeader("Content-Type"); err != nil {
			panic(err)
		}
		fmt.Println(head)
		var bs []byte
		if bs, err = handle.ReadBodyAll(); err != nil {
			panic(err)
		}
		fmt.Println(string(bs))
	}

}
