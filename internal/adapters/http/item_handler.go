package http

import (
	"encoding/json"
	"errors"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"github.com/teamcubation/go-items-challenge/internal/domain"
	"github.com/teamcubation/go-items-challenge/internal/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/presenter"
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
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_REQUEST_BODY", "Invalid request payload", map[string]interface{}{
			"error": "Check all missing fields and try again",
		}))
		return
	}
	if err := utils.ValidateStruct(&itm); err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_REQUEST_BODY", "Invalid request payload", map[string]interface{}{
			"error": "Check all missing fields and try again",
		}))
		return
	}
	createdItem, err := h.itemService.CreateItem(r.Context(), &itm)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) { //trocar para domain tirar de presenter
			middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_ITEM_NOT_FOUND", "Item not found", map[string]interface{}{
				"error": "Check item code",
			}))
			return
		}

		if errors.Is(err, domain.ErrInvalidCategory) {
			middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_CATEGORY", "Invalid category", map[string]interface{}{
				"error": "Check category",
			}))
			return
		}

		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdItem); err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
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
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_ID", "Invalid ID format", map[string]interface{}{
			"error": "Check ID format",
		}))
		return
	}
	var itm item.Item
	if err := json.NewDecoder(r.Body).Decode(&itm); err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_REQUEST_BODY", "Invalid request payload", map[string]interface{}{
			"error": "Check all missing fields and try again",
		}))
		return
	}
	if err := utils.ValidateStruct(&itm); err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_REQUEST_BODY", "Invalid request payload", map[string]interface{}{
			"error": "Check all missing fields and try again",
		}))
		return
	}
	itm.ID = id
	updatedItem, err := h.itemService.UpdateItem(r.Context(), &itm)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_ITEM_NOT_FOUND", "Item not found", map[string]interface{}{
				"error": "Check item code",
			}))
			return
		}
		if errors.Is(err, domain.ErrInvalidCategory) {
			middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_CATEGORY", "Invalid category", map[string]interface{}{
				"error": "Check category",
			}))
			return
		}

		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedItem); err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
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
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_ID", "Invalid ID format", map[string]interface{}{
			"error": "Check ID format",
		}))
		return
	}
	deletedItem, err := h.itemService.DeleteItem(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_ITEM_NOT_FOUND", "Item not found", map[string]interface{}{
				"error": "Check item code",
			}))
			return
		}

		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(deletedItem); err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}
}

// GetItemByID recupera um item pelo ID
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
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_ID", "Invalid ID format", map[string]interface{}{
			"error": "Check ID format",
		}))
		return
	}
	itm, err := h.itemService.GetItemByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_ITEM_NOT_FOUND", "Item not found", map[string]interface{}{
				"error": "Check item code",
			}))
			return
		}

		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(itm); err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
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
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_PAGE", "Invalid page number", map[string]interface{}{
			"error": "Check page number",
		}))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INVALID_LIMIT", "Invalid limit number", map[string]interface{}{
			"error": "Check limit number",
		}))
		return
	}

	items, _, err := h.itemService.ListItems(r.Context(), status, limit, page)
	if err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		middleware.ErrorHandlingMiddleware(w, presenter.New("ERR_INTERNAL_SERVER", "Internal server error", map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}
}
