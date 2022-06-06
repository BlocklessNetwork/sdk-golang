import "blockless-sdk-golang"

func main() {
	if handle, err := http.HttpOpen("https://www.163.com", http.NewDefaultHttpOptions()); err != nil {
		panic(err)
	}

}