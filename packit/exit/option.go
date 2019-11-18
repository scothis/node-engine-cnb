package exit

import "io"

type Option func(handler Handler) Handler

func WithStderr(stderr io.Writer) Option {
	return func(handler Handler) Handler {
		handler.stderr = stderr
		return handler
	}
}

func WithExitFunc(e func(int)) Option {
	return func(handler Handler) Handler {
		handler.exitFunc = e
		return handler
	}
}
