package database

import "consumer/internal/models"

type InterfacePostgresDB interface {
	AddOrderStruct(order models.Order) (int, error)
	AddOrder(order models.Order, items []int, payment int) (int, error)
	AddItemsWithCategory(items []models.Item) ([]int, error)
	AddPayment(payment models.Payment) (int, error)
	AddDeliveryMan(delivery_man models.DeliveryMan) (int, error)
	AddOrderStatus(order_id int) error
	GetOrder(order_uuid string) (models.Order, error)
	AddDeliveryMach(order_id int, delivery_man_id int) error

	UpdateStatus(order_id int, status string) error
	CheckDeliveryStart(delivery_man_id int) (bool, error)
}
