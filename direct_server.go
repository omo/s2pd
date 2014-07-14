package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var _ = fmt.Printf

var knownHeaderNames = []string{
	"Content-Type",
	"Content-Length",
	"ETag",
	"Last-Modified",
}

type DirectServer struct {
}

var servingPattern = regexp.MustCompile("(/$)|(\\.html$)|(\\.rdf$)|(\\.xml$)")

func (self *DirectServer) ShouldServe(url *url.URL) bool {
	return servingPattern.MatchString(url.Path)
}

func (self *DirectServer) ServeHTTP(w http.ResponseWriter, r *http.Request, url *url.URL) {
	resp, err := http.Get(url.String())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(resp.StatusCode)
	h := w.Header()
	for _, n := range knownHeaderNames {
		val := resp.Header.Get(n)
		if 0 < len(val) {
			h.Add(n, val)
		}
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(bytes)
}

func MakeDirectServer() *DirectServer {
	return &DirectServer{}
}
