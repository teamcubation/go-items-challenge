package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"strings"
)

type ElasticsearchUserAdapter struct {
	client *elasticsearch.Client
}

func NewElasticsearchUserAdapter() (*ElasticsearchUserAdapter, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &ElasticsearchUserAdapter{client: client}, nil
}

type esUserRepository struct {
	esAdapter *ElasticsearchUserAdapter
}

func NewEsUserRepository(esAdapter *ElasticsearchUserAdapter) *esUserRepository {
	return &esUserRepository{esAdapter: esAdapter}
}

func (r *esUserRepository) CreateUser(ctx context.Context, u *user.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "users",
		DocumentID: fmt.Sprintf("%d", u.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.esAdapter.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}
	return nil
}

func (r *esUserRepository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	query := fmt.Sprintf(`{"query": {"match": {"username": "%s"}}}`, username)
	req := esapi.SearchRequest{
		Index: []string{"users"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(ctx, r.esAdapter.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching document: %s", res.String())
	}

	var users []user.User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}
