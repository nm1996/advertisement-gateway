package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestConcurrentRouteHandler(t *testing.T) {
	// Create a new Gin router and define the route.
	router := gin.Default()
	router.POST("/advertisements", func(c *gin.Context) {
		c.Status(200)
	})

	numRequests := 10

	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Send the requests concurrently.
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()

			// Create a new HTTP request and send it to the route handler function.
			req, err := http.NewRequest("POST", "/advertisements", bytes.NewBufferString(`{
				"name": "Product A",
				"price": 19.99,
				"description": "This is the description for Product A",
				"image": "https://example.com/product-a.jpg",
				"author": "John Smith",
				"date": "2023-04-05T00:00:00Z"
			}`))

			if err != nil {
				t.Error(err)
			}

			// Use the gin test engine to handle the request and get the response.
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			// Check that the response status code is 200 OK.
			if resp.Code != http.StatusOK {
				t.Errorf("unexpected status code: got %d, want %d", resp.Code, http.StatusOK)
			}
		}()
	}

	// Wait for all goroutines to finish.
	wg.Wait()
}
