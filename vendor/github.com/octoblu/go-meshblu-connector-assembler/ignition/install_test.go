package ignition_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/octoblu/go-meshblu-connector-assembler/ignition"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

type Transaction struct {
	Request        *http.Request
	ResponseBody   string
	ResponseStatus int
}

var _ = Describe("Installing", func() {
	var (
		server       *httptest.Server
		transactions map[string]*Transaction
	)

	BeforeEach(func() {
		transactions = make(map[string]*Transaction)

		server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			key := fmt.Sprintf("%v %v", request.Method, request.URL.Path)

			transaction, ok := transactions[key]
			if !ok {
				Fail(fmt.Sprintf("Received an unexpected message: %v", key))
			}

			transaction.Request = request
			response.Header().Add("Content-Type", "application/octet-stream")
			response.WriteHeader(transaction.ResponseStatus)
			response.Write([]byte(transaction.ResponseBody))
		}))
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Install", func() {
		It("should exist", func() {
			Expect(ignition.Install).NotTo(BeNil())
		})
	})

	Describe("InstallWithoutDefaults", func() {
		var Fs afero.Fs

		Describe("When everything goes well", func() {
			BeforeEach(func() {
				transactions["GET /path/to/ignition"] = &Transaction{ResponseStatus: 200, ResponseBody: "Nelle Reed"}

				ignitionURL := fmt.Sprintf("%v/path/to/ignition", server.URL)
				Fs = afero.NewMemMapFs()
				options := ignition.InstallOptions{IgnitionURL: ignitionURL, IgnitionPath: "/local/path/to/ignition"}
				ignition.InstallWithoutDefaults(options, Fs)
			})

			It("Should download the file", func() {
				transaction := transactions["GET /path/to/ignition"]
				Expect(transaction.Request).NotTo(BeNil())
			})

			It("Should create a directory", func() {
				exists, err := afero.Exists(Fs, "/local/path/to")

				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			})

			It("Should create a file", func() {
				exists, err := afero.Exists(Fs, "/local/path/to/ignition")

				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			})

			It("Should store the response from the server", func() {
				data, err := afero.ReadFile(Fs, "/local/path/to/ignition")

				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal("Nelle Reed"))
			})
		})

		Describe("When the server doesn't respond", func() {
			var err error

			BeforeEach(func() {
				options := ignition.InstallOptions{IgnitionURL: "http://0.0.0.0:0", IgnitionPath: "/local/path/to/ignition"}
				err = ignition.InstallWithoutDefaults(options, afero.NewMemMapFs())
			})

			It("Should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Describe("When the server yields a 404", func() {
			var err error

			BeforeEach(func() {
				transactions["GET /path/to/ignition"] = &Transaction{ResponseStatus: 404}

				ignitionURL := fmt.Sprintf("%v/path/to/ignition", server.URL)
				options := ignition.InstallOptions{IgnitionURL: ignitionURL, IgnitionPath: "/local/path/to/ignition"}
				err = ignition.InstallWithoutDefaults(options, afero.NewMemMapFs())
			})

			It("Should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("Ignition download expected 200, Received: 404"))
			})
		})

		Describe("When filesystem cannot be written to", func() {
			var err error

			BeforeEach(func() {
				transactions["GET /path/to/ignition"] = &Transaction{ResponseStatus: 200, ResponseBody: "Nelle Reed"}

				ignitionURL := fmt.Sprintf("%v/path/to/ignition", server.URL)
				options := ignition.InstallOptions{IgnitionURL: ignitionURL, IgnitionPath: "/../../ignition"}
				err = ignition.InstallWithoutDefaults(options, afero.NewReadOnlyFs(afero.NewMemMapFs()))
			})

			It("Should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
