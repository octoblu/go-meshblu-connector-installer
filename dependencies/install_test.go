package dependencies_test

import (
	"github.com/octoblu/go-meshblu-connector-installer/dependencies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Install", func() {
	It("should be a thing", func() {
		Expect(dependencies.Install).NotTo(BeNil())
	})
})
