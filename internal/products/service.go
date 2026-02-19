package products

import (
	"context"

	repo "github.com/P8ace/sampleGoApi/internal/adapters/postgres/repo"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductById(ctx context.Context, id int64) (repo.Product, error)
}

type svcImpl struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svcImpl{repo: repo}
}

func (s *svcImpl) ListProducts(ctx context.Context) ([]repo.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	return products, err
}

func (s *svcImpl) FindProductById(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductsById(ctx, id)
}
