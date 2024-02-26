//go:build js && wasm

package async

import (
	"errors"
	"testing"
)

func TestAsync(t *testing.T) {
	ev := "Hello World"

	p := NewPromise(func() (any, error) {
		return ev, nil
	})

	val, err := Await(p)
	if err != nil {
		t.Fatal(err)
	}

	if val.String() != ev {
		t.Fatalf("expected %s but got %s", val, ev)
	}
}

func TestCatch(t *testing.T) {
	exp := errors.New("uh-oh")

	p := NewPromise(func() (any, error) {
		return nil, exp
	})

	_, err := Await(p)
	if errors.Is(exp, err) {
		t.Fatalf("expected %v got %v", exp, err)
	}
}
