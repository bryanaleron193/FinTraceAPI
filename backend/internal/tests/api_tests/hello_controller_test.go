package api_tests

import (
	"net/http"
	api_clients "simple-gin-backend/internal/tests/clients"
	"testing"
)

func TestGetHelloWorld(t *testing.T) {
	// Set up the Gin router using the reusable function
	client := api_clients.NewTestClient(false)

	// Perform a request to the route
	response := client.PerformRequest("GET", "/", nil, nil)
	// Assert the response
	api_clients.AssertResponse(t, response, http.StatusOK, `"Hello, world!"`)
}
