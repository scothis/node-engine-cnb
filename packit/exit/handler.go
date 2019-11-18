package exit

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var Fail = errors.New("failed")

type Handler struct {
	stderr   io.Writer
	exitFunc func(int)
}

func NewHandler(options ...Option) Handler {
	handler := Handler{
		stderr:   os.Stderr,
		exitFunc: os.Exit,
	}

	for _, option := range options {
		handler = option(handler)
	}

	return handler
}

func (h Handler) Error(err error) {
	fmt.Fprintln(h.stderr, err)

	var code int
	switch err {
	case Fail:
		code = 100
	case nil:
		code = 0
	default:
		code = 1
	}

	h.exitFunc(code)
}
