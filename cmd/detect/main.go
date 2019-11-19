package main

import (
	"github.com/cloudfoundry/node-engine-cnb/node"
	"github.com/cloudfoundry/node-engine-cnb/packit"
)

func main() {
	nvmrcParser := node.NewNvmrcParser()
	buildpackYMLParser := node.NewBuildpackYMLVersionParser()

	packit.Detect(node.Detect(nvmrcParser, buildpackYMLParser))
}
