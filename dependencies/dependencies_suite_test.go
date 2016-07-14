package dependencies_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDependencies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependencies Suite")
}
