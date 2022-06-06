package main

import http "github.com/txlabs/blockless-sdk-golang/http"

func main() {
	if _, err := http.HttpRequest("https://www.163.com", http.NewDefaultHttpOptions("GET")); err != nil {
		panic(err)
	}

}
