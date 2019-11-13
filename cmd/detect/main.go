package main

import "os"
import "github.com/cloudfoundry/node-engine-cnb/packit"
import "github.com/cloudfoundry/node-engine-cnb/node"

func main() {
	packit.Detect(os.Args, node.Detect())
}
