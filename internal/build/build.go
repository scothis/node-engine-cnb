package build

import (
	"fmt"

	"github.com/cloudfoundry/libjavabuildpack"
)

const NodeDependency = "node"

type Node struct {
	BuildContribution, LaunchContribution bool
	CacheLayer                            libjavabuildpack.DependencyCacheLayer
	LaunchLayer                           libjavabuildpack.DependencyLaunchLayer
}

func NewNode(builder libjavabuildpack.Build) (n Node, planExists bool, e error) {
	bp, planExists := builder.BuildPlan[NodeDependency]
	if !planExists {
		return Node{}, false, nil
	}

	deps, err := builder.Buildpack.Dependencies()
	if err != nil {
		return Node{}, false, err
	}

	dep, err := deps.Best(NodeDependency, bp.Version, builder.Stack)
	if err != nil {
		return Node{}, false, err
	}

	node := Node{}

	if _, contributeBuild := bp.Metadata["build"]; contributeBuild {
		node.BuildContribution = true
		node.CacheLayer = builder.Cache.DependencyLayer(dep)
	}

	if _, contributeLaunch := bp.Metadata["launch"]; contributeLaunch {
		node.LaunchContribution = true
		node.LaunchLayer = builder.Launch.DependencyLayer(dep)
	}

	return node, true, nil
}

var environment = map[string]string{
	"NODE_ENV":              "production",
	"NODE_MODULES_CACHE":    "true",
	"NODE_VERBOSE":          "false",
	"NPM_CONFIG_PRODUCTION": "true",
	"NPM_CONFIG_LOGLEVEL":   "error",
	"WEB_MEMORY":            "512",
	"WEB_CONCURRENCY":       "1",
}

func (n Node) Contribute() error {
	if n.BuildContribution {
		err := n.CacheLayer.Contribute(func(artifact string, layer libjavabuildpack.DependencyCacheLayer) error {
			layer.Logger.SubsequentLine("Expanding to %s", layer.Root)
			if err := libjavabuildpack.ExtractTarGz(artifact, layer.Root, 1); err != nil {
				return err
			}

			layer.Logger.SubsequentLine("Writing NODE_HOME")
			layer.OverrideEnv("NODE_HOME", layer.Root)

			for key, value := range environment {
				layer.Logger.SubsequentLine("Writing " + key)
				layer.OverrideEnv(key, value)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	if n.LaunchContribution {
		err := n.LaunchLayer.Contribute(func(artifact string, layer libjavabuildpack.DependencyLaunchLayer) error {
			layer.Logger.SubsequentLine("Expanding to %s", layer.Root)
			if err := libjavabuildpack.ExtractTarGz(artifact, layer.Root, 1); err != nil {
				return err
			}

			layer.Logger.SubsequentLine("Writing profile.d/NODE_HOME")
			layer.WriteProfile("NODE_HOME", fmt.Sprintf("export NODE_HOME=\"%s\"", layer.Root))

			for key, value := range environment {
				layer.Logger.SubsequentLine("Writing profile.d/" + key)
				layer.WriteProfile(key, fmt.Sprintf("export %s=\"%s\"", key, value))
			}

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
