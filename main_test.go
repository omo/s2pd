package main

import (
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestHello(t *testing.T) {
}

func makeTestRequest(host string, path string) *http.Request {
	u, e := url.Parse(path)
	if e != nil {
		log.Fatal(e)
	}

	return &http.Request{
		Method: "GET",
		Proto:  "http",
		URL:    u,
		Host:   host,
	}
}

func TestMapFromTDiaryRequest(t *testing.T) {
	mapper := &URLMapper{LivingSite: "living.aws", ArchiveSite: "archive.aws"}
	req := makeTestRequest("steps.dodgson.org", "/?date=20091231")
	Expect(mapper.MapToURL(req).String(), "http://archive.aws/2009/12/31/", t)
}

func TestMapFromTDiaryAtomRequest(t *testing.T) {
	mapper := &URLMapper{LivingSite: "living.aws", ArchiveSite: "archive.aws"}
	req := makeTestRequest("steps.dodgson.org", "/index.rdf")
	Expect(mapper.MapToURL(req).String(), "http://living.aws/atom.xml", t)
}

func TestMapFromOctopressRequest(t *testing.T) {
	mapper := &URLMapper{LivingSite: "living.aws", ArchiveSite: "archive.aws"}
	req := makeTestRequest("steps.dodgson.org", "/b/2014/07/05/life-of-touch/")
	Expect(mapper.MapToURL(req).String(), "http://living.aws/b/2014/07/05/life-of-touch/", t)
}

func MustParse(str string) *url.URL {
	u, e := url.Parse(str)
	if e != nil {
		log.Fatal(e)
	}

	return u
}

func TestShouldServeDirectly(t *testing.T) {
	ExpectTrue(ShouldServeDirectly(MustParse("http://example.com/")), "/", t)
	ExpectTrue(ShouldServeDirectly(MustParse("http://example.com/foo/bar/")), "/", t)
	ExpectTrue(ShouldServeDirectly(MustParse("http://example.com/atom.xml")), "atom.xml", t)
	ExpectTrue(ShouldServeDirectly(MustParse("http://example.com/index.rdf")), "index.rdf", t)
	ExpectTrue(ShouldServeDirectly(MustParse("http://example.com/index.html")), "index.html", t)
	ExpectTrue(!ShouldServeDirectly(MustParse("http://example.com/hoge.css")), "css", t)
	ExpectTrue(!ShouldServeDirectly(MustParse("http://example.com/fuga.jpeg")), "css", t)
}
