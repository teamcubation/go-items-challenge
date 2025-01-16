package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teamcubation/go-items-challenge/internal/adapters/client"

	"github.com/stretchr/testify/assert"
)

// func TestIsAValidCategory_Success(t *testing.T) {
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{"name": "electronics", "active": true}`))
// 	}))
// 	defer server.Close()

// 	categoryClient := client.NewCategoryClient(server.URL)

// 	isValid, err := categoryClient.IsAValidCategory(context.Background(), 1)

// 	assert.NoError(t, err)
// 	assert.True(t, isValid)
// }

func TestIsAValidCategory_NotActive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"name": "toys", "active": false}`))
		if err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	defer server.Close()

	categoryClient := client.NewCategoryClient(server.URL)

	isValid, err := categoryClient.IsAValidCategory(context.Background(), 2)

	assert.NoError(t, err)
	assert.False(t, isValid)

}

func TestIsAValidCategory_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	categoryClient := client.NewCategoryClient(server.URL)

	isValid, err := categoryClient.IsAValidCategory(context.Background(), 3)

	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestIsAValidCategory_RequestFailure(t *testing.T) {
	categoryClient := client.NewCategoryClient("http://invalid-url")

	isValid, err := categoryClient.IsAValidCategory(context.Background(), 4)

	assert.Error(t, err)
	assert.False(t, isValid)
}
