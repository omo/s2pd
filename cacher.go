package main

import (
	"fmt"
	"net/http"
	"net/url"
)

var _ = fmt.Printf

type CacherRequest interface {
}

type CacherGet struct {
	url          *url.URL
	responseChan chan *CacheEntry
}

type CacherSet struct {
	url   *url.URL
	entry *CacheEntry
}

type CacherReset struct {
}

type CacheEntry struct {
	Body []byte
	Head http.Header
}

type Cacher struct {
	requestChan chan CacherRequest
	cachedData  map[string]*CacheEntry
}

func (self *Cacher) AskSet(url *url.URL, entry *CacheEntry) {
	self.requestChan <- &CacherSet{url: url, entry: entry}
}

func (self *Cacher) AskReset() {
	self.requestChan <- &CacherReset{}
}

func (self *Cacher) AskGet(url *url.URL) chan *CacheEntry {
	resChan := make(chan *CacheEntry)
	self.requestChan <- &CacherGet{url: url, responseChan: resChan}
	return resChan
}

func (self *Cacher) serve() {
	for r := range self.requestChan {
		switch rr := r.(type) {
		case *CacherGet:
			rr.responseChan <- self.cachedData[rr.url.String()]
		case *CacherSet:
			self.cachedData[rr.url.String()] = rr.entry
		case *CacherReset:
			self.cachedData = make(map[string]*CacheEntry)
		}
	}
}

func MakeCacher() *Cacher {
	req_chan := make(chan CacherRequest)
	self := &Cacher{requestChan: req_chan, cachedData: make(map[string]*CacheEntry)}
	go self.serve()

	return self
}
