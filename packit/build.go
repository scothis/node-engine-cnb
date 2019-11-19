package packit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/cloudfoundry/node-engine-cnb/packit/exit"
)

type BuildFunc func(BuildContext) (BuildResult, error)

type BuildContext struct {
	Stack      string
	WorkingDir string
	Plan       BuildpackPlan
	Layers     Layers
}

type BuildResult struct {
	Plan      BuildpackPlan
	Layers    []Layer
	Processes []Process
}

type Process struct {
	Type    string   `toml:"type"`
	Command string   `toml:"command"`
	Args    []string `toml:"args"`
	Direct  bool     `toml:"direct"`
}

type BuildpackPlanEntry struct {
	Name     string                 `toml:"name"`
	Version  string                 `toml:"version"`
	Metadata map[string]interface{} `toml:"metadata"`
}

type BuildpackPlan struct {
	Entries []BuildpackPlanEntry `toml:"entries"`
}

func Build(f BuildFunc, options ...Option) {
	config := Config{
		exitHandler: exit.NewHandler(),
		args:        os.Args,
	}

	for _, option := range options {
		config = option(config)
	}

	var (
		layersPath = config.args[1]
		planPath   = config.args[3]
	)

	pwd, err := os.Getwd()
	if err != nil {
		config.exitHandler.Error(err)
		return
	}

	var plan BuildpackPlan
	_, err = toml.DecodeFile(planPath, &plan)
	if err != nil {
		config.exitHandler.Error(err)
		return
	}

	result, err := f(BuildContext{
		Stack:      os.Getenv("CNB_STACK_ID"),
		WorkingDir: pwd,
		Plan:       plan,
		Layers: Layers{
			Path: layersPath,
		},
	})
	if err != nil {
		config.exitHandler.Error(err)
		return
	}

	err = writeTOML(planPath, result.Plan)
	if err != nil {
		config.exitHandler.Error(err)
		return
	}

	for _, layer := range result.Layers {
		err = writeTOML(filepath.Join(layer.Path, fmt.Sprintf("%s.toml", layer.Name)), layer)
		if err != nil {
			config.exitHandler.Error(err)
			return
		}
	}

	if len(result.Processes) > 0 {
		var launch struct {
			Processes []Process `toml:"processes"`
		}
		launch.Processes = result.Processes

		err = writeTOML(filepath.Join(layersPath, "launch.toml"), launch)
		if err != nil {
			config.exitHandler.Error(err)
			return
		}
	}
}

func writeTOML(path string, value interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(value)
}
