package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	proxyUrl, _ := url.Parse("http://localhost:8080")

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	resp, err := client.Get("http://www.google.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
