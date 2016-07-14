package dependencydownloader_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/octoblu/go-meshblu-connector-installer/dependencydownloader"
	"github.com/octoblu/go-meshblu-connector-installer/osruntime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Transaction struct {
	Request        *http.Request
	ResponseBody   string
	ResponseStatus int
}

var _ = Describe("Downloader", func() {
	var (
		server       *httptest.Server
		transactions map[string]*Transaction
		sut          dependencydownloader.Downloader
		err          error
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
			response.Header().Add("Content-Type", "applicaton/json")
			response.WriteHeader(transaction.ResponseStatus)
			response.Write([]byte(transaction.ResponseBody))
		}))
	})

	AfterEach(func() {
		server.Close()
	})

	It("should be a thing", func() {
		Expect(dependencydownloader.Download).NotTo(BeNil())
	})

	Describe("Construction", func() {
		Describe("New", func() {
			BeforeEach(func() {
				sut = dependencydownloader.New("v1.2.3", "v3.2.1")
			})

			It("Should construct a new Downloader", func() {
				Expect(sut).NotTo(BeNil())
			})
		})

		Describe("NewWithoutDefaults", func() {
			Describe("When called with valid parameters", func() {
				BeforeEach(func() {
					sut, err = dependencydownloader.NewWithoutDefaults("v1.2.3", "v3.2.1", "https://github.com", osruntime.New())
				})

				It("Should construct a new Downloader", func() {
					Expect(sut).NotTo(BeNil())
				})

				It("Should return a nil error", func() {
					Expect(err).To(BeNil())
				})
			})

			Describe("When called with an invalid url", func() {
				BeforeEach(func() {
					sut, err = dependencydownloader.NewWithoutDefaults("v1.2.3", "v3.2.1", "", osruntime.New())
				})

				It("Should not construct a new Downloader", func() {
					Expect(sut).To(BeNil())
				})

				It("Should return an error", func() {
					Expect(err).NotTo(BeNil())
				})
			})
		})
	})

	Describe("with an instance", func() {
		BeforeEach(func() {
			sut, err = dependencydownloader.NewWithoutDefaults("v1.2.3", "v3.2.1", server.URL, osruntime.OSRuntime{GOOS: "darwin", GOARCH: "amd64"})
			if err != nil {
				Fail(fmt.Sprintf("NewWithoutDefaults returned an unexpected error: %v", err.Error()))
			}
		})

		Describe("when Download is called", func() {
			var assemblerTransaction *Transaction
			var dependencyManagerTransaction *Transaction

			BeforeEach(func() {
				assemblerTransaction = &Transaction{ResponseStatus: 200, ResponseBody: ""}
				transactions["GET /octoblu/go-meshblu-connector-assembler/releases/download/v1.2.3/meshblu-connector-assembler-darwin-amd64"] = assemblerTransaction
				dependencyManagerTransaction = &Transaction{ResponseStatus: 200, ResponseBody: ""}
				transactions["GET /octoblu/go-meshblu-connector-dependency-manager/releases/download/v3.2.1/meshblu-connector-dependency-manager-darwin-amd64"] = dependencyManagerTransaction

				sut.Download()
			})

			It("should pull down the correct assembler", func() {
				Expect(assemblerTransaction.Request).NotTo(BeNil())
			})

			It("should pull down the correct dependency manager", func() {
				Expect(dependencyManagerTransaction.Request).NotTo(BeNil())
			})
		})
	})
})
