package node

import (
	"sort"

	"github.com/cloudfoundry/node-engine-cnb/packit"
)

type PlanEntryResolver struct{}

func NewPlanEntryResolver() PlanEntryResolver {
	return PlanEntryResolver{}
}

func (r PlanEntryResolver) Resolve(entries []packit.BuildpackPlanEntry) packit.BuildpackPlanEntry {
	var priorities = map[string]int{
		"buildpack.yml": 3,
		"package.json":  2,
		".nvmrc":        1,
		"":              -1,
	}

	sort.Slice(entries, func(i, j int) bool {
		leftSource := entries[i].Metadata["VersionSource"]
		left, _ := leftSource.(string)

		rightSource := entries[j].Metadata["VersionSource"]
		right, _ := rightSource.(string)

		return priorities[left] > priorities[right]
	})

	return entries[0]
}
