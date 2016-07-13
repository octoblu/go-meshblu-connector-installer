package onetimetoken_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/octoblu/go-meshblu-connector-installer/onetimetoken"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("onetimetoken", func() {
	var sut onetimetoken.OTP
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
			sut = onetimetoken.New("otp")
		})

		It("should return an instance", func() {
			Expect(sut).NotTo(BeNil())
		})
	})

	Describe("->NewWithURLOverride", func() {
		Describe("With a valid url", func() {
			BeforeEach(func() {
				sut, err = onetimetoken.NewWithURLOverride("otp", server.URL)
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
				sut, err = onetimetoken.NewWithURLOverride("otp", "")
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
			sut, _ = onetimetoken.NewWithURLOverride("otp", server.URL)
		})

		Describe("sut.ExchangeForInformation", func() {
			Describe("when called", func() {
				BeforeEach(func() {
					sut.ExchangeForInformation()
				})

				It("Should call GET /retrieve/otp on the server", func() {
					Expect(lastRequest).NotTo(BeNil())
					Expect(lastRequest.Method).To(Equal("GET"))
					Expect(lastRequest.URL.Path).To(Equal("/retrieve/otp"))
				})
			})

			Describe("when the server yields a response", func() {
				var info onetimetoken.OTPInformation

				BeforeEach(func() {
					nextResponseStatus = 200
					nextResponseBody = `{"uuid":"some-uuid","token":"some-token"}`
					info = sut.ExchangeForInformation()
				})

				XIt("Return the uuid", func() {
					Expect(info.UUID).To(Equal("some-uuid"))
				})

				XIt("Return the token", func() {
					Expect(info.Token).To(Equal("some-token"))
				})
			})
		})
	})
})
