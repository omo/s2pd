package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func MakeEntry(text string) *CacheEntry {
	buf := new(bytes.Buffer)
	buf.WriteString(text)
	return &CacheEntry{
		Body: buf.Bytes(),
		Head: http.Header(map[string][]string{}),
	}
}

func TestHelloCacher(t *testing.T) {
	bodyText := "Hello"
	target := MakeCacher()
	to_cache := MakeEntry(bodyText)
	u := &url.URL{Path: "/hello"}
	target.AskSet(u, to_cache)
	Expect(string((<-target.AskGet(u)).Body), bodyText, t)

	target.AskReet()
	ExpectTrue(nil == (<-target.AskGet(u)), "AskReset", t)
}

func TestCacherGetFail(t *testing.T) {
	target := MakeCacher()
	to_cache := MakeEntry("hello")
	target.AskSet(&url.URL{Path: "/hello"}, to_cache)
	to_get := target.AskGet(&url.URL{Path: "/bye"})
	got := <-to_get

	ExpectTrue(got == nil, "not found", t)
}
