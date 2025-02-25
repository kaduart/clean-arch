package graph

import (
	"context"
	"fmt"
	"log"

	"github.com/devfull/25-clean-architeture/internal/infra/graph/model"
	"github.com/devfull/25-clean-architeture/internal/usecase"
)

func (r *mutationResolver) CreateOrder(ctx context.Context, input *model.OrderInput) (*model.Order, error) {

	dto := usecase.OrderInputDTO{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}

	output, err := r.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &model.Order{
		ID:         input.ID,
		Price:      float64(output.Price),
		Tax:        float64(output.Tax),
		FinalPrice: float64(output.FinalPrice),
	}, nil
}

func (r *queryResolver) ListOrders(ctx context.Context) ([]*model.Order, error) {

	output, err := r.ListOrdersUseCase.Execute()
	if err != nil {
		log.Printf("[GraphQL] Erro no use case: %v", err)
		return nil, fmt.Errorf("erro interno ao buscar pedidos")
	}

	var orders []*model.Order
	for _, dto := range output {

		orders = append(orders, &model.Order{
			ID:         dto.ID,
			Price:      dto.Price,
			Tax:        dto.Tax,
			FinalPrice: dto.FinalPrice,
		})

	}

	return orders, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
