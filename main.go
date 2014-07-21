package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var _ = fmt.Printf

type MainHandler struct {
	mapper *URLMapper
	direct *DirectServer
	cacher *Cacher
}

func (self *MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/clear" {
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

		return
	}

	pair := self.mapper.MapToURLPair(r)
	if self.direct.ShouldServe(pair.Stored) {
		self.direct.ServeHTTP(w, r, pair.Stored)
	} else {
		http.Redirect(w, r, pair.Front.String(), http.StatusMovedPermanently)
	}
}

func main() {
	mapper := &URLMapper{
		Frontend:     "steps.dodgson.org",
		LivingStore:  "blog.dodgson.org.s3-website-us-east-1.amazonaws.com",
		ArchiteStore: "bn.dodgson.org.s3-website-us-east-1.amazonaws.com",
	}

	cacher := MakeCacher()
	direct := MakeDirectServer(cacher)

	s := &http.Server{
		Addr:           ":8090",
		Handler:        &MainHandler{mapper: mapper, direct: direct, cacher: cacher},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
