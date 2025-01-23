package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/user"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/teamcubation/go-items-challenge/cmd/api/server"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ctx         context.Context
	cancel      context.CancelFunc
	mockServer  *httptest.Server
	db          *gorm.DB
	pgContainer testcontainers.Container
	apiURL      string
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Tests Suite")
}

var _ = BeforeSuite(func() {
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Minute)

	var err error
	pgContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
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
	Expect(err).NotTo(HaveOccurred())

	postgresPort, _ := pgContainer.MappedPort(ctx, "5432/tcp")
	postgresHost, err := pgContainer.Host(ctx)
	Expect(err).NotTo(HaveOccurred())

	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", postgresHost, postgresPort.Port())
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())

	err = db.AutoMigrate(&user.User{}, &item.Item{})
	Expect(err).NotTo(HaveOccurred())

	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/categories/1" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":   1,
				"name": "Bebidas",
			})
		} else {
			http.NotFound(w, r)
		}
	}))

	r := server.SetupRouter(db, mockServer.URL)

	go func() {
		server.NewServer(r, ":8081").Start(ctx)
	}()

	apiURL = "http://localhost:8081"
	time.Sleep(2 * time.Second)
})

var _ = AfterSuite(func() {
	if pgContainer != nil {
		err := pgContainer.Terminate(ctx)
		Expect(err).NotTo(HaveOccurred())
	}
	if mockServer != nil {
		mockServer.Close()
	}
	cancel()
})

var _ = Describe("Integration Tests", func() {
	BeforeEach(func() {

	})

	Describe("Item management", func() {
		Context("Creating an item", func() {
			It("should create an item successfully", func() {
				itemPayload := `{
					"code": "434343",
					"title": "Caipirinha",
					"description": "Bebida brasileña",
					"category_id": 1,
					"price": 112.10,
					"stock": 5
				}`

				resp, err := http.Post(fmt.Sprintf("%s/api/items", apiURL), "application/json", strings.NewReader(itemPayload))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				var count int64
				err = db.Model(&item.Item{}).Where("code = ?", "434343").Count(&count).Error
				Expect(err).ShouldNot(HaveOccurred())
				Expect(count).To(Equal(int64(1)))
			})
		})

		// Context("Fetching a item by ID", func() {
		// 	It("should fetch an Item successfully", func() {
		// 		resp, err := http.Get(fmt.Sprintf("%s/v1/items/1", apiURL))
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		// 		var category map[string]interface{}
		// 		err = json.NewDecoder(resp.Body).Decode(&category)
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 		Expect(category["id"]).To(Equal(float64(1)))
		// 		Expect(category["name"]).To(Equal("Bebidas"))
		// 	})
		// })
	})
})
