# blockless-golang-sdk
![](blockless.png)

This is blockless sandbox golang sdk, we can use it develop the app with the golang. 

The sdk is work with tinygo, the built-in support json , http  and others.

The example for http request run in blockless sandbox.

```go
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
```

### How to compile

1. Install the tinygo, the url is https://tinygo.org/.

2. use the follow command to compile.

```bash
$ tinygo build -o hello_http.wasi -target wasi .
```
