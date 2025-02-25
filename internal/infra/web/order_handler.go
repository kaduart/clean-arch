package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devfull/25-clean-architeture/internal/usecase"
)

type WebOrderHandler struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewWebOrderHandler(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrderUseCase *usecase.ListOrdersUseCase) *WebOrderHandler {
	return &WebOrderHandler{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrderUseCase,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Unable to parse body", http.StatusBadRequest)
		return
	}

	output, err := h.CreateOrderUseCase.Execute(dto)
	if err != nil {
		http.Error(w, "Unable to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, "Error encoder order", http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {

	orders, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		log.Printf("Failed to list orders: %v", err)
		return
	}

	// Converta cada Order para CreateOrderResponse
	/* var pbOrders []*pb.CreateOrderResponse
	for _, order := range orders {
		pbOrder := &pb.CreateOrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		}
		pbOrders = append(pbOrders, pbOrder)
	} */

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("[WEB] Erro de serialização: %v", err)
		http.Error(w, "Erro ao formatar resposta", http.StatusInternalServerError)
	}

	//return &pb.OrderListResponse{Orders: pbOrders}, nil

}
