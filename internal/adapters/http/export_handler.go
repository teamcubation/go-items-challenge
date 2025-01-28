package http

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"github.com/teamcubation/go-items-challenge/internal/ports/in"
)

type ExportHandler struct {
	itemService in.ItemService
}

func NewExportHandler(itemService in.ItemService) *ExportHandler {
	return &ExportHandler{itemService: itemService}
}

func (h *ExportHandler) Export(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		http.Error(w, "format query parameter is required", http.StatusBadRequest)
		return
	}

	switch format {
	case "CSV":
		h.ExportCSV(w, r)
	default:
		http.Error(w, "Unsupported format", http.StatusBadRequest)
	}
}

func (h *ExportHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	items, _, err := h.itemService.ListItems(r.Context(), "", 100, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=items.csv")
	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"ID", "Code", "Title", "Description", "Price", "Stock", "CategoryID"})

	for _, item := range items {
		err := writer.Write([]string{
			strconv.Itoa(item.ID),
			item.Code,
			item.Title,
			item.Description,
			strconv.FormatFloat(item.Price, 'f', 2, 64),
			strconv.Itoa(item.Stock),
			strconv.Itoa(item.CategoryID),
		})
		if err != nil {
			log.Printf("Error writing CSV: %v", err)
			http.Error(w, "Error writing CSV", http.StatusInternalServerError)
			return
		}
	}
}
