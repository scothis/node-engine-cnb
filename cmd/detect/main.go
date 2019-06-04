package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/nodejs-cnb/node"
	"github.com/cloudfoundry/nodejs-cnb/nvmrc"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/detect"
)

func main() {
	context, err := detect.DefaultDetect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create a default detection context: %s", err)
		os.Exit(101)
	}

	code, err := runDetect(context)
	if err != nil {
		context.Logger.Info(err.Error())
	}

	os.Exit(code)
}

func runDetect(context detect.Detect) (int, error) {
	version := context.BuildPlan[node.Dependency].Version

	nvmrcPath := filepath.Join(context.Application.Root, ".nvmrc")
	nvmrcExists, err := helper.FileExists(nvmrcPath)
	if err != nil {
		return context.Fail(), err
	}

	nvmrcVersion := version
	if nvmrcExists {
		nvmrcVersion, err = nvmrc.GetVersion(nvmrcPath, context.Logger)
		version = nvmrcVersion
		if err != nil {
			return context.Fail(), err
		}
	}

	buildpackYAMLPath := filepath.Join(context.Application.Root, "buildpack.yml")
	bpYmlExists, err := helper.FileExists(buildpackYAMLPath)
	if err != nil {
		return detect.FailStatusCode, err
	}

	buildpackYmlVersion := version
	if bpYmlExists {
		buildpackYmlVersion, err = readBuildpackYamlVersion(buildpackYAMLPath)
		version = buildpackYmlVersion
		if err != nil {
			return detect.FailStatusCode, err
		}
	}

	if bpYmlExists && nvmrcExists && buildpackYmlVersion != nvmrcVersion {
		context.Logger.Info("There is a mismatch between versions in the buildpack.yml and the .nvmrc, buildpack.yml will take precedence")
	}

	return context.Pass(buildplan.BuildPlan{
		node.Dependency: buildplan.Dependency{
			Version:  version,
			Metadata: buildplan.Metadata{"launch": true},
		},
	})

}

func readBuildpackYamlVersion(buildpackYAMLPath string) (string, error) {
	buf, err := ioutil.ReadFile(buildpackYAMLPath)
	if err != nil {
		return "", err
	}

	config := struct {
		Node struct {
			Version string `yaml:"version"`
		} `yaml:"node"`
	}{}
	if err := yaml.Unmarshal(buf, &config); err != nil {
		return "", err
	}

	return config.Node.Version, nil
}
