package graph

import "github.com/devfull/25-clean-architeture/internal/usecase"

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}
