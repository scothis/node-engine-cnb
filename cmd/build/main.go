package main

import (
	"github.com/cloudfoundry/node-engine-cnb/node"
	"github.com/cloudfoundry/node-engine-cnb/packit"
)

func main() {
	packit.Build(node.Build(nil, nil))
}
