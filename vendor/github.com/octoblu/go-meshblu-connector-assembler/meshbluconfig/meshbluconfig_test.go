package meshbluconfig_test

import (
	"github.com/octoblu/go-meshblu-connector-assembler/meshbluconfig"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Meshbluconfig", func() {
	Describe("Write", func() {
		It("should not be nil", func() {
			Expect(meshbluconfig.Write).NotTo(BeNil())
		})
	})

	Describe("WriteWithFS", func() {
		var fs afero.Fs

		BeforeEach(func() {
			fs = afero.NewMemMapFs()
		})

		Describe("When called with a directory and meshblu params", func() {
			BeforeEach(func() {
				options := meshbluconfig.Options{
					DirPath:  "/path/to/connector",
					UUID:     "the-uuid",
					Token:    "a-token",
					Hostname: "a-host",
					Port:     100,
				}
				meshbluconfig.WriteWithFS(options, fs)
			})

			It("Should write meshbluJSON file in memory file system", func() {
				ok, err := afero.Exists(fs, "/path/to/connector/meshblu.json")
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())
			})

			It("Should match properties", func() {
				meshbluJson, err := afero.ReadFile(fs, "/path/to/connector/meshblu.json")
				Expect(err).To(BeNil())
				Expect(meshbluJson).Should(MatchJSON(`{"uuid": "the-uuid", "token": "a-token", "hostname": "a-host", "port": 100}`))
			})
		})
	})
})
