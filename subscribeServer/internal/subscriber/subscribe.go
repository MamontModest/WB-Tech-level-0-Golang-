package subscriber

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"subscribe/internal/entity"
)

const (
	clusterID = "test-cluster"
	clientID  = "order-subscriber"
	channel   = "order-notification"
)

type Subscriber struct {
	clusterId,
	clientId,
	channel string
}

func New() *Subscriber {
	return &Subscriber{
		clusterId: clusterID,
		clientId:  clientID,
		channel:   channel,
	}
}

func messageHandler(msg *stan.Msg) {
	order := MessageOrder{}
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Println(err)
		return
	}
	if err = order.Validate(); err != nil {
		log.Println(err)
		return
	}
	repo := NewRepository()
	err = repo.CreateOrder(order, context.Background())
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Subscriber) ConnectToSubscribe() {
	conn, err := stan.Connect(s.clusterId, s.clientId)
	if err != nil {
		log.Println(err)
	}
	_, err = conn.Subscribe(s.channel, messageHandler, stan.StartWithLastReceived())
	if err != nil {
		log.Println(err)
	}
	log.Println("subscribe complete")

}

type MessageOrder struct {
	entity.Order
}
