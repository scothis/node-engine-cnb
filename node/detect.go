package node

import "github.com/cloudfoundry/node-engine-cnb/packit"

func Detect() packit.DetectFunc {
	return func(packit.DetectContext) (packit.DetectResult, error) {
		return packit.DetectResult{}, nil
	}
}
