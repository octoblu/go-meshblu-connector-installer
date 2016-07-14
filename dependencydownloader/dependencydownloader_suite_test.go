package dependencydownloader_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDependencydownloader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependencydownloader Suite")
}
