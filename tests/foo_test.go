package tests

import (
	"aske-w/itu-minitwit/foo"
	"testing"
)

func TestFoo(t *testing.T) {
	result := foo.Foo(3)
	expected := 5
	if result != expected {
		t.Error(t.Name(), ": expected '", expected, "' but got '", result, "'")
	}
}

func TestFooAnother(t *testing.T) {
	result := foo.Foo(3)
	expected := 4
	if result != expected {
		t.Error(t.Name(), ": expected '", expected, "' but got '", result, "'")
	}
}
