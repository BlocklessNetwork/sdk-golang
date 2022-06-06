# blockless-golang-sdk

This is blockless sandbox golang sdk, we can use it develop the app with the golang. 

The sdk is work with tinygo.

The example for http request run in blockless sandbox.

```go
import "blockless/http"

func main() {
    if handle, err := http.HttpOpen("https://www.163.com", http.NewDefaultHttpOptions()); err != nil {
        panic(err)
    }

}
```

