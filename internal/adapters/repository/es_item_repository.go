package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/teamcubation/go-items-challenge/internal/ports/out"
	"strings"
	"time"
)

type Item struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CategoryID  int       `json:"category_id"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   int       `json:"updated_by"`
}

type ElasticsearchAdapter struct {
	client *elasticsearch.Client
}

func NewElasticsearchAdapter() (*ElasticsearchAdapter, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &ElasticsearchAdapter{client: client}, nil
}

type esItemRepository struct {
	esAdapter *ElasticsearchAdapter
}

func NewEsItemRepository(esAdapter *ElasticsearchAdapter) out.ItemRepository {
	return &esItemRepository{esAdapter: esAdapter}
}

func (r *esItemRepository) CreateItem(ctx context.Context, itm *item.Item) (*item.Item, error) {
	data, err := json.Marshal(itm)
	if err != nil {
		return nil, err
	}

	req := esapi.IndexRequest{
		Index:      "items",
		DocumentID: fmt.Sprintf("%d", itm.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.esAdapter.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error indexing document: %s", res.String())
	}
	return itm, nil
}

func (r *esItemRepository) GetItemByID(ctx context.Context, id int) (*item.Item, error) {
	req := esapi.GetRequest{
		Index:      "items",
		DocumentID: fmt.Sprintf("%d", id),
	}

	res, err := req.Do(ctx, r.esAdapter.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting document: %s", res.String())
	}

	var itm item.Item
	if err := json.NewDecoder(res.Body).Decode(&itm); err != nil {
		return nil, err
	}

	return &itm, nil
}

func (r *esItemRepository) UpdateItem(ctx context.Context, itm *item.Item) (*item.Item, error) {
	data, err := json.Marshal(itm)
	if err != nil {
		return nil, err
	}

	req := esapi.UpdateRequest{
		Index:      "items",
		DocumentID: fmt.Sprintf("%d", itm.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.esAdapter.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error updating document: %s", res.String())
	}
	return itm, nil
}

func (r *esItemRepository) DeleteItem(ctx context.Context, id int) (*item.Item, error) {
	req := esapi.DeleteRequest{
		Index:      "items",
		DocumentID: fmt.Sprintf("%d", id),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.esAdapter.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error deleting document: %s", res.String())
	}
	return nil, nil
}

func (e *ElasticsearchAdapter) SearchItems(ctx context.Context, query string, options map[string]interface{}) ([]Item, error) {
	req := esapi.SearchRequest{
		Index: []string{"items"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching items: %s", res.String())
	}

	var items []Item
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *esItemRepository) ItemExistsByCode(ctx context.Context, code string) bool {
	query := fmt.Sprintf(`{"query": {"match": {"code": "%s"}}}`, code)
	items, err := r.esAdapter.SearchItems(ctx, query, nil)
	return err == nil && len(items) > 0
}

func (r *esItemRepository) ListItems(ctx context.Context, status string, limit int, page int) (*item.Response, error) {
	query := fmt.Sprintf(`{"query": {"match": {"status": "%s"}}}`, status)
	repoItems, err := r.esAdapter.SearchItems(ctx, query, map[string]interface{}{
		"from": (page - 1) * limit,
		"size": limit,
	})
	if err != nil {
		return nil, err
	}

	var items []item.Item
	for _, repoItem := range repoItems {
		items = append(items, item.Item{
			ID:          repoItem.ID,
			Code:        repoItem.Code,
			Title:       repoItem.Title,
			Description: repoItem.Description,
			CategoryID:  repoItem.CategoryID,
			Price:       repoItem.Price,
			Stock:       repoItem.Stock,
			Status:      repoItem.Status,
			CreatedAt:   repoItem.CreatedAt,
			UpdatedAt:   repoItem.UpdatedAt,
			CreatedBy:   repoItem.CreatedBy,
			UpdatedBy:   repoItem.UpdatedBy,
		})
	}

	totalItems := len(items)
	totalPages := (totalItems + limit - 1) / limit
	return &item.Response{
		TotalPages: totalPages,
		Data:       items,
	}, nil
}
