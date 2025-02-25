package service

import (
	"context"
	"fmt"

	"github.com/devfull/25-clean-architeture/internal/infra/grpc/protofiles/pb"
	"github.com/devfull/25-clean-architeture/internal/usecase"
)

// temos um ordeservice que recebe um usecase de criacao de order
type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewOrderService(CreateOrderUseCase usecase.CreateOrderUseCase,
	listUseCase *usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: CreateOrderUseCase,
		ListOrdersUseCase:  listUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	fmt.Println("CreateOrder called")
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.Blank) (*pb.OrderListResponse, error) {
	output, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var orders []*pb.CreateOrderResponse
	for _, o := range output {
		orders = append(orders, &pb.CreateOrderResponse{
			Id:         o.ID,
			Price:      float32(o.Price),
			Tax:        float32(o.Tax),
			FinalPrice: float32(o.FinalPrice),
		})
	}

	return &pb.OrderListResponse{Orders: orders}, nil
}
