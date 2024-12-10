package service

import (
	"consumer/internal/config"
	"consumer/internal/models"
)

type InterfaceService interface {
	GetOrderSrv(orderUUID string) (models.Order, error)
	CreateOrder(order models.Order) (int, error)
	Read(cfg config.Config)

	UpdateStatusSrv(orderId int, status string) error
	GiveOrderDelivery(orderId int, delivery_man_id int) error
	CreateDeliveryMan(delivery_man models.DeliveryMan) (int, error)
	//AddDeliveryMan(delivery_man models.DeliveryMan) (int, error)
	//CheckDeliveryStart(delivery_man_id int) (bool, error)
}

//mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mocks
