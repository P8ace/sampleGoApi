package products

import "context"

type Service interface {
	ListProducts(ctx context.Context) error
}

type svcImpl struct {
	//repository
}

func NewService() Service {
	return &svcImpl{}
}

func (s *svcImpl) ListProducts(ctx context.Context) error {
	return nil
}
