package serviceconfig_test

import (
	"github.com/octoblu/go-meshblu-connector-assembler/serviceconfig"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceConfig", func() {
	Describe("Write", func() {
		It("should not be nil", func() {
			Expect(serviceconfig.Write).NotTo(BeNil())
		})
	})

	Describe("WriteWithFS", func() {
		var fs afero.Fs

		BeforeEach(func() {
			fs = afero.NewMemMapFs()
		})

		Describe("When called with a directory and Service params", func() {
			BeforeEach(func() {
				options := serviceconfig.Options{
					ServiceName:   "a-service-name",
					ServiceType:   "Service",
					DisplayName:   "a-display-name",
					Description:   "a-description",
					ConnectorName: "a-connector-name",
					GithubSlug:    "a-github-slug",
					Tag:           "some-tag",
					BinPath:       "/some/bin/path",
					Dir:           "/path/to/connector/",
					LogDir:        "/path/to/logs/",
				}
				serviceconfig.WriteWithFS(options, fs)
			})

			It("Should write service json file in memory file system", func() {
				ok, err := afero.Exists(fs, "/path/to/connector/service.json")
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())
			})

			It("Should match properties", func() {
				serviceJson, err := afero.ReadFile(fs, "/path/to/connector/service.json")
				Expect(err).To(BeNil())
				Expect(serviceJson).Should(MatchJSON(`{
				  "ServiceName": "a-service-name",
					"ServiceType": "Service",
				  "DisplayName": "a-display-name",
				  "Description": "a-description",
				  "ConnectorName": "a-connector-name",
				  "GithubSlug": "a-github-slug",
				  "Tag": "some-tag",
				  "BinPath": "/some/bin/path",
				  "Dir": "/path/to/connector/",
				  "Stderr": "/path/to/logs/connector-error.log",
				  "Stdout": "/path/to/logs/connector.log"
				}`))
			})
		})
	})
})
