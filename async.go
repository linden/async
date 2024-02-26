//go:build js && wasm

package async

import (
	"errors"
	"syscall/js"
)

var (
	Promise = js.Global().Get("Promise")
	Error   = js.Global().Get("Error")
)

func NewPromise(fn func() (any, error)) js.Value {
	var h js.Func

	// create a promise executor.
	h = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer h.Release()

		res := args[0]
		rej := args[1]

		// run the executor.
		go func() {
			val, err := fn()

			if err != nil {
				// reject the promise.
				// convert the Go error to a javascript Error.
				rej.Invoke(Error.New(err.Error()))
				return
			}

			// resolve the promise.
			res.Invoke(val)
		}()

		return nil
	})

	// create a new promise.
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/Promise.
	return Promise.New(h)
}

func Await(p js.Value) (*js.Value, error) {
	valChan := make(chan js.Value)
	errChan := make(chan error)

	// handle when the promise has been fufilled.
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/then.
	h := js.FuncOf(func(this js.Value, args []js.Value) any {
		valChan <- args[0]
		return nil
	})
	defer h.Release()

	p.Call("then", h)

	// handle when the promise has been broken.
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/catch.
	h = js.FuncOf(func(this js.Value, args []js.Value) any {
		// convert the javascript Error to a Go error.
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error/toString.
		errChan <- errors.New(args[0].Get("message").String())
		return nil
	})
	defer h.Release()

	p.Call("catch", h)

	select {
	case val := <-valChan:
		return &val, nil

	case err := <-errChan:
		return nil, err
	}
}
