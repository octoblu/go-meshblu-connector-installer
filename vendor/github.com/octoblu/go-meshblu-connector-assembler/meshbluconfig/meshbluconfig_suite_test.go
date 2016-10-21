package meshbluconfig_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMeshbluconfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Meshbluconfig Suite")
}
