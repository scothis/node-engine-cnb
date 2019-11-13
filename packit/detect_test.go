package packit_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/node-engine-cnb/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		workingDir string
		tmpDir     string
	)

	it.Before(func() {
		var err error
		workingDir, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		tmpDir, err = filepath.EvalSymlinks(tmpDir)
		Expect(err).NotTo(HaveOccurred())

		Expect(os.Chdir(tmpDir)).To(Succeed())
	})

	it.After(func() {
		Expect(os.Chdir(workingDir)).To(Succeed())
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	it("provides the detect context to the given DetectFunc", func() {
		var context packit.DetectContext

		packit.Detect(nil, func(ctx packit.DetectContext) (packit.DetectResult, error) {
			context = ctx

			return packit.DetectResult{}, nil
		})

		Expect(context).To(Equal(packit.DetectContext{
			WorkingDir: tmpDir,
		}))
	})

	context("failure cases", func() {
		context("when the $PWD is set to a directory that does not exist", func() {
			var workingDir string

			it.Before(func() {
				workingDir = os.Getenv("PWD")

				Expect(os.Setenv("PWD", "/no/such/directory")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Setenv("PWD", workingDir)).To(Succeed())
			})

			it("returns an error", func() {
				var context packit.DetectContext

				packit.Detect(nil, func(ctx packit.DetectContext) (packit.DetectResult, error) {
					context = ctx

					return packit.DetectResult{}, nil
				})

				Expect(context).To(Equal(packit.DetectContext{
					WorkingDir: tmpDir,
				}))
			})
		})
	})
}
