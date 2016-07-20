package testserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Server has transactions and can be closed
type Server interface {
	// Close shuts down the server and blocks until all outstanding
	// requests on this server have completed.
	Close()
	// Get retrieves the transaction registered for a given
	// method/path combination. May return nil
	Get(method, path string) *Transaction
	// Set registers a transaction for a given
	// method/path combination. May be set to nil,
	// however, a request to that method/path will cause the server
	// to call the Fail handler
	Set(method, path string, transaction *Transaction)
	// URL returns the url the test server is running at
	URL() string
}

type transactionServer struct {
	transactions map[string]*Transaction
	server       *httptest.Server
}

// FailHandler is a function that will be called
// when something unexpected happens in the test
// server.
type FailHandler func(string, ...int)

// New constructs a new Server
func New(onFail FailHandler) Server {
	transactions := make(map[string]*Transaction)

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		key := fmt.Sprintf("%v %v", request.Method, request.URL.Path)

		transaction, ok := transactions[key]
		if !ok {
			onFail(fmt.Sprintf("Received an unexpected message: %v", key))
		}

		transaction.Request = request
		response.Header().Add("Content-Type", "application/octet-stream")
		response.WriteHeader(transaction.ResponseStatus)

		if len(transaction.ResponseBody) > 0 {
			response.Write(transaction.ResponseBody)
		}
		response.Write([]byte(transaction.ResponseBodyStr))
	}))

	return &transactionServer{transactions: transactions, server: server}
}

func (server *transactionServer) Close() {
	server.server.Close()
}

func (server *transactionServer) Get(method, path string) *Transaction {
	key := fmt.Sprintf("%v %v", method, path)
	return server.transactions[key]
}

func (server *transactionServer) Set(method, path string, transaction *Transaction) {
	key := fmt.Sprintf("%v %v", method, path)
	server.transactions[key] = transaction
}

func (server *transactionServer) URL() string {
	return server.server.URL
}
