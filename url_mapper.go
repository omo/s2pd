package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// --- URLMappter ---

type URLMapper struct {
	LivingSite  string
	ArchiveSite string
}

func (self *URLMapper) mapWithSamePathAt(url *url.URL, host string) *url.URL {
	urlToReturn := *url
	urlToReturn.Host = host
	urlToReturn.Scheme = "http"
	return &urlToReturn
}

func (self *URLMapper) mapToLivingSite(url *url.URL) *url.URL {
	return self.mapWithSamePathAt(url, self.LivingSite)
}

func (self *URLMapper) mapToArchiveSite(url *url.URL) *url.URL {
	return self.mapWithSamePathAt(url, self.ArchiveSite)
}

func (self *URLMapper) mapToLivingSiteAtom() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   self.LivingSite,
		Path:   "/atom.xml",
	}
}

var dateQueryPattern = regexp.MustCompile("([[:digit:]]{4})([[:digit:]]{2})([[:digit:]]{2})")

//
// Note that there are no such things like:
// - RSS for backnumbers (There is nothing new coming)
// - Assets for tDiary (We support only linked pages and old RSS readers.)
//
func (self *URLMapper) MapToURL(req *http.Request) *url.URL {
	// Asset and HTML files for backnumbers
	if 0 == strings.Index(req.URL.Path, "/bn/") {
		return self.mapToArchiveSite(req.URL)
	}

	// HTML files for current blogs
	if 0 == strings.Index(req.URL.Path, "/b/") {
		return self.mapToLivingSite(req.URL)
	}

	// RSS for tDiary
	if req.URL.Path == "/index.rdf" || req.URL.Path == "/no_comments.rdf" {
		return self.mapToLivingSiteAtom()
	}

	// HTML for tDiary
	values := req.URL.Query()
	date, hasDate := values["date"]
	if hasDate {
		dateMatches := dateQueryPattern.FindStringSubmatch(date[0])
		if nil == dateMatches {
			// This isn't expected. Fallback.
			return self.mapToLivingSite(req.URL)
		}

		path := fmt.Sprintf("/%s/%s/%s/", dateMatches[1], dateMatches[2], dateMatches[3])
		return &url.URL{
			Scheme: "http",
			Host:   self.ArchiveSite,
			Path:   path,
		}
	}

	// Anything else. Probably Assets and Atom for current blogs.
	return self.mapToLivingSite(req.URL)
}
