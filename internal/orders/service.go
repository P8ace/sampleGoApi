package orders

import (
	"context"
	"fmt"

	"github.com/P8ace/sampleGoApi/internal/adapters/postgres/repo"
	"github.com/jackc/pgx/v5"
)

type OrderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int64 `json:"quantity"`
}

type CreateOrder struct {
	CustomerID int64       `json:"customerId"`
	Items      []OrderItem `json:"items"`
}

type Service interface {
	PlaceOrder(ctx context.Context, order CreateOrder) (repo.Order, error)
}

type OrderService struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewOrderService(repo *repo.Queries, db *pgx.Conn) *OrderService {
	return &OrderService{
		repo: repo,
		db:   db,
	}
}

func (o *OrderService) PlaceOrder(ctx context.Context, order CreateOrder) (repo.Order, error) {
	// validate payload
	if order.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer Id is required")
	}
	if len(order.Items) == 0 {
		return repo.Order{}, fmt.Errorf("At least one item is required")
	}

	tx, err := o.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := o.repo.WithTx(tx)

	// create an order
	neworder, err := qtx.CreateOrder(ctx, order.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	// look if the product exists
	for _, item := range order.Items {
		product, err := qtx.FindProductsById(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, err
		}

		if product.Quantity < int32(item.Quantity) {
			return repo.Order{}, err
		}
		// create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:      neworder.ID,
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PriceInCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, err
		}

	}
	// create order item in the database
	tx.Commit(ctx)

	return neworder, nil
}
