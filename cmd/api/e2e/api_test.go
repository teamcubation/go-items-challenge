package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/teamcubation/go-items-challenge/cmd/api/server"
	"github.com/teamcubation/go-items-challenge/internal/adapters/client"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	ctx         context.Context
	cancel      context.CancelFunc
	mockServer  *httptest.Server
	postgresDB  *gorm.DB
	pgContainer testcontainers.Container
	apiURL      string
	token       string
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.ctx, suite.cancel = context.WithTimeout(context.Background(), 3*time.Minute)

	postgresContainer, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "testpassword",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	suite.Require().Nil(err)
	suite.pgContainer = postgresContainer

	postgresPort, _ := postgresContainer.MappedPort(suite.ctx, "5432/tcp")
	postgresHost, err := postgresContainer.Host(suite.ctx)
	suite.Require().Nil(err)

	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", postgresHost, postgresPort.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().Nil(err)

	err = db.AutoMigrate(&user.User{}, &item.Item{})
	suite.Require().Nil(err)
	suite.postgresDB = db

	suite.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/categories/1" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			response := client.Category{
				Name:   "sports",
				Active: true,
			}

			err := json.NewEncoder(w).Encode(&response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.NotFound(w, r)
		}
	}))

	r := server.SetupRouter(db, suite.mockServer.URL)

	go func() {
		server.NewServer(r, ":8081").Start(suite.ctx)
	}()
	time.Sleep(2 * time.Second)

	userPayload := `{
		"username": "test",
		"password": "ultrasecretpass"
	}`
	resp, err := http.Post("http://localhost:8081/register", "application/json", strings.NewReader(userPayload))
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusCreated, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TeardownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	suite.mockServer.Close()
	suite.Require().NoError(err)
	suite.cancel()
}

func (suite *IntegrationTestSuite) SetupTest() {
	userPayload := `{
		"username": "test",
		"password": "ultrasecretpass"
	}`
	resp, err := http.Post("http://localhost:8081/login", "application/json", strings.NewReader(userPayload))
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	suite.Assert().NoError(err)

	type LoginResponse struct {
		Token string `json:"token"`
	}
	var response LoginResponse
	err = json.Unmarshal(bodyBytes, &response)
	suite.Assert().NoError(err)

	suite.token = response.Token

	item := item.Item{
		Code:        "999999",
		Title:       "Bola",
		Description: "bola river plate",
		CategoryID:  1,
		Price:       112100,
		Stock:       5,
	}
	err = suite.postgresDB.Create(&item).Error
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TeardownTest() {
	err := suite.postgresDB.Exec("DELETE FROM items").Error
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TestCreateItemEndpoint() {
	itemPayload := `{
		"code": "434343",
		"title": "caipirinha",
		"description": "bola river plate",
		"category_id": 1,
		"price": 112100,
		"stock": 5
	}`

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8081/api/items", strings.NewReader(itemPayload))
	suite.Assert().NoError(err)

	// Agregar el header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	resp, err := client.Do(req)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusCreated, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestGetItemEndpoint() {
	resp, err := http.Get("http://localhost:8081/api/items/1")
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var item item.Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	suite.Require().NoError(err)
	suite.Require().Equal("434343", item.Code)
	suite.Require().Equal("caipirinha", item.Title)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
