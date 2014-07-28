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

func makeTestMapper() *URLMapper {
	return &URLMapper{
		Frontend:     MustParse("http://front/"),
		LivingStore:  MustParse("http://living.aws/"),
		ArchiveStore: MustParse("http://archive.aws/"),
	}
}

func TestMapFromTDiaryRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/?date=20091231")
	Expect(mapper.GetMapping(req).Front.String(), "http://front/bn/2009/12/31/", t)
	ExpectTrue(mapper.GetMapping(req).Stored == nil, req.URL.String(), t)
}

func TestMapFromTDiaryAtomRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/index.rdf")
	Expect(mapper.GetMapping(req).Stored.String(), "http://living.aws/atom.xml", t)
}

func TestMapFromOctopressRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/b/2014/07/05/life-of-touch/")
	Expect(mapper.GetMapping(req).Stored.String(), "http://living.aws/b/2014/07/05/life-of-touch/", t)
}
func TestMapFromOctopressImageRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "http://front/images/line-tile.png?1383981792")
	Expect(mapper.GetMapping(req).Front.String(), "http://front/images/line-tile.png?1383981792", t)
	Expect(mapper.GetMapping(req).Stored.String(), "http://living.aws/images/line-tile.png?1383981792", t)
}

func TestMapFromBacknumberRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/bn/2011/11/04/")
	Expect(mapper.GetMapping(req).Stored.String(), "http://archive.aws/bn/2011/11/04/", t)
}

func TestMapFromBacknumberAssetRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/bn/stylesheets/style.css")
	Expect(mapper.GetMapping(req).Stored.String(), "http://archive.aws/bn/stylesheets/style.css", t)
}

func TestShouldServeDirectly(t *testing.T) {
	s := &DirectServer{}
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/")), "/", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/foo/bar/")), "/", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/atom.xml")), "atom.xml", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/index.rdf")), "index.rdf", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/index.html")), "index.html", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/hoge.css")), "css", t)
	ExpectTrue(!s.ShouldServe(nil), "nil", t)
	ExpectTrue(!s.ShouldServe(MustParse("http://example.com/fuga.jpeg")), "jpeg", t)
}
