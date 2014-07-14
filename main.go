package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var _ = fmt.Printf

type MainHandler struct {
	mapper *URLMapper
	direct *DirectServer
}

func (self *MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := self.mapper.MapToURL(r)
	if self.direct.ShouldServe(url) {
		// FIXME: I might want to direct old sites to the actual URL.
		self.direct.ServeHTTP(w, r, url)
	} else {
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
	}
}

func main() {
	mapper := &URLMapper{
		LivingSite:  "blog.dodgson.org.s3-website-us-east-1.amazonaws.com",
		ArchiveSite: "bn.dodgson.org.s3-website-us-east-1.amazonaws.com",
	}

	s := &http.Server{
		Addr:           ":8090",
		Handler:        &MainHandler{mapper: mapper, direct: MakeDirectServer()},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
