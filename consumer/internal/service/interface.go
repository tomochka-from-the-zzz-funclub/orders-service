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
}

//mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mocks
