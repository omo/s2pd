package main

import (
	"testing"
)

// Copied from github.com/eknkc/amber/amber_test.go
func Expect(cur, expected string, t *testing.T) {
	if cur != expected {
		t.Fatalf("Expected {%s} got {%s}.", expected, cur)
	}
}

func ExpectOK(err error, t *testing.T) {
	if err != nil {
		t.Fatal("Should be OK")
	}
}

func ExpectTrue(ok bool, subject string, t *testing.T) {
	if !ok {
		t.Fatalf("%s should be OK", subject)
	}
}
