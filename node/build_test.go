package node_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/node-engine-cnb/node"
	"github.com/cloudfoundry/node-engine-cnb/node/fakes"
	"github.com/cloudfoundry/node-engine-cnb/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir         string
		entryResolver     *fakes.EntryResolver
		dependencyManager *fakes.DependencyManager

		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = ioutil.TempDir("", "layers")
		Expect(err).NotTo(HaveOccurred())

		entryResolver = &fakes.EntryResolver{}
		entryResolver.ResolveCall.Returns.BuildpackPlanEntry = packit.BuildpackPlanEntry{
			Name:    "node",
			Version: "~10",
			Metadata: map[string]interface{}{
				"VersionSource": "buildpack.yml",
			},
		}

		dependencyManager = &fakes.DependencyManager{}
		dependencyManager.ResolveCall.Returns.Dependency = node.Dependency{}

		build = node.Build(entryResolver, dependencyManager)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
	})

	it("returns a result that installs node", func() {
		result, err := build(packit.BuildContext{
			Stack: "some-stack",
			Plan: packit.BuildpackPlan{
				Entries: []packit.BuildpackPlanEntry{
					{
						Name:    "node",
						Version: "~10",
						Metadata: map[string]interface{}{
							"VersionSource": "buildpack.yml",
						},
					},
				},
			},
			Layers: packit.Layers{Path: layersDir},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(packit.BuildResult{
			Plan: packit.BuildpackPlan{
				Entries: []packit.BuildpackPlanEntry{
					{
						Name:    "node",
						Version: "~10",
						Metadata: map[string]interface{}{
							"VersionSource": "buildpack.yml",
						},
					},
				},
			},
			Layers: []packit.Layer{
				{
					Name:   "node",
					Path:   filepath.Join(layersDir, "node"),
					Build:  false,
					Launch: true,
					Cache:  false,
				},
			},
		}))

		Expect(filepath.Join(layersDir, "node")).To(BeADirectory())

		Expect(entryResolver.ResolveCall.Receives.BuildpackPlanEntrySlice).To(Equal([]packit.BuildpackPlanEntry{
			{
				Name:    "node",
				Version: "~10",
				Metadata: map[string]interface{}{
					"VersionSource": "buildpack.yml",
				},
			},
		}))

		Expect(dependencyManager.ResolveCall.Receives.Entry).To(Equal(packit.BuildpackPlanEntry{
			Name:    "node",
			Version: "~10",
			Metadata: map[string]interface{}{
				"VersionSource": "buildpack.yml",
			},
		}))
		Expect(dependencyManager.ResolveCall.Receives.Stack).To(Equal("some-stack"))
		Expect(dependencyManager.InstallCall.Receives.Dependency).To(Equal(node.Dependency{}))
		Expect(dependencyManager.InstallCall.Receives.Layer).To(Equal(packit.Layer{
			Name:   "node",
			Path:   filepath.Join(layersDir, "node"),
			Build:  false,
			Launch: true,
			Cache:  false,
		}))
	})

	context("failure cases", func() {
		context("when a dependency cannot be resolved", func() {
			it.Before(func() {
				dependencyManager.ResolveCall.Returns.Error = errors.New("failed to resolve dependency")
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name:    "node",
								Version: "~10",
								Metadata: map[string]interface{}{
									"VersionSource": "buildpack.yml",
								},
							},
						},
					},
					Layers: packit.Layers{Path: layersDir},
				})
				Expect(err).To(MatchError("failed to resolve dependency"))
			})
		})

		context("when a dependency cannot be installed", func() {
			it.Before(func() {
				dependencyManager.InstallCall.Returns.Error = errors.New("failed to install dependency")
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name:    "node",
								Version: "~10",
								Metadata: map[string]interface{}{
									"VersionSource": "buildpack.yml",
								},
							},
						},
					},
					Layers: packit.Layers{Path: layersDir},
				})
				Expect(err).To(MatchError("failed to install dependency"))
			})
		})

		context("when the layers directory cannot be written to", func() {
			it.Before(func() {
				Expect(os.Chmod(layersDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(layersDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{
								Name:    "node",
								Version: "~10",
								Metadata: map[string]interface{}{
									"VersionSource": "buildpack.yml",
								},
							},
						},
					},
					Layers: packit.Layers{Path: layersDir},
				})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})
	})
}
