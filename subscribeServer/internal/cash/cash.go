package cash

import (
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
	"subscribe/internal/entity"
	"time"
)

type Cash interface {
	GetOrder(orderId string) (entity.Order, bool)
	SaveOrder(order entity.Order)
	Reload()
}

const (
	Expiration = 5 * time.Minute
)

var LastTimeUpdate = time.Now().UTC()

type cash struct {
	Repository     *cache.Cache
	LastTimeUpdate time.Time
	Expiration     time.Duration
}

func NewInit() Cash {
	c := cash{
		Repository:     cache.New(3*time.Minute, 5*time.Minute),
		LastTimeUpdate: LastTimeUpdate,
		Expiration:     Expiration,
	}
	c.Reload()
	return c
}

func (c cash) SaveOrder(order entity.Order) {
	c.Repository.Set(order.OrderUid, order, cache.NoExpiration)
	if c.LastTimeUpdate.Add(c.Expiration).Before(time.Now().UTC()) {
		c.Reload()
		c.LastTimeUpdate = time.Now().UTC()
	}
}

func (c cash) GetOrder(orderId string) (entity.Order, bool) {
	ord, find := c.Repository.Get(orderId)
	if !find {
		repo := NewRepository()
		order, err := repo.GetOrderById(orderId, context.Background())
		if err != nil {
			fmt.Println(err)
			return entity.Order{}, false
		}
		c.SaveOrder(order)
		return order, true
	}
	return ord.(entity.Order), true
}

func (c cash) Reload() {
	repo := NewRepository()
	orders, err := repo.UpdateCash(context.Background())
	if err != nil {
		panic(err)
	}
	for _, v := range orders {
		c.SaveOrder(v)
	}
}
