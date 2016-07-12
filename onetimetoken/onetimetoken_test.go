package onetimetoken_test

import (
	"github.com/octoblu/go-meshblu-connector-installer/onetimetoken"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Onetimetoken", func() {
	Describe("constructor", func() {
		It("Should be a function", func() {
			Expect(onetimetoken.New).NotTo(BeNil())
		})
	})
})
