package nadlers

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/file"
	"net/http"
	"subscribe/internal/cash"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, cash cash.Cash) {
	res := resource{cash: cash}
	r.Get("/order/", res.get)
	r.Get("/order", file.Content("/internal/view/main.html"))
	r.Get("/*", file.Server(file.PathMap{
		"/": "/internal/static/"}))
}

type resource struct {
	cash cash.Cash
}

func (r resource) get(c *routing.Context) error {
	orderId := c.Request.FormValue("orderId")
	order, find := r.cash.GetOrder(orderId)
	if !find {
		return c.WriteWithStatus(nil, http.StatusNotFound)
	}
	return c.WriteWithStatus(order, 200)
}
