package api_tests

import (
	"testing"

	"simple-gin-backend/internal/tests/testutils"
)

func TestMain(m *testing.M) {
	// Initialize the test suite with the database
	testutils.InitializeTestSuite(m)
}