package packit

type Option func(config Config) Config

//go:generate faux --interface ExitHandler --output fakes/exit_handler.go
type ExitHandler interface {
	Error(error)
}

func WithExitHandler(exitHandler ExitHandler) Option {
	return func(config Config) Config {
		config.exitHandler = exitHandler
		return config
	}
}

func WithArgs(args []string) Option {
	return func(config Config) Config {
		config.args = args
		return config
	}
}
