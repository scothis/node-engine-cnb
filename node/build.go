package node

import (
	"github.com/cloudfoundry/node-engine-cnb/packit"
)

type Dependency struct{}

//go:generate faux --interface EntryResolver --output fakes/entry_resolver.go
type EntryResolver interface {
	Resolve([]packit.BuildpackPlanEntry) packit.BuildpackPlanEntry
}

//go:generate faux --interface DependencyManager --output fakes/dependency_manager.go
type DependencyManager interface {
	Resolve(entry packit.BuildpackPlanEntry, stack string) (Dependency, error)
	Install(Dependency, packit.Layer) error
}

func Build(entries EntryResolver, dependencies DependencyManager) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		entry := entries.Resolve(context.Plan.Entries)

		dependency, err := dependencies.Resolve(entry, context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		nodeLayer, err := context.Layers.Get(Node, packit.LaunchLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		err = dependencies.Install(dependency, nodeLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		return packit.BuildResult{
			Plan:   context.Plan,
			Layers: []packit.Layer{nodeLayer},
		}, nil
	}
}
