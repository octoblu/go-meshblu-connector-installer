package osruntime_test

import (
	"runtime"

	"github.com/octoblu/go-meshblu-connector-installer/osruntime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OSRuntime", func() {
	Describe("New", func() {
		It("Should construct a new runtime instance that matches the current process", func() {
			sut := osruntime.New()
			Expect(sut.GOARCH).To(Equal(runtime.GOARCH))
			Expect(sut.GOOS).To(Equal(runtime.GOOS))
		})
	})
})
