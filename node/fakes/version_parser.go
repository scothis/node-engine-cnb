package fakes

import "sync"

type VersionParser struct {
	ParseCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Version string
			Err     error
		}
		Stub func(string) (string, error)
	}
}

func (f *VersionParser) Parse(param1 string) (string, error) {
	f.ParseCall.Lock()
	defer f.ParseCall.Unlock()
	f.ParseCall.CallCount++
	f.ParseCall.Receives.Path = param1
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub(param1)
	}
	return f.ParseCall.Returns.Version, f.ParseCall.Returns.Err
}
