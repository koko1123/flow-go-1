package rest

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/access"
	"github.com/onflow/flow-go/engine/access/rest/middleware"
)

// all route names
const (
	getTransactionByIDRoute          = "getTransactionByID"
	getTransactionResultByIDRoute    = "getTransactionResultByID"
	createTransactionRoute           = "createTransaction"
	getBlocksByIDRoute               = "getBlocksByIDs"
	getBlocksByHeightRoute           = "getBlocksByHeight"
	getCollectionByIDRoute           = "getCollectionByID"
	executeScriptRoute               = "executeScript"
	getBlockPayloadByIDRoute         = "getBlockPayloadByID"
	getExecutionResultByBlockIDRoute = "getExecutionResultByBlockID"
	getExecutionResultByIDRoute      = "getExecutionResultByID"
	getAccountRoute                  = "getAccount"
)

// NewServer returns an HTTP server initialized with the REST API handler
func NewServer(backend access.API, listenAddress string, logger zerolog.Logger) (*http.Server, error) {

	router, err := initRouter(backend, logger)
	if err != nil {
		return nil, err
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
			http.MethodHead},
	})

	return &http.Server{
		Addr:         listenAddress,
		Handler:      c.Handler(router),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}, nil
}

func initRouter(backend access.API, logger zerolog.Logger) (*mux.Router, error) {
	router := mux.NewRouter().StrictSlash(true)
	v1SubRouter := router.PathPrefix("/v1").Subrouter()

	// common middleware for all request
	v1SubRouter.Use(middleware.LoggingMiddleware(logger))
	v1SubRouter.Use(middleware.QueryExpandable())
	v1SubRouter.Use(middleware.QuerySelect())

	var linkGenerator LinkGenerator = NewLinkGeneratorImpl(v1SubRouter)

	// create a schema validation
	validation, err := newSchemaValidation()
	if err != nil {
		return nil, err
	}

	for _, r := range routeDefinitions() {
		h := NewHandler(logger, backend, r.handler, linkGenerator, validation)
		v1SubRouter.
			Methods(r.method).
			Path(r.pattern).
			Name(r.name).
			Handler(h)
	}
	return router, nil
}

type routeDefinition struct {
	name       string
	method     string
	pattern    string
	validators []ApiValidatorFunc
	handler    ApiHandlerFunc
}

func routeDefinitions() []routeDefinition {
	return []routeDefinition{
		// Transactions
		{
			method:  http.MethodGet,
			pattern: "/transactions/{id}",
			name:    getTransactionByIDRoute,
			handler: getTransactionByID,
		}, {
			method:  http.MethodPost,
			pattern: "/transactions",
			name:    createTransactionRoute,
			handler: createTransaction,
		},
		// Transaction Results
		{
			method:  http.MethodGet,
			pattern: "/transaction_results/{id}",
			name:    getTransactionResultByIDRoute,
			handler: getTransactionResultByID,
		},
		// Blocks
		{
			method:  http.MethodGet,
			pattern: "/blocks/{id}",
			name:    getBlocksByIDRoute,
			handler: getBlocksByIDs,
		}, {
			method:  http.MethodGet,
			pattern: "/blocks",
			name:    getBlocksByHeightRoute,
			handler: getBlocksByHeight,
		},
		// Block Payload
		{
			method:  http.MethodGet,
			pattern: "/blocks/{id}/payload",
			name:    getBlockPayloadByIDRoute,
			handler: getBlockPayloadByID,
		},
		// Execution Result
		{
			method:  http.MethodGet,
			pattern: "/execution_results/{id}",
			name:    getExecutionResultByIDRoute,
			handler: getExecutionResultByID,
		},
		{
			method:  http.MethodGet,
			pattern: "/execution_results",
			name:    getExecutionResultByBlockIDRoute,
			handler: getExecutionResultsByBlockIDs,
		},
		// Collections
		{
			method:  http.MethodGet,
			pattern: "/collections/{id}",
			name:    getCollectionByIDRoute,
			handler: getCollectionByID,
		},
		// Scripts
		{
			method:  http.MethodPost,
			pattern: "/scripts",
			name:    executeScriptRoute,
			handler: executeScript,
		},
		// Accounts
		{
			method:  http.MethodGet,
			pattern: "/accounts/{address}",
			name:    getAccountRoute,
			handler: getAccount,
		}}
}
