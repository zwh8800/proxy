package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type h struct{}

func (h) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxyTo := r.RequestURI[1:] // remove '/'

	resp, err := client.Get(proxyTo)
	if err != nil {
		return
	}
	if resp.StatusCode == 302 {
		loc := resp.Header.Get("Location")
		w.Header().Add("Location", "/"+loc)
		w.WriteHeader(302)
		return
	}

	_, _ = io.Copy(w, resp.Body)
}

func main() {
	addr := ":8898"
	if len(os.Args) >= 2 {
		addr = os.Args[1]
	}
	log.Fatal(http.ListenAndServe(addr, h{}))
}
