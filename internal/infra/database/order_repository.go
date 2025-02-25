package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/devfull/25-clean-architeture/internal/entity"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.DB.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")

	if err != nil {
		fmt.Printf("Error preparing statement: %v", err)
		return err
	}

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *OrderRepository) List() ([]*entity.Order, error) {
	resp, err := r.DB.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		log.Printf("Erro na query: %v", err)
		return nil, err
	}
	defer resp.Close()
	
	orders := []*entity.Order{}
	for resp.Next() {
		var order entity.Order
		if err := resp.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			log.Printf("Erro na scannnn: %v", err)
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}
