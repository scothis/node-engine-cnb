package fakes

import (
	"sync"

	"github.com/cloudfoundry/node-engine-cnb/node"
	"github.com/cloudfoundry/node-engine-cnb/packit"
)

type DependencyManager struct {
	InstallCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Dependency node.Dependency
			Layer      packit.Layer
		}
		Returns struct {
			Error error
		}
		Stub func(node.Dependency, packit.Layer) error
	}
	ResolveCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Entry packit.BuildpackPlanEntry
			Stack string
		}
		Returns struct {
			Dependency node.Dependency
			Error      error
		}
		Stub func(packit.BuildpackPlanEntry, string) (node.Dependency, error)
	}
}

func (f *DependencyManager) Install(param1 node.Dependency, param2 packit.Layer) error {
	f.InstallCall.Lock()
	defer f.InstallCall.Unlock()
	f.InstallCall.CallCount++
	f.InstallCall.Receives.Dependency = param1
	f.InstallCall.Receives.Layer = param2
	if f.InstallCall.Stub != nil {
		return f.InstallCall.Stub(param1, param2)
	}
	return f.InstallCall.Returns.Error
}
func (f *DependencyManager) Resolve(param1 packit.BuildpackPlanEntry, param2 string) (node.Dependency, error) {
	f.ResolveCall.Lock()
	defer f.ResolveCall.Unlock()
	f.ResolveCall.CallCount++
	f.ResolveCall.Receives.Entry = param1
	f.ResolveCall.Receives.Stack = param2
	if f.ResolveCall.Stub != nil {
		return f.ResolveCall.Stub(param1, param2)
	}
	return f.ResolveCall.Returns.Dependency, f.ResolveCall.Returns.Error
}
