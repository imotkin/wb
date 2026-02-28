package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/imotkin/L0/internal/entity"
)

func NewOrder() entity.Order {
	return entity.Order{
		UID:         uuid.New(),
		TrackNumber: uuid.NewString(),
		Entry:       "WBIL",
		Delivery: entity.Delivery{
			Name:    "Иван Иванов",
			Phone:   "+79999999999",
			Zip:     "101000",
			City:    "Москва",
			Address: "Площадь Мира, стр. 15",
			Region:  "Центральный",
			Email:   "ivanov@example.com",
		},
		Payment: entity.Payment{
			Transaction:  uuid.New(),
			RequestID:    uuid.NewString(),
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []entity.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         uuid.New(),
				Name:        "Product 1",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "ABC",
				Status:      202,
			},
			{
				ChrtID:      9934931,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         uuid.New(),
				Name:        "Product 2",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "DEF",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "sign-123",
		CustomerID:        uuid.NewString(),
		DeliveryService:   "DHL",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now().Truncate(0),
		Shard:             "1",
	}
}

func migrationsPath(t *testing.T, name string) string {
	t.Helper()

	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		migrationsDir := filepath.Join(dir, name)

		f, err := os.Stat(migrationsDir)
		if err == nil && f.IsDir() {
			return migrationsDir
		}

		f, err = os.Stat(filepath.Join(dir, "go.mod"))
		if err == nil && !f.IsDir() {
			t.Fatalf("project root was reached")
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("migrations dir is not found")
		}

		dir = parent
	}
}

func TestIntegrationPostgres(t *testing.T) {
	ctx := context.Background()

	container, err := pg.Run(
		ctx,
		"postgres:16-alpine",
		pg.WithDatabase("orders"),
		pg.WithUsername("user"),
		pg.WithPassword("secret"),
		pg.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to stop test container: %v\n", err)
		}
	})

	endpoint, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	postgres, err := New(ctx, endpoint)
	require.NoError(t, err)

	err = postgres.MigrateUp(ctx, migrationsPath(t, "migrations"))
	require.NoError(t, err)

	order := NewOrder()

	t.Run("AddOrder", func(t *testing.T) {
		inserted, err := postgres.AddOrder(ctx, order)
		require.NoError(t, err)
		require.True(t, inserted)
	})

	t.Run("GetOrder", func(t *testing.T) {
		got, err := postgres.GetOrder(ctx, order.UID)
		require.NoError(t, err)

		require.Equal(t, order, got)
	})

	t.Run("List", func(t *testing.T) {
		orders := make([]entity.Order, 0, 11)
		orders = append(orders, order) // add previous test order

		for range 10 {
			order := NewOrder()
			postgres.AddOrder(ctx, order)
			orders = append(orders, order)
		}

		got, err := postgres.List(ctx)
		require.NoError(t, err)

		require.Equal(t, orders, got)
	})
}
