package dependencies_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/octoblu/go-meshblu-connector-installer/dependencies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Transaction struct {
	Request        *http.Request
	ResponseBody   string
	ResponseStatus int
}

var _ = Describe("Download", func() {
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
			response.Header().Add("Content-Type", "applicaton/json")
			response.WriteHeader(transaction.ResponseStatus)
			response.Write([]byte(transaction.ResponseBody))
		}))
	})

	AfterEach(func() {
		server.Close()
	})

	It("should be a thing", func() {
		Expect(dependencies.Download).NotTo(BeNil())
	})

	Describe("when called with an invalid url", func() {
		var err error

		BeforeEach(func() {
			// https: //github.com/octoblu/go-${projectName}/releases/download/${tag}/${projectName}-${platform}
			err = dependencies.DownloadWithURLAndRuntime("v1.2.0", "v3.2.1", "", dependencies.Runtime{})
		})

		It("should return an error", func() {
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("when called with an assembler and dependency manager version", func() {
		var assemblerTransaction *Transaction
		var dependencyManagerTransaction *Transaction

		BeforeEach(func() {
			assemblerTransaction = &Transaction{ResponseStatus: 200, ResponseBody: ""}
			transactions["GET /octoblu/go-meshblu-connector-assembler/releases/download/v1.2.0/meshblu-connector-assembler-darwin-amd64"] = assemblerTransaction
			dependencyManagerTransaction = &Transaction{ResponseStatus: 200, ResponseBody: ""}
			transactions["GET /octoblu/go-meshblu-connector-dependency-manager/releases/download/v3.2.1/meshblu-connector-dependency-manager-darwin-amd64"] = dependencyManagerTransaction

			dependencies.DownloadWithURLAndRuntime("v1.2.0", "v3.2.1", server.URL, dependencies.Runtime{GOOS: "darwin", GOARCH: "amd64"})
		})

		It("should pull down the correct assembler", func() {
			Expect(assemblerTransaction.Request).NotTo(BeNil())
		})

		It("should pull down the correct dependency manager", func() {
			Expect(dependencyManagerTransaction.Request).NotTo(BeNil())
		})
	})
})
