package api_clients

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/routes"
	"simple-gin-backend/internal/tests/factories"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// SetupTestRouter sets up the Gin router for testing purposes
func SetupTestRouter() *gin.Engine {
	router := gin.Default()
	// Register common middleware, routes, etc.
	routes.RegisterRoutes(router)
	return router
}

// GenerateTestJWT generates a JWT token for testing
func GenerateTestJWT(userID uint) (string, error) {
	// Define the token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AssertResponse checks the response status and body
func AssertResponse(t *testing.T, response *httptest.ResponseRecorder, expectedStatus int, expectedBody string) {
	assert.Equal(t, expectedStatus, response.Code)
	if expectedBody != "" {
		assert.JSONEq(t, expectedBody, response.Body.String())
	}
}

// TestClient represents a reusable client for API testing
type TestClient struct {
	Router *gin.Engine
	Token  string
	User   models.User
}

// PerformRequest is a helper function to perform HTTP requests in tests
func (client *TestClient) PerformRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonBody)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, _ := http.NewRequest(method, path, reqBody)

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Add the token to the Authorization header if available
	if client.Token != "" {
		req.Header.Set("Authorization", "Bearer "+client.Token)
	}

	// Set additional headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	client.Router.ServeHTTP(w, req)
	return w
}

// NewTestClientWithAuth initializes a new TestClient with authentication and creates a user in DB
func NewTestClient(authenticate bool) *TestClient {
	router := SetupTestRouter()
	client := &TestClient{Router: router}

	if authenticate {
		user := factories.UserFactory()
		token, err := GenerateTestJWT(user.ID)
		if err != nil {
			log.Fatalf("Failed to generate JWT token: %v", err)
		}
		client.Token = token
		client.User = user
	}
	return client
}
