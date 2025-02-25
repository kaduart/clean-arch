package usecase

import (
	"github.com/devfull/25-clean-architeture/internal/entity"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(repo entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{OrderRepository: repo}
}

func (uc *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := uc.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	var output []OrderOutputDTO
	for _, order := range orders {
		output = append(output, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return output, nil
}
