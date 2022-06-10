package main

import (
	"fmt"

	http "github.com/txlabs/blockless-sdk-golang/http"
	"github.com/txlabs/blockless-sdk-golang/jsonparser"
)

func main() {
	var handle *http.HttpHandle
	var err error
	opts := http.NewDefaultHttpOptions("POST")
	opts.Body = `{}`
	if handle, err = http.HttpRequest("https://demo.bls.dev/tokens", opts); err != nil {
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
		var item []byte
		var jsonType jsonparser.ValueType
		if bs, err = handle.ReadBodyAll(); err != nil {
			panic(err)
		}
		if item, jsonType, _, err = jsonparser.Get(bs, "tokens"); err != nil {
			panic(err)
		}
		if jsonType == jsonparser.NotExist {
			fmt.Println("tokens is not exits")
			return
		}
		if jsonType == jsonparser.Array {
			jsonparser.ArrayEach(item, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if dataType == jsonparser.String {
					fmt.Println(string(value))
				}
			})
		}
	}
}
