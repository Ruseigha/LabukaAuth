package integration

import (
	"os"
	"testing"
	"time"

	"github.com/Ruseigha/LabukaAuth/test/testutil"
)

func TestMain(m *testing.M) {
	// Wait for MongoDB to be ready
	// WHY: In CI/CD, MongoDB container might still be starting
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	if err := testutil.WaitForMongoDB(mongoURI, 30*time.Second); err != nil {
		panic("MongoDB not available: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}