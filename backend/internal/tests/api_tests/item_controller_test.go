package api_tests

import (
	"encoding/json"
	"net/http"
	"simple-gin-backend/internal/models"
	api_clients "simple-gin-backend/internal/tests/clients"
	"simple-gin-backend/internal/tests/factories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {
	// Initialize the test client with authentication (creates a user and token)
	client := api_clients.NewTestClient(true)

	// Create items for the authenticated user using the factory
	item1 := factories.ItemFactory(client.User.ID)
	item2 := factories.ItemFactory(client.User.ID)

	user2 := factories.UserFactory()
	factories.ItemFactory(user2.ID)

	// Perform a GET request to /items
	response := client.PerformRequest("GET", "/items", nil, nil)

	// Assert that the response code is 200 OK
	assert.Equal(t, http.StatusOK, response.Code)

	// Define the expected items response (make sure the structure matches the JSON output)
	expectedBody := []models.Item{
		item1,
		item2,
	}

	// Unmarshal the response to a slice of items and assert it matches the expected data
	var actualItems []models.Item
	err := json.Unmarshal(response.Body.Bytes(), &actualItems)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Assert that the returned items match what was created
	assert.Equal(t, expectedBody, actualItems)
}
