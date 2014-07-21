package main

import (
	"log"
	"net/url"
)

type LogAccessOKEntry interface {
	Note() string
}

func LogAccessOK(url *url.URL, code int, resp LogAccessOKEntry) {
	log.Printf("A %s %d %s\n", url.String(), code, resp.Note())
}

func LogAccessError(url *url.URL, err error) {
	log.Printf("E %s %s\n", url.String(), err.Error())
}

func LogAccessRedirect(from *url.URL, to *url.URL) {
	log.Printf("R %s %s\n", from.String(), to.String())
}

func LogAccessMisc(url *url.URL) {
	log.Printf("M %s\n", url.String())
}
