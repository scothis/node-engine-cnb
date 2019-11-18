package exit_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitExit(t *testing.T) {
	suite := spec.New("packit/exit", spec.Report(report.Terminal{}))
	suite("Handler", testHandler)
	suite.Run(t)
}
