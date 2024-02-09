package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func handleRequest(client *http.Client, resp http.ResponseWriter, req *http.Request) {
	req.RequestURI = ""
	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}
	if req.URL.Host == "" {
		req.URL.Host = req.Host
	}

	log.Printf("Forwarding request to: %s", req.URL.String())

	serverResp, err := client.Do(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer serverResp.Body.Close()

	log.Printf("Received response with status code: %d", serverResp.StatusCode)

	copyHeader(resp.Header(), serverResp.Header)
	resp.WriteHeader(serverResp.StatusCode)
	io.Copy(resp, serverResp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("tcp", addr)
			},
		},
	}

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		handleRequest(client, resp, req)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
