package node_test

import (
	"testing"

	"github.com/cloudfoundry/node-engine-cnb/node"
	"github.com/cloudfoundry/node-engine-cnb/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testPlanEntryResolver(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		resolver node.PlanEntryResolver
	)

	it.Before(func() {
		resolver = node.NewPlanEntryResolver()
	})

	context("when a buildpack.yml entry is included", func() {
		it("resolves the best plan entry", func() {
			entry := resolver.Resolve([]packit.BuildpackPlanEntry{
				{
					Name:    "node",
					Version: "package-json-version",
					Metadata: map[string]interface{}{
						"VersionSource": "package.json",
					},
				},
				{
					Name:    "node",
					Version: "other-version",
				},
				{
					Name:    "node",
					Version: "buildpack-yml-version",
					Metadata: map[string]interface{}{
						"VersionSource": "buildpack.yml",
					},
				},
				{
					Name:    "node",
					Version: "nvmrc-version",
					Metadata: map[string]interface{}{
						"VersionSource": ".nvmrc",
					},
				},
			})
			Expect(entry).To(Equal(packit.BuildpackPlanEntry{
				Name:    "node",
				Version: "buildpack-yml-version",
				Metadata: map[string]interface{}{
					"VersionSource": "buildpack.yml",
				},
			}))
		})
	})

	context("when a package.json entry is included", func() {
		it("resolves the best plan entry", func() {
			entry := resolver.Resolve([]packit.BuildpackPlanEntry{
				{
					Name:    "node",
					Version: "package-json-version",
					Metadata: map[string]interface{}{
						"VersionSource": "package.json",
					},
				},
				{
					Name:    "node",
					Version: "other-version",
				},
				{
					Name:    "node",
					Version: "nvmrc-version",
					Metadata: map[string]interface{}{
						"VersionSource": ".nvmrc",
					},
				},
			})
			Expect(entry).To(Equal(packit.BuildpackPlanEntry{
				Name:    "node",
				Version: "package-json-version",
				Metadata: map[string]interface{}{
					"VersionSource": "package.json",
				},
			}))
		})
	})

	context("when a .nvmrc entry is included", func() {
		it("resolves the best plan entry", func() {
			entry := resolver.Resolve([]packit.BuildpackPlanEntry{
				{
					Name:    "node",
					Version: "other-version",
				},
				{
					Name:    "node",
					Version: "nvmrc-version",
					Metadata: map[string]interface{}{
						"VersionSource": ".nvmrc",
					},
				},
			})
			Expect(entry).To(Equal(packit.BuildpackPlanEntry{
				Name:    "node",
				Version: "nvmrc-version",
				Metadata: map[string]interface{}{
					"VersionSource": ".nvmrc",
				},
			}))
		})
	})

	context("when an unknown source entry is included", func() {
		it("resolves the best plan entry", func() {
			entry := resolver.Resolve([]packit.BuildpackPlanEntry{
				{
					Name:    "node",
					Version: "other-version",
				},
			})
			Expect(entry).To(Equal(packit.BuildpackPlanEntry{
				Name:    "node",
				Version: "other-version",
			}))
		})
	})
}
