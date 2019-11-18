package exit_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/cloudfoundry/node-engine-cnb/packit/exit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testHandler(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		exitCode int
		stderr   *bytes.Buffer
		handler  exit.Handler
	)

	it.Before(func() {
		stderr = bytes.NewBuffer([]byte{})

		handler = exit.NewHandler(exit.WithStderr(stderr), exit.WithExitFunc(func(c int) { exitCode = c }))
	})

	it("prints the error message and exits with the right error code", func() {
		handler.Error(errors.New("some-error-message"))
		Expect(stderr).To(ContainSubstring("some-error-message"))
	})

	context("when the error is nil", func() {
		it("exits with code 0", func() {
			handler.Error(nil)
			Expect(exitCode).To(Equal(0))
		})
	})

	context("when the error is non-nil", func() {
		it("exits with code 1", func() {
			handler.Error(errors.New("failed"))
			Expect(exitCode).To(Equal(1))
		})
	})

	context("when the error is exit.Fail", func() {
		it("exits with code 1", func() {
			handler.Error(exit.Fail)
			Expect(exitCode).To(Equal(100))
		})
	})
}
