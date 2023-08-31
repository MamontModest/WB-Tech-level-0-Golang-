package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Payment struct {
	OrderUid     string
	DeliveryCost decimal.Decimal `json:"delivery_cost" validate:"required"`
	GoodsTotal   int             `json:"goods_total" validate:"required"`
	CustomFee    decimal.Decimal `json:"custom_fee" validate:"required"`
	Amount       decimal.Decimal `json:"amount" validate:"required"`
	Transaction  string          `json:"transaction" validate:"required"`
	RequestId    string          `json:"request_id"`
	Currency     string          `json:"currency" validate:"required"`
	Provider     string          `json:"provider" validate:"required"`
	Bank         string          `json:"bank" validate:"required"`
	PaymentDt    time.Duration   `json:"payment_dt" validate:"required"`
}
