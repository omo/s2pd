package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var _ = fmt.Printf

func MustParse(str string) *url.URL {
	u, e := url.Parse(str)
	if e != nil {
		log.Fatal(e)
	}

	return u
}

type MainHandler struct {
	mapper *URLMapper
	direct *DirectServer
	cacher *Cacher
}

func (self *MainHandler) serveClear(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "text/html")
		io.WriteString(w,
			"<!DOCTYPE html>\n<form method=POST><input type='submit' value='Clear Cache'></form>")
	} else if r.Method == "POST" {
		self.cacher.AskReset()
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(200)
		io.WriteString(w, "OK")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	LogAccessMisc(r.URL)
	return
}

func (self *MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/clear" {
		self.serveClear(w, r)
		return
	}

	pair := self.mapper.GetMapping(r)
	if self.direct.ShouldServe(pair.Stored) {
		self.direct.ServeHTTP(w, r, pair.Stored)
	} else {
		u := pair.GetURLToRedirect()
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		LogAccessRedirect(r.URL, u)
	}
}

var kFrontend = flag.String("frontend", "http://steps.dodgson.org/", "Frontend URL")
var kLast = flag.String("Last", "http://blog.dodgson.org.s3-website-us-east-1.amazonaws.com/", "Active Blog CDN URL")
var kArchive = flag.String("archive", "http://bn.dodgson.org.s3-website-us-east-1.amazonaws.com/", "Old Blog CDN URL")

func main() {
	flag.Parse()

	mapper := &URLMapper{
		Frontend:     MustParse(*kFrontend),
		LastStore:  MustParse(*kLast),
		ArchiveStore: MustParse(*kArchive),
	}

	cacher := MakeCacher()
	direct := MakeDirectServer(cacher)

	s := &http.Server{
		Addr:           ":8300",
		Handler:        &MainHandler{mapper: mapper, direct: direct, cacher: cacher},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
