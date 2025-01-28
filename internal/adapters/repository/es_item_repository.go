package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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

func (e *ElasticsearchAdapter) CreateItem(ctx context.Context, item Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "items",
		DocumentID: fmt.Sprintf("%d", item.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}
	return nil
}

func (e *ElasticsearchAdapter) UpdateItem(ctx context.Context, item Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      "items",
		DocumentID: fmt.Sprintf("%d", item.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating document: %s", res.String())
	}
	return nil
}

func (e *ElasticsearchAdapter) DeleteItem(ctx context.Context, id int) error {
	req := esapi.DeleteRequest{
		Index:      "items",
		DocumentID: fmt.Sprintf("%d", id),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
	}
	return nil
}

func (e *ElasticsearchAdapter) SearchItems(ctx context.Context, query string, filters map[string]interface{}) ([]Item, error) {
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"multi_match": map[string]interface{}{
						"query":  query,
						"fields": []string{"title", "description"},
					}},
				},
				"filter": filters,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex("items"),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	var items []Item
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var item Item
		source := hit.(map[string]interface{})["_source"]
		data, _ := json.Marshal(source)
		json.Unmarshal(data, &item)
		items = append(items, item)
	}

	return items, nil
}
