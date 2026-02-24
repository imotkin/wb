package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/imotkin/L0/internal/entity"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, url string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Postgres{pool: pool}, nil
}

func (p *Postgres) MigrateUp(ctx context.Context, path string) error {
	goose.SetLogger(goose.NopLogger())

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(p.pool)
	defer db.Close()

	return goose.UpContext(ctx, db, path)
}

func (p *Postgres) MigrateDown(ctx context.Context, path string) error {
	goose.SetLogger(goose.NopLogger())

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(p.pool)
	defer db.Close()

	return goose.DownContext(ctx, db, path)
}

func (p *Postgres) List(ctx context.Context) ([]entity.Order, error) {
	query :=
		`SELECT
            o.id, o.track_number, o.entry, o.locale, o.internal_signature,
            o.customer_id, o.delivery_service, o.shardkey, o.sm_id,
            o.date_created, o.oof_shard,
            d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
            p.transaction, p.request_id, p.currency, p.provider, p.amount,
            EXTRACT(EPOCH FROM p.payment_dt)::BIGINT, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
			(SELECT COALESCE(jsonb_agg(item), '[]'::jsonb) 
     			FROM (SELECT * FROM items WHERE order_id = o.id) item
    		) AS items_json
        FROM orders o
        LEFT JOIN deliveries d ON d.order_id = o.id
        LEFT JOIN payments p ON p.order_id = o.id
	`

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
		IsoLevel:   pgx.RepeatableRead,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	var order entity.Order
	var itemsJSON []byte

	for rows.Next() {
		fields := []any{
			&order.UID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID,
			&order.DateCreated, &order.Shard,

			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,

			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
			&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
			&order.Payment.CustomFee,

			&itemsJSON,
		}

		err = rows.Scan(fields...)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(itemsJSON, &order.Items)
		if err != nil {
			return nil, fmt.Errorf("decode order items: %w", err)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (p *Postgres) GetOrder(ctx context.Context, id uuid.UUID) (entity.Order, error) {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
		IsoLevel:   pgx.ReadCommitted,
	})
	if err != nil {
		return entity.Order{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	orderQuery := `
        SELECT
            o.id, o.track_number, o.entry, o.locale, o.internal_signature,
            o.customer_id, o.delivery_service, o.shardkey, o.sm_id,
            o.date_created, o.oof_shard,
            d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
            p.transaction, p.request_id, p.currency, p.provider, p.amount,
            EXTRACT(EPOCH FROM p.payment_dt)::BIGINT, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
        FROM orders o
        LEFT JOIN deliveries d ON d.order_id = o.id
        LEFT JOIN payments p ON p.order_id = o.id
        WHERE o.id = $1`

	var order entity.Order

	fields := []any{
		&order.UID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID,
		&order.DateCreated, &order.Shard,

		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,

		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	}

	err = tx.QueryRow(ctx, orderQuery, id).Scan(fields...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, entity.ErrOrderNotFound
		}

		return entity.Order{}, err
	}

	itemsQuery := `
		SELECT chrt_id, track_number, price, rid, name, 
			   sale, size, total_price, nm_id, brand, status
          FROM items WHERE order_id = $1`

	rows, err := tx.Query(ctx, itemsQuery, order.UID)
	if err != nil {
		return entity.Order{}, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Item])
	if err != nil {
		return entity.Order{}, err
	}

	order.Items = items

	return order, tx.Commit(ctx)
}

func (p *Postgres) AddOrder(ctx context.Context, order entity.Order) (bool, error) {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
		IsoLevel:   pgx.ReadCommitted,
	})
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	inserted, err := p.addOrder(ctx, tx, order)
	if err != nil {
		return false, fmt.Errorf("failed to add order: %w", err)
	}

	if !inserted {
		return true, nil
	}

	err = p.addDelivery(ctx, tx, order.UID, order.Delivery)
	if err != nil {
		return false, fmt.Errorf("failed to add delivery: %w", err)
	}

	err = p.addPayment(ctx, tx, order.UID, order.Payment)
	if err != nil {
		return false, fmt.Errorf("failed to add payment: %w", err)
	}

	for _, item := range order.Items {
		err = p.addItem(ctx, tx, order.UID, item)
		if err != nil {
			return false, fmt.Errorf("failed to add item %q: %w", item.Name, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *Postgres) addOrder(ctx context.Context, tx pgx.Tx, order entity.Order) (bool, error) {
	query := `
		INSERT INTO orders (
			id, track_number, entry, locale, internal_signature, customer_id, 
			delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, DEFAULT)
		ON CONFLICT DO NOTHING`

	fields := []any{
		order.UID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
	}

	tag, err := tx.Exec(ctx, query, fields...)

	return tag.RowsAffected() > 0, err
}

func (p *Postgres) addItem(ctx context.Context, tx pgx.Tx, orderID uuid.UUID, item entity.Item) error {
	query := `
		INSERT INTO items (
			order_id, chrt_id, track_number, price, rid, name,
			sale, size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	fields := []any{
		orderID,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.RID,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmID,
		item.Brand,
		item.Status,
	}

	_, err := tx.Exec(ctx, query, fields...)

	return err
}

func (p *Postgres) addDelivery(ctx context.Context, tx pgx.Tx, orderID uuid.UUID, delivery entity.Delivery) error {
	query := `
		INSERT INTO deliveries (
			order_id, name, phone, zip,
			city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	fields := []any{
		orderID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	}

	_, err := tx.Exec(ctx, query, fields...)

	return err
}

func (p *Postgres) addPayment(ctx context.Context, tx pgx.Tx, orderID uuid.UUID, payment entity.Payment) error {
	query := `
		INSERT INTO payments (
			order_id, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, to_timestamp($7), $8, $9, $10, $11)`

	fields := []any{
		orderID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	}

	_, err := tx.Exec(ctx, query, fields...)

	return err
}
