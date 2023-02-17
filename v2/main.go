package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/fvbock/endless"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Private  bool   `json:"private"`
}

//go:embed config.json
var config []byte
var C Config

func init() {
	err := json.Unmarshal(config, &C)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   C.Endpoint,
	})

	// create a handler for the reverse proxy
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// update the headers to allow for SSL redirection
		r.URL.Host = C.Endpoint
		r.URL.Scheme = "https"
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = C.Endpoint
		// call the reverse proxy
		proxy.ServeHTTP(w, r)
	})

	// start the server on port 8080
	endless.ListenAndServe(":8080", nil)
}
