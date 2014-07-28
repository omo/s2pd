package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type URLMapping struct {
	Front  *url.URL
	Stored *url.URL
}

func (self *URLMapping) GetURLToRedirect() *url.URL {
	if nil != self.Stored {
		return self.Stored
	}

	return self.Front
}

// --- URLMappter ---

type URLMapper struct {
	Frontend     *url.URL
	LivingStore  *url.URL
	ArchiveStore *url.URL
}

func (self *URLMapper) mapWithSamePathAt(url *url.URL, host *url.URL) *url.URL {
	urlToReturn := *url
	urlToReturn.Host = host.Host
	urlToReturn.Scheme = host.Scheme
	return &urlToReturn
}

func (self *URLMapper) mapToLivingStore(url *url.URL) *url.URL {
	return self.mapWithSamePathAt(url, self.LivingStore)
}

func (self *URLMapper) mapToArchiveStore(url *url.URL) *url.URL {
	return self.mapWithSamePathAt(url, self.ArchiveStore)
}

func (self *URLMapper) mapToFrontend(url *url.URL) *url.URL {
	return self.mapWithSamePathAt(url, self.Frontend)
}

func (self *URLMapper) mapToLivingStoreAtom() *url.URL {
	return &url.URL{
		Scheme: self.LivingStore.Scheme,
		Host:   self.LivingStore.Host,
		Path:   "/atom.xml",
	}
}

// TODO(morrita): Should be configurable (theoretically.)
var frontHostWhitelist = map[string]bool{
	"s2p.flakiness.es": true,
	"localhost":        true,
}

func (self *URLMapper) GetMapping(req *http.Request) URLMapping {
	mapping := self.getMappingFromPath(req)
	if mapping.Front.Host != req.Host && !frontHostWhitelist[req.Host] {
		return URLMapping{Front: mapping.Front, Stored: nil}
	}

	return mapping
}

var dateQueryPattern = regexp.MustCompile("([[:digit:]]{4})([[:digit:]]{2})([[:digit:]]{2})")

//
// Note that there are no such things like:
// - RSS for backnumbers (There is nothing new coming)
// - Assets for tDiary (We support only linked pages and old RSS readers.)
//
func (self *URLMapper) getMappingFromPath(req *http.Request) URLMapping {
	front, can_be_stored := self.GetFront(req)
	if !can_be_stored {
		return URLMapping{Front: front, Stored: nil}
	}

	stored := self.GetStored(front)
	return URLMapping{Front: front, Stored: stored}
}

func (self *URLMapper) GetFront(req *http.Request) (*url.URL, bool) {
	values := req.URL.Query()
	date, hasDate := values["date"]
	if !hasDate {
		return self.mapToFrontend(req.URL), true
	}

	dateMatches := dateQueryPattern.FindStringSubmatch(date[0])
	if nil == dateMatches {
		return self.mapToFrontend(req.URL), true
	}

	// HTML for tDiary
	path := fmt.Sprintf("/bn/%s/%s/%s/", dateMatches[1], dateMatches[2], dateMatches[3])
	return &url.URL{
		Scheme: self.Frontend.Scheme,
		Host:   self.Frontend.Host,
		Path:   path,
	}, false
}

func (self *URLMapper) GetStored(url *url.URL) *url.URL {
	if 0 == strings.Index(url.Path, "/bn/") {
		return self.mapToArchiveStore(url)
	}

	if 0 == strings.Index(url.Path, "/b/") {
		return self.mapToLivingStore(url)
	}

	// RSS for tDiary
	if url.Path == "/index.rdf" || url.Path == "/no_comments.rdf" {
		return self.mapToLivingStoreAtom()
	}

	// Anything else. Probably They are:
	// - Assets and Atom for current blogs or
	// - Some non-article pages.
	return self.mapToLivingStore(url)
}
