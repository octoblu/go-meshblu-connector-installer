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
		It("should produce an instance", func() {
			Expect(onetimetoken.New("otp")).NotTo(BeNil())
		})
	})

	Describe("->NewWithURLOverride", func() {
		BeforeEach(func() {
			sut = onetimetoken.NewWithURLOverride("otp", server.URL)
		})

		It("should produce an instance", func() {
			Expect(sut).NotTo(BeNil())
		})
	})

	Describe("with an instance", func() {
		BeforeEach(func() {
			sut = onetimetoken.NewWithURLOverride("otp", server.URL)
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

				It("Return the uuid", func() {
					Expect(info.UUID).To(Equal("some-uuid"))
				})

				It("Return the token", func() {
					Expect(info.Token).To(Equal("some-token"))
				})
			})
		})
	})
})
