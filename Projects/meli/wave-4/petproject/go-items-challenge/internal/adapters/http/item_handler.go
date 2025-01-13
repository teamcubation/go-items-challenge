package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/teamcubation/go-items-challenge/internal/ports/in"
	"github.com/teamcubation/go-items-challenge/pkg/log"
)

type ItemHandler struct {
	itemService in.ItemService
}

func NewItemHandler(itemService in.ItemService) *ItemHandler {
	return &ItemHandler{itemService: itemService}
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var itm item.Item
	if err := json.NewDecoder(r.Body).Decode(&itm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdItem, err := h.itemService.CreateItem(r.Context(), &itm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(createdItem)
}

func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	var itm item.Item
	if err := json.NewDecoder(r.Body).Decode(&itm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	itm.ID = id
	updatedItem, err := h.itemService.UpdateItem(r.Context(), &itm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updatedItem)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	deletedItem, err := h.itemService.DeleteItem(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(deletedItem)
}

func (h *ItemHandler) GetItemById(w http.ResponseWriter, r *http.Request) {
	ctx := log.Context(r)
	logger := log.GetFromContext(ctx)
	logger.Info("Entering ItemHandler: GetItemById()")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	itm, err := h.itemService.GetItemById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if itm == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(itm)
}

func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		http.Error(w, "Invalid page", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	items, _, err := h.itemService.ListItems(r.Context(), status, limit, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}
