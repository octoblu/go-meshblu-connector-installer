package onetimetoken_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOnetimetoken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Onetimetoken Suite")
}
