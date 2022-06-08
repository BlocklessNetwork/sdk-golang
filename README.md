# blockless-golang-sdk

This is blockless sandbox golang sdk, we can use it develop the app with the golang. 

The sdk is work with tinygo.

The example for http request run in blockless sandbox.

```go
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
```

How to compile?

1. Install the tinygo, the url is https://tinygo.org/.

2. use the follow command to compile.

```bash
$ tinygo build -o hello_http.wasi -target wasi .
```