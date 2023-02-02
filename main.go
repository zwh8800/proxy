package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	prefix string
	addr   string
)

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type h struct{}

func (h) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxyTo := strings.Replace(r.RequestURI, prefix, "", 1)

	resp, err := client.Get(proxyTo)
	if err != nil {
		return
	}
	if resp.StatusCode == 302 {
		loc := resp.Header.Get("Location")
		w.Header().Add("Location", prefix+loc)
		w.WriteHeader(302)
		return
	}

	_, _ = io.Copy(w, resp.Body)
}

func main() {
	flag.StringVar(&prefix, "prefix", "/hn-iptv/", "")
	flag.StringVar(&addr, "addr", ":8098", "")
	flag.Parse()
	log.Fatal(http.ListenAndServe(addr, h{}))
}
