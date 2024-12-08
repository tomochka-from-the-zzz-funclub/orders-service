package database

import "consumer/internal/models"

type InterfacePostgresDB interface {
	AddOrderStruct(order models.Order) (int, error)
	AddOrder(order models.Order, items []int, delivery int, payment int) (int, error)
	AddItemsWithCategory(items []models.Item) ([]int, error)
	AddPayment(payment models.Payment) (int, error)
	AddDeliveryMan(delivery models.DeliveryMan) (int, error)
	AddOrderStatus(order_id int) error
	GetOrder(order_uuid string) (models.Order, error)

	UpdateStatus(order_id int, status string) error
}
