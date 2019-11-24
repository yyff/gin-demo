package main

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type Order struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

func getUserOrders(ctx context.Context, db *sql.DB, id int) ([]*Order, error) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM orders WHERE user_id=?`, id)
	if err != nil {
		return nil, errors.Wrap(err, "querying orders by user")
	}
	defer rows.Close()
	orders := []*Order{}
	for rows.Next() {
		order := &Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID); err != nil {
			return nil, errors.Wrap(err, "scan order")
		}
		orders = append(orders, order)
	}

	return orders, nil
}
func getOrder(ctx context.Context, db *sql.DB, id int) (*Order, error) {
	row := db.QueryRowContext(ctx, `SELECT * FROM orders WHERE id=?`, id)
	var order Order
	if err := row.Scan(&order.ID, &order.UserID, &order.ProductID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "querying order")
	}
	return &order, nil
}

func createOrder(ctx context.Context, db *sql.DB, order *Order) (int, error) {
	var orderID int
	result, err := db.ExecContext(ctx, "INSERT INTO orders (user_id, product_id) VALUES (?,?) ", order.UserID, order.ProductID)
	if err != nil {
		return -1, err
	}
	rowID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	orderID = int(rowID)
	return orderID, nil
}
