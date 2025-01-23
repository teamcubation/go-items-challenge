package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/teamcubation/go-items-challenge/internal/ports/out"
)

type categoryClient struct {
	client  *resty.Client
	baseURL string
}

func NewCategoryClient(baseURL string) out.CategoryClient {
	client := resty.New().
		SetBaseURL(baseURL).                   // Sets the base URL
		SetTimeout(10 * time.Second).          // Sets the timeout for requests
		SetRetryCount(3).                      // Retries up to 3 times if it fails
		SetRetryWaitTime(1 * time.Second).     // Time between retries
		SetRetryMaxWaitTime(10 * time.Second). // Maximum wait time between retries
		AddRetryCondition(func(r *resty.Response, _ error) bool {
			return r.StatusCode() >= 500
		})

	return &categoryClient{
		client:  client,
		baseURL: baseURL,
	}
}

type Category struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func (c *categoryClient) IsAValidCategory(_ context.Context, id int) (bool, error) {
	endpoint := fmt.Sprintf("/v1/categories/%d", id)

	var response Category

	// Makes the GET request
	resp, err := c.client.R().Get(endpoint)
	if err != nil {
		return false, fmt.Errorf("error making GET request: %w", err)
	}

	if resp.IsError() {
		return false, fmt.Errorf("error: %s", resp.String())
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return false, fmt.Errorf("error unmarshaling manually: %s", resp.String())
	}

	// Returns the decoded category
	return response.Active, nil
}
