package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
)

func TestValidateItemMiddleware_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up a router with the middleware and a simple handler to inspect the context
	router := gin.New()
	router.Use(middleware.ValidateItem())
	router.POST("/items", func(c *gin.Context) {
		newItem, exists := c.Get("newItem")
		assert.True(t, exists)

		item, ok := newItem.(item.Item)
		assert.True(t, ok)
		assert.Equal(t, "Test Item", item.Title)
		assert.Equal(t, 100.00, item.Price)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Create a request with a valid JSON body
	itemJSON := `{"title": "Test Item", "price": 100.00}`
	req, _ := http.NewRequest("POST", "/items", bytes.NewBufferString(itemJSON))
	rec := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status": "success"}`, rec.Body.String())
}

func TestValidateItemMiddleware_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up a router with the middleware
	router := gin.New()
	router.Use(middleware.ValidateItem())
	router.POST("/items", func(c *gin.Context) {
		// This should not be reached if middleware works correctly
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Create a request with an invalid JSON body
	invalidJSON := `{"title": 123}` // Assuming "title" should be a string
	req, _ := http.NewRequest("POST", "/items", bytes.NewBufferString(invalidJSON))
	rec := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error": "invalid request body"}`, rec.Body.String())
}
