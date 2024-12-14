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
	GetStatusSrv(order_id string) (string, string, error)
	//AddDeliveryMan(delivery_man models.DeliveryMan) (int, error)
	//CheckDeliveryStart(delivery_man_id int) (bool, error)

	FindPhoneUser(phone string) error

	Registration(user models.User) error
	CheckRegistration(user string) error

	GenerateRandomToken(phone string) (string, error)
	ValidateToken(tokenString string) (string, error)

	CheckAdmin(id_admin int) error
	ValidateTokenAdmin(tokenString string) (int, error)
	GenerateAdminToken(id int) (string, error)
}

//mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mocks
