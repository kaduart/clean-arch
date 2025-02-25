package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devfull/25-clean-architeture/internal/usecase"
)

type ListOrdersHandler struct {
	ListOrdersUseCase usecase.ListOrdersUseCase
}

func NewListOrdersHandler(uc usecase.ListOrdersUseCase) *ListOrdersHandler {
	return &ListOrdersHandler{ListOrdersUseCase: uc}
}

func (h *ListOrdersHandler) List(w http.ResponseWriter, r *http.Request) {
	orders, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		http.Error(w, "Failed to list orders", http.StatusInternalServerError)
		return
	}

	log.Printf("BATEUU no LIST %v", orders)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
