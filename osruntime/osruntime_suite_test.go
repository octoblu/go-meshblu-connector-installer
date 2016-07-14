package osruntime_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOsruntime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Osruntime Suite")
}
