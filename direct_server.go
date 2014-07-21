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
	cacher *Cacher
}

var servingPattern = regexp.MustCompile("(/$)|(\\.html$)|(\\.rdf$)|(\\.xml$)")

func (self *DirectServer) ShouldServe(url *url.URL) bool {
	return nil != url && servingPattern.MatchString(url.Path)
}

type DirectServerResponse struct {
	statusCode int
	head       http.Header
	body       []byte
	hit        bool
}

func (self *DirectServer) get(url *url.URL) (*DirectServerResponse, error) {
	cached := <-self.cacher.AskGet(url)
	if nil != cached {
		return &DirectServerResponse{
			statusCode: 200,
			head:       cached.Head,
			body:       cached.Body,
			hit:        true,
		}, nil
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	if resp.StatusCode == 200 {
		self.cacher.AskSet(url, &CacheEntry{Body: bytes, Head: resp.Header})
	}

	return &DirectServerResponse{
		statusCode: resp.StatusCode,
		body:       bytes,
		head:       resp.Header,
		hit:        false,
	}, nil
}

func (self *DirectServerResponse) Note() string {
	if self.hit {
		return "HIT"
	}

	return "MISS"
}

func (self *DirectServer) ServeHTTP(w http.ResponseWriter, r *http.Request, url *url.URL) {
	resp, err := self.get(url)
	if err != nil {
		LogAccessError(r.URL, err)
		http.Error(w, err.Error(), 500)
		return
	}

	h := w.Header()
	for _, n := range knownHeaderNames {
		val := resp.head.Get(n)
		if 0 < len(val) {
			h.Add(n, val)
		}
	}

	since := r.Header.Get("If-Modified-Since")
	last := resp.head.Get("Last-Modified")
	if 0 < len(last) && since == last && resp.statusCode == 200 {
		w.WriteHeader(http.StatusNotModified)
		LogAccessOK(r.URL, http.StatusNotModified, resp)
		return
	}

	w.WriteHeader(resp.statusCode)
	w.Write(resp.body)
	LogAccessOK(r.URL, resp.statusCode, resp)
}

func MakeDirectServer(cacher *Cacher) *DirectServer {
	return &DirectServer{cacher: cacher}
}
