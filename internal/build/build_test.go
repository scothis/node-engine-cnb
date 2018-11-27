package build_test

import (
	"path/filepath"
	"testing"

	"github.com/sclevine/spec/report"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/cloudfoundry/nodejs-cnb/internal/build"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func TestUnitBuild(t *testing.T) {
	RegisterTestingT(t)
	spec.Run(t, "build", testBuilds, spec.Report(report.Terminal{}))
}

func testBuilds(t *testing.T, when spec.G, it spec.S) {
	when("NewNode", func() {
		var (
			stubNodeFixture string
			f               test.BuildFactory
		)

		it.Before(func() {
			f = test.NewBuildFactory(t)
			stubNodeFixture = filepath.Join("lifecycle", "stub-node.tar.gz")
		})

		it("returns true if a build plan exists", func() {
			f.AddBuildPlan(t, build.NodeDependency, libbuildpack.BuildPlanDependency{})
			f.AddDependency(t, build.NodeDependency, stubNodeFixture)

			_, ok, err := build.NewNode(f.Build)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		})

		it("returns false if a build plan does not exist", func() {
			_, ok, err := build.NewNode(f.Build)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeFalse())
		})

		it("does not contribute node to the cache or launch layer when build and launch are not set", func() {
			f.AddBuildPlan(t, build.NodeDependency, libbuildpack.BuildPlanDependency{
				Metadata: libbuildpack.BuildPlanDependencyMetadata{},
			})
			f.AddDependency(t, build.NodeDependency, stubNodeFixture)

			nodeDep, _, err := build.NewNode(f.Build)
			Expect(err).NotTo(HaveOccurred())

			Expect(nodeDep.Contribute()).To(Succeed())

			cacheLayerRoot := filepath.Join(f.Build.Cache.Root, build.NodeDependency)
			launchLayerRoot := filepath.Join(f.Build.Launch.Root, build.NodeDependency)
			Expect(filepath.Join(cacheLayerRoot, "stub.txt")).NotTo(BeAnExistingFile())
			Expect(filepath.Join(launchLayerRoot, "stub.txt")).NotTo(BeAnExistingFile())
		})

		it("contributes node to the cache layer when build is true", func() {
			f.AddBuildPlan(t, build.NodeDependency, libbuildpack.BuildPlanDependency{
				Metadata: libbuildpack.BuildPlanDependencyMetadata{
					"build": true,
				},
			})
			f.AddDependency(t, build.NodeDependency, stubNodeFixture)

			nodeDep, _, err := build.NewNode(f.Build)
			Expect(err).NotTo(HaveOccurred())

			Expect(nodeDep.Contribute()).To(Succeed())

			layerRoot := filepath.Join(f.Build.Cache.Root, build.NodeDependency)
			Expect(filepath.Join(layerRoot, "stub.txt")).To(BeARegularFile())
		})

		it("contributes node to the launch layer when launch is true", func() {
			f.AddBuildPlan(t, build.NodeDependency, libbuildpack.BuildPlanDependency{
				Metadata: libbuildpack.BuildPlanDependencyMetadata{
					"launch": true,
				},
			})
			f.AddDependency(t, build.NodeDependency, stubNodeFixture)

			nodeDep, _, err := build.NewNode(f.Build)
			Expect(err).NotTo(HaveOccurred())

			Expect(nodeDep.Contribute()).To(Succeed())

			layerRoot := filepath.Join(f.Build.Launch.Root, build.NodeDependency)
			Expect(filepath.Join(layerRoot, "stub.txt")).To(BeARegularFile())
		})
	})
}
