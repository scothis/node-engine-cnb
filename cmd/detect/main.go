package main

import (
	"github.com/cloudfoundry/node-engine-cnb/node"
	"github.com/cloudfoundry/node-engine-cnb/packit"
)

func main() {
	buildpackYMLParser := node.NewBuildpackYMLParser()

	packit.Detect(node.Detect(nil, buildpackYMLParser))
}
