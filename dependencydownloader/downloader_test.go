package dependencydownloader_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"github.com/octoblu/go-meshblu-connector-installer/dependencydownloader"
	"github.com/octoblu/go-meshblu-connector-installer/osruntime"
	"github.com/spf13/afero"

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
					sut, err = dependencydownloader.NewWithoutDefaults("v1.2.3", "v3.2.1", "https://github.com", osruntime.New(), afero.NewMemMapFs())
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
					sut, err = dependencydownloader.NewWithoutDefaults("v1.2.3", "v3.2.1", "", osruntime.New(), afero.NewMemMapFs())
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

	Describe("with a OSX instance", func() {
		var Fs afero.Fs

		BeforeEach(func() {
			Fs = afero.NewMemMapFs()

			sut, err = dependencydownloader.NewWithoutDefaults("v1.2.3", "v3.2.1", server.URL, osruntime.OSRuntime{GOOS: "darwin", GOARCH: "amd64"}, Fs)
			Expect(err).To(BeNil(), "NewWithoutDefaults returned an unexpected error")
		})

		Describe("when Download is called", func() {
			var assemblerTransaction *Transaction
			var dependencyManagerTransaction *Transaction

			BeforeEach(func() {
				assemblerTransaction = &Transaction{ResponseStatus: 200, ResponseBody: `{"foo":"bar"}`}
				transactions["GET /octoblu/go-meshblu-connector-assembler/releases/download/v1.2.3/meshblu-connector-assembler-darwin-amd64"] = assemblerTransaction
				dependencyManagerTransaction = &Transaction{ResponseStatus: 200, ResponseBody: `{"test1":"test2"}`}
				transactions["GET /octoblu/go-meshblu-connector-dependency-manager/releases/download/v3.2.1/meshblu-connector-dependency-manager-darwin-amd64"] = dependencyManagerTransaction

				sut.Download()
			})

			It("should pull down the correct assembler", func() {
				Expect(assemblerTransaction.Request).NotTo(BeNil())
			})

			It("should pull down the correct dependency manager", func() {
				Expect(dependencyManagerTransaction.Request).NotTo(BeNil())
			})

			Describe("the assembler", func() {
				var fileBytes []byte

				BeforeEach(func() {
					Home, _ := os.LookupEnv("HOME")
					path := path.Join(Home, ".octoblu/MeshbluConnectors/bin/assembler-installer")

					fileBytes, err = afero.ReadFile(Fs, path)
					Expect(err).To(BeNil())
				})

				It("should save the assembler to the file system", func() {
					Expect(string(fileBytes)).To(Equal(`{"foo":"bar"}`))
				})

				XIt("should be executable (gotta figure this one out)", func() {
				})
			})

			Describe("the dependency manager", func() {
				var fileBytes []byte

				BeforeEach(func() {
					Home, _ := os.LookupEnv("HOME")
					path := path.Join(Home, ".octoblu/MeshbluConnectors/bin/dependency-manager")

					fileBytes, err = afero.ReadFile(Fs, path)
					Expect(err).To(BeNil())
				})

				It("should save the assembler to the file system", func() {
					Expect(string(fileBytes)).To(Equal(`{"test1":"test2"}`))
				})

				XIt("should be executable (gotta figure this one out)", func() {
				})
			})
		})
	})
})
