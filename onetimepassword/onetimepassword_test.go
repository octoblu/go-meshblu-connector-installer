package onetimepassword_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/octoblu/go-meshblu-connector-installer/onetimepassword"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("onetimepassword", func() {
	var sut onetimepassword.OTP
	var err error
	var server *httptest.Server
	var lastRequest *http.Request
	var nextResponseBody string
	var nextResponseStatus int

	BeforeEach(func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lastRequest = r

			w.Header().Add("Content-Type", "applicaton/json")
			w.WriteHeader(nextResponseStatus)
			w.Write([]byte(nextResponseBody))
		}))
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("->New", func() {
		BeforeEach(func() {
			sut = onetimepassword.New("otp")
		})

		It("should return an instance", func() {
			Expect(sut).NotTo(BeNil())
		})
	})

	Describe("->NewWithURLOverride", func() {
		Describe("With a valid url", func() {
			BeforeEach(func() {
				sut, err = onetimepassword.NewWithURLOverride("otp", server.URL)
			})

			It("should produce an instance", func() {
				Expect(sut).NotTo(BeNil())
			})

			It("should not yield an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("With an invalid url", func() {
			BeforeEach(func() {
				sut, err = onetimepassword.NewWithURLOverride("otp", "")
			})

			It("should not produce an instance", func() {
				Expect(sut).To(BeNil())
			})

			It("should yield an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("with an instance", func() {
		BeforeEach(func() {
			sut, _ = onetimepassword.NewWithURLOverride("otp", server.URL)
		})

		Describe("sut.GetInformation", func() {
			Describe("when called", func() {
				BeforeEach(func() {
					sut.GetInformation()
				})

				It("Should call GET /retrieve/otp on the server", func() {
					Expect(lastRequest).NotTo(BeNil())
					Expect(lastRequest.Method).To(Equal("GET"))
					Expect(lastRequest.URL.Path).To(Equal("/retrieve/otp"))
				})
			})

			Describe("when the server is unavailable", func() {
				var info *onetimepassword.OTPInformation

				BeforeEach(func() {
					server.Close()
					nextResponseStatus = 200
					nextResponseBody = `{"uuid":"some-uuid","token":"some-token"}`
					info, err = sut.GetInformation()
				})

				It("Should return no info", func() {
					Expect(info).To(BeNil())
				})

				It("Should return an error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Describe("when the server yields json that doesn't match the schema", func() {
				var info *onetimepassword.OTPInformation

				BeforeEach(func() {
					nextResponseStatus = 200
					nextResponseBody = `{"uuid":123}`
					info, err = sut.GetInformation()
				})

				It("Should return no info", func() {
					Expect(info).To(BeNil())
				})

				It("Should return an error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Describe("when the server yields a 404", func() {
				var info *onetimepassword.OTPInformation

				BeforeEach(func() {
					nextResponseStatus = 404
					nextResponseBody = "Not Found"
					info, err = sut.GetInformation()
				})

				It("Should return no info", func() {
					Expect(info).To(BeNil())
				})

				It("Should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("Received non 200: 404"))
					Expect(err.Error()).To(ContainSubstring("Not Found"))
				})
			})

			Describe("when the server yields a valid response", func() {
				var info *onetimepassword.OTPInformation

				BeforeEach(func() {
					nextResponseStatus = 200
					nextResponseBody = `{"uuid":"some-uuid","token":"some-token"}`
					info, err = sut.GetInformation()
				})

				It("Return the uuid", func() {
					Expect(info.UUID).To(Equal("some-uuid"))
				})

				It("Return the token", func() {
					Expect(info.Token).To(Equal("some-token"))
				})

				It("Should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("sut.Expire", func() {
			Describe("when called", func() {
				BeforeEach(func() {
					sut.Expire()
				})

				It("Should call GET /expire/otp on the server", func() {
					Expect(lastRequest).NotTo(BeNil())
					Expect(lastRequest.Method).To(Equal("GET"))
					Expect(lastRequest.URL.Path).To(Equal("/expire/otp"))
				})
			})

			Describe("when the server is unavailable", func() {
				BeforeEach(func() {
					server.Close()
					nextResponseStatus = 200
					nextResponseBody = `{"uuid":"some-uuid","token":"some-token"}`
					err = sut.Expire()
				})

				It("Should return an error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Describe("when the server yields a 404", func() {
				BeforeEach(func() {
					nextResponseStatus = 404
					nextResponseBody = "Not Found"
					err = sut.Expire()
				})

				It("Should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("Received non 200: 404"))
					Expect(err.Error()).To(ContainSubstring("Not Found"))
				})
			})

			Describe("when the server yields a valid response", func() {
				var info *onetimepassword.OTPInformation

				BeforeEach(func() {
					nextResponseStatus = 200
					nextResponseBody = `{
						"uuid":  "some-uuid",
						"token": "some-token",
						"metadata": {
						  "githubSlug": "octoblu/meshblu-connector-say-hello",
							"connectorAssemblerVersion": "v13.0.0",
							"dependencyManagerVersion": "v3.0.2",
							"ignitionVersion": "v6.0.0",
							"connector": "say-hello",
							"tag": "v6.0.0"
						}
					}`
					info, err = sut.GetInformation()
				})

				It("Return the uuid", func() {
					Expect(info.UUID).To(Equal("some-uuid"))
				})

				It("Return the token", func() {
					Expect(info.Token).To(Equal("some-token"))
				})

				It("Return the Metadata", func() {
					Expect(info.Metadata.GithubSlug).To(Equal("octoblu/meshblu-connector-say-hello"))
					Expect(info.Metadata.ConnectorAssemblerVersion).To(Equal("v13.0.0"))
					Expect(info.Metadata.DependencyManagerVersion).To(Equal("v3.0.2"))
					Expect(info.Metadata.IgnitionVersion).To(Equal("v6.0.0"))
					Expect(info.Metadata.Connector).To(Equal("say-hello"))
					Expect(info.Metadata.Tag).To(Equal("v6.0.0"))
				})

				It("Should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
