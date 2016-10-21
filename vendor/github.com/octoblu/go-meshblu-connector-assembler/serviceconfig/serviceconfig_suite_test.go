package serviceconfig_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServiceconfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Serviceconfig Suite")
}
