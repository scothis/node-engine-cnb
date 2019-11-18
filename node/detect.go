package node

import (
	"path/filepath"

	"github.com/cloudfoundry/node-engine-cnb/packit"
)

//go:generate faux --interface VersionParser --output fakes/version_parser.go
type VersionParser interface {
	Parse(path string) (version string, err error)
}

type BuildPlanMetadata struct {
	VersionSource string
}

func Detect(nvmrcParser, buildpackYMLParser VersionParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		var requirements []packit.BuildPlanRequirement
		version, err := nvmrcParser.Parse(filepath.Join(context.WorkingDir, NvmrcSource))
		if err != nil {
			return packit.DetectResult{}, err
		}

		if version != "" {
			requirements = append(requirements, packit.BuildPlanRequirement{
				Name:    Node,
				Version: version,
				Metadata: BuildPlanMetadata{
					VersionSource: NvmrcSource,
				},
			})
		}

		version, err = buildpackYMLParser.Parse(filepath.Join(context.WorkingDir, BuildpackYMLSource))
		if err != nil {
			return packit.DetectResult{}, err
		}

		if version != "" {
			requirements = append(requirements, packit.BuildPlanRequirement{
				Name:    Node,
				Version: version,
				Metadata: BuildPlanMetadata{
					VersionSource: BuildpackYMLSource,
				},
			})
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: Node},
				},
				Requires: requirements,
			},
		}, nil
	}
}
