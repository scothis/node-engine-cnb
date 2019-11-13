package packit_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

func TestUnitPackit(t *testing.T) {
	RegisterTestingT(t)

	suite := spec.New("packit", spec.Report(report.Terminal{}))

	suite("Detect", testDetect)

	suite.Run(t)
}
