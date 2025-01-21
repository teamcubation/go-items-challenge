package http

import (
	"encoding/json"
	"github.com/teamcubation/go-items-challenge/internal/utils"
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

// CreateItem cria um novo item
// @Summary Cria um novo item
// @Description Cria um novo item com os dados fornecidos no corpo da requisição
// @Tags items
// @Accept json
// @Produce json
// @Param item body item.Item true "Informações do item"
// @Success 200 {object} item.Item
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /items [post]
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var itm item.Item
	if err := json.NewDecoder(r.Body).Decode(&itm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := utils.ValidateStruct(&itm); err != nil {
		http.Error(w, "missing or invalid fields in the body", http.StatusBadRequest)
		return
	}
	createdItem, err := h.itemService.CreateItem(r.Context(), &itm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(createdItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateItem atualiza um item existente
// @Summary Atualiza um item existente
// @Description Atualiza um item existente com os dados fornecidos no corpo da requisição
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "ID do item"
// @Param item body item.Item true "Informações do item"
// @Success 200 {object} item.Item
// @Failure 404 {string} string "ID de item inválido"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /items/{id} [put]
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
	if err := utils.ValidateStruct(&itm); err != nil {
		http.Error(w, "missing or invalid fields in the body", http.StatusBadRequest)
		return
	}
	itm.ID = id
	updatedItem, err := h.itemService.UpdateItem(r.Context(), &itm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteItem deleta um item existente
// @Summary Deleta um item existente
// @Description Deleta um item existente com o ID fornecido
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "ID do item"
// @Success 200 {object} item.Item
// @Failure 400 {string} string "ID de Item não encontrado"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /items/{id} [delete]
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
	if err := json.NewEncoder(w).Encode(deletedItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetItemById recupera um item pelo ID
// @Summary Recupera um item pelo ID
// @Description Recupera um item existente com o ID fornecido
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "ID do item"
// @Success 200 {object} item.Item
// @Failute 400 {string} string "ID de item inválido"
// @Failure 404 {string} string "Item não encontrado"
// @Router /items/{id} [get]
func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	ctx := log.Context(r)
	logger := log.GetFromContext(ctx)
	logger.Info("Entering ItemHandler: GetItemById()")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	itm, err := h.itemService.GetItemByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if itm == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(itm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ListItems lista os itens
// @Summary Lista os itens
// @Description Lista os itens com base nos parâmetros fornecidos
// @Tags items
// @Accept json
// @Produce json
// @Param status query string false "Status do item"
// @Param limit query int false "Limite de itens por página"
// @Param page query int false "Página"
// @Success 200 {object} []item.Item
// @Failure 400 {string} string "Página inválida"
// @Failure 400 {string} string "Limite inválido"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /items [get]
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
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
