package test

import (
	"bytes"
	"encoding/json"
	"gateway/controllers"
	"gateway/interfaces"
	"gateway/services"
	"gateway/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const ADV_TEST_JSON = `
	[
		{
			"name" : "Product A",
			"price" : 19.99,
			"description" : "This is the description for Product A",
			"image" : "https://example.com/product-a.jpg",
			"author" : "John Doe",
			"date" : "2023-04-01T00:00:00Z"
		},
		{
			"name" : "Product A",
			"price" : 19.99,
			"description" : "This is the description for Product A",
			"image" : "https://example.com/product-a.jpg",
			"author" : "John Doe",
			"date" : "2023-04-01T00:00:00Z"
		},
		{
			"name" : "Product A",
			"price" : 19.99,
			"description" : "This is the description for Product A",
			"image" : "https://example.com/product-a.jpg",
			"author" : "John Doe",
			"date" : "2023-04-01T00:00:00Z"
		}
	]
`

type TestSuite struct {
	Router     *gin.Engine
	TestServer *httptest.Server
	Service    interfaces.AdvertisementService
	Controller interfaces.AdvertisementController
}

func (suite *TestSuite) setup() {
	suite.Router = gin.Default()
	suite.TestServer = httptest.NewServer(suite.Router)

	// Initialize your service and controller objects here
	logger := log.New(os.Stdout, "[myapp_test] ", log.LstdFlags|log.Lmicroseconds|log.LUTC)
	queue := utils.NewSafeQueue()
	suite.Service = services.NewAdvertisementService(logger, queue)
	suite.Controller = controllers.NewAdvertisementController(suite.Service, logger)

	suite.Router.POST("/advertisements", suite.Controller.ReceiveAdvertisements)
	suite.Router.GET("/get-queue-size", suite.Controller.GetAdvertisementsCount)
}

func (suite *TestSuite) teardown() {
	suite.TestServer.Close()
}

func (suite *TestSuite) TestConcurrentRouteHandler(t *testing.T) {
	numRequests := 100

	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Send the requests concurrently.
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()

			// Create a new HTTP request and send it to the route handler function.
			req, err := http.NewRequest("POST", "/advertisements", bytes.NewBufferString(ADV_TEST_JSON))

			if err != nil {
				t.Error(err)
			}

			// Use the gin test engine to handle the request and get the response.
			resp := httptest.NewRecorder()
			suite.Router.ServeHTTP(resp, req)

			// Check that the response status code is 200 OK.
			if resp.Code != http.StatusOK {
				t.Errorf("unexpected status code: got %d, want %d", resp.Code, http.StatusOK)
			}
		}()
	}

	// Wait for all goroutines to finish.
	wg.Wait()
	time.Sleep(1 * time.Second)
}

func (suite *TestSuite) TestQueueResponseCount(t *testing.T) {

	// Create a new HTTP request and send it to the route handler function.
	req, err := http.NewRequest("GET", "/get-queue-size", nil)

	if err != nil {
		t.Error(err)
	}

	// Use the gin test engine to handle the request and get the response.
	resp := httptest.NewRecorder()
	suite.Router.ServeHTTP(resp, req)

	// Check that the response status code is 200 OK.
	if resp.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %d, want %d", resp.Code, http.StatusOK)
	}

	var responseDataHolder map[string]string
	err = json.Unmarshal(resp.Body.Bytes(), &responseDataHolder)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if responseDataHolder["message"] != "300" {
		t.Errorf("unexpected response data: got %v, want %d", responseDataHolder["message"], 300)
	}
}

// Run test logic

func TestSuiteRunner(t *testing.T) {
	suite := &TestSuite{}
	suite.setup()
	defer suite.teardown()

	t.Run("Test1", suite.TestConcurrentRouteHandler)
	t.Run("Test2", suite.TestQueueResponseCount)
	// Run more tests here
}
