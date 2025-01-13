package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	http2 "github.com/teamcubation/go-items-challenge/internal/adapters/http"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/teamcubation/go-items-challenge/internal/ports/in/mocks"
)

func setupRouter(handler *http2.ItemHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/items", handler.CreateItem).Methods(http.MethodPost)
	r.HandleFunc("/items/{id}", handler.UpdateItem).Methods(http.MethodPut)
	r.HandleFunc("/items/{id}", handler.DeleteItem).Methods(http.MethodDelete)
	r.HandleFunc("/items/{id}", handler.GetItemById).Methods(http.MethodGet)
	r.HandleFunc("/items", handler.ListItems).Methods(http.MethodGet)
	return r
}

func TestItemHandler_CreateItem(t *testing.T) {
	mockService := new(mocks.ItemService)
	handler := http2.NewItemHandler(mockService)
	router := setupRouter(handler)

	newItem := &item.Item{ID: 1, Code: "ABC", Stock: 50}

	mockService.On("CreateItem", mock.Anything, newItem).Return(newItem, nil)

	reqBody, _ := json.Marshal(newItem)
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var respItem item.Item
	err := json.Unmarshal(rec.Body.Bytes(), &respItem)
	assert.NoError(t, err)
	assert.Equal(t, newItem, &respItem)

	mockService.AssertExpectations(t)
}

func TestItemHandler_UpdateItem(t *testing.T) {
	mockService := new(mocks.ItemService)
	handler := http2.NewItemHandler(mockService)
	router := setupRouter(handler)

	itemID := 1
	existingItem := &item.Item{ID: itemID, Code: "XYZ", Stock: 10}
	mockService.On("UpdateItem", mock.Anything, existingItem).Return(existingItem, nil)

	reqBody, _ := json.Marshal(existingItem)
	req := httptest.NewRequest(http.MethodPut, "/items/"+strconv.Itoa(itemID), bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var respItem item.Item
	err := json.Unmarshal(rec.Body.Bytes(), &respItem)
	assert.NoError(t, err)
	assert.Equal(t, existingItem, &respItem)

	mockService.AssertExpectations(t)
}

func TestItemHandler_DeleteItem(t *testing.T) {
	mockService := new(mocks.ItemService)
	handler := http2.NewItemHandler(mockService)
	router := setupRouter(handler)

	itemID := 1
	deletedItem := &item.Item{ID: itemID, Code: "123"}
	mockService.On("DeleteItem", mock.Anything, itemID).Return(deletedItem, nil)

	req := httptest.NewRequest(http.MethodDelete, "/items/"+strconv.Itoa(itemID), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var respItem item.Item
	err := json.Unmarshal(rec.Body.Bytes(), &respItem)
	assert.NoError(t, err)
	assert.Equal(t, deletedItem, &respItem)

	mockService.AssertExpectations(t)
}

func TestItemHandler_GetItemById(t *testing.T) {
	mockService := new(mocks.ItemService)
	handler := http2.NewItemHandler(mockService)
	router := setupRouter(handler)

	itemID := 1
	foundItem := &item.Item{ID: itemID, Code: "456"}
	mockService.On("GetItemById", mock.Anything, itemID).Return(foundItem, nil)

	req := httptest.NewRequest(http.MethodGet, "/items/"+strconv.Itoa(itemID), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var respItem item.Item
	err := json.Unmarshal(rec.Body.Bytes(), &respItem)
	assert.NoError(t, err)
	assert.Equal(t, foundItem, &respItem)

	mockService.AssertExpectations(t)
}

func TestItemHandler_ListItems(t *testing.T) {
	mockService := new(mocks.ItemService)
	handler := http2.NewItemHandler(mockService)
	router := setupRouter(handler)

	// Create a list of pointers to items
	items := []*item.Item{
		{ID: 1, Code: "ABC", Stock: 10},
		{ID: 2, Code: "XYZ", Stock: 15},
	}
	mockService.On("ListItems", mock.Anything, "", 10, 1).Return(items, 2, nil)

	req := httptest.NewRequest(http.MethodGet, "/items?limit=10&page=1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var respItems []*item.Item
	err := json.Unmarshal(rec.Body.Bytes(), &respItems)
	assert.NoError(t, err)
	assert.Equal(t, items, respItems)

	mockService.AssertExpectations(t)
}
