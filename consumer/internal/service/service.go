package service

import (
	//store "consumer/internal/cache"
	"consumer/internal/database"
	myLog "consumer/internal/logger"
	"consumer/internal/models"
	"encoding/json"
	"fmt"

	"consumer/internal/config"

	"github.com/IBM/sarama"
)

type Srv struct {
	//cache store.Cache
	db database.InterfacePostgresDB
}

func NewSrv(cfg config.Config) *Srv {
	base := database.NewPostgres(cfg)
	//c, err := base.GetAllOrders()

	// if err != nil {
	// 	myLog.Log.Warnf("cache not restored")
	// 	//c = make(map[int]models.Order)
	// }
	myLog.Log.Debugf("cache restored!")
	return &Srv{
		//cache: store.NewStore(c),
		db: base,
	}
}

func (s *Srv) GetOrderSrv(orderUUID string) (models.Order, error) {
	myLog.Log.Debugf("GetOrderSrv with id: %+v", orderUUID)
	//order, err := s.cache.Get(orderUUID)
	// if err != nil {
	orderdb, err := s.db.GetOrder(orderUUID)
	if err != nil {
		return models.Order{}, err
	}
	//s.cache.Add(orderdb)
	return orderdb, nil
	//}
	//return order, nil
}

func (s *Srv) CreateOrder(order models.Order) (int, error) {
	myLog.Log.Debugf("CreateOrder")
	id, err := s.db.AddOrderStruct(order)
	if err != nil {
		return 0, err
	}
	// err = s.db.AddOrderStatus(order)
	// if err != nil {
	// 	return 0, err
	// }
	//s.cache.Add(order)
	return id, nil
}

func (s *Srv) CreateDeliveryMan(delivery_man models.DeliveryMan) (int, error) {
	myLog.Log.Debugf("CreateDeliveryMan")
	id, err := s.db.AddDeliveryMan(delivery_man)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Srv) UpdateStatusSrv(orderId int, status string) error {
	myLog.Log.Debugf("UpdateStatusSrv")
	err := s.db.UpdateStatus(orderId, status)
	if err != nil {
		myLog.Log.Errorf("Error UpdateStatusSrv: %+v", err.Error())
		return err
	}
	myLog.Log.Debugf("Sucsses Update satus order with id: %+v", orderId)
	return nil
}

func (s *Srv) GiveOrderDelivery(orderId int, delivery_man_id int) error {
	myLog.Log.Debugf("GiveOrderDelivery")
	err := s.db.AddDeliveryMach(orderId, delivery_man_id)
	if err != nil {
		myLog.Log.Errorf("Error GiveOrderDelivery: %+v", err.Error())
	}
	return err
}

// func (s *Srv) AddDeliveryMan(delivery_man models.DeliveryMan) error {
// 	myLog.Log.Debugf("AddDeliveryMan")
// 	err := s.db.AddDeliveryMan(delivery_man)
// 	if err != nil {
// 		myLog.Log.Errorf("Error AddDeliveryMan: %+v", err.Error())
// 	}
// 	return err
// }

func (s *Srv) CheckDeliveryStart(delovery_man_id int) (bool, error) {
	res, err := s.db.CheckDeliveryStart(delovery_man_id)
	if err != nil {
		myLog.Log.Errorf("Error SRV CheckDeliveryStart: %+v", err.Error())
		return res, err
	}
	myLog.Log.Errorf("Succes SRV CheckDeliveryStart: %+v", err.Error())
	return res, nil
}

func (s *Srv) Read(cfg config.Config) {
	myLog.Log.Debugf("Start Read")
	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)}, nil)
	if err != nil {
		myLog.Log.Errorf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	partConsumer, err := consumer.ConsumePartition(cfg.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		myLog.Log.Errorf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				myLog.Log.Debugf("Channel closed, exiting")
				return
			}

			var receivedMessage models.Order
			err := json.Unmarshal(msg.Value, &receivedMessage)

			if err != nil {
				myLog.Log.Debugf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			myLog.Log.Debugf("Received message: %+v\n", receivedMessage)

			s.CreateOrder(receivedMessage)
			//s.cache.Add(receivedMessage)
			myLog.Log.Debugf("success")
		}
	}
}
