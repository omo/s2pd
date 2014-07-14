package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var servingPattern = regexp.MustCompile("(/$)|(\\.html$)|(\\.rdf$)|(\\.xml$)")

func ShouldServeDirectly(url *url.URL) bool {
	return servingPattern.MatchString(url.Path)
}

// --- URLMappter ---

type URLMapper struct {
	LivingSite  string
	ArchiveSite string
}

func (self *URLMapper) mapToLivingSite(req *http.Request) *url.URL {
	urlToReturn := *(req.URL)
	urlToReturn.Host = self.LivingSite
	urlToReturn.Scheme = "http"
	return &urlToReturn
}

func (self *URLMapper) mapToLivingSiteAtom() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   self.LivingSite,
		Path:   "/atom.xml",
	}
}

var datePattern = regexp.MustCompile("([[:digit:]]{4})([[:digit:]]{2})([[:digit:]]{2})")

func (self *URLMapper) MapToURL(req *http.Request) *url.URL {
	values := req.URL.Query()
	date, hasDate := values["date"]
	if !hasDate {
		if req.URL.Path == "/index.rdf" || req.URL.Path == "/no_comments.rdf" {
			return self.mapToLivingSiteAtom()
		}

		return self.mapToLivingSite(req)
	}

	dateMatches := datePattern.FindStringSubmatch(date[0])
	if nil == dateMatches {
		return self.mapToLivingSite(req)
	}

	path := fmt.Sprintf("/%s/%s/%s/", dateMatches[1], dateMatches[2], dateMatches[3])
	return &url.URL{
		Scheme: "http",
		Host:   self.ArchiveSite,
		Path:   path,
	}
}

// --- MainHandler ---

type MainHandler struct {
	mapper *URLMapper
}

func (self *MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := self.mapper.MapToURL(r)
	http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
}

func main() {
	mapper := &URLMapper{
		LivingSite:  "blog.dodgson.org.s3-website-us-east-1.amazonaws.com",
		ArchiveSite: "bn.dodgson.org.s3-website-us-east-1.amazonaws.com",
	}

	s := &http.Server{
		Addr:           ":8090",
		Handler:        &MainHandler{mapper: mapper},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
