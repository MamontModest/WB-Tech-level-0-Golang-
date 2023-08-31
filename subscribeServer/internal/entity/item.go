package entity

import "github.com/shopspring/decimal"

type Item struct {
	OrderUid    string
	ChrtId      int             `json:"chrt_id" validate:"required"`
	TrackNumber string          `json:"track_number" validate:"required"`
	Price       decimal.Decimal `json:"price" validate:"min=0"`
	Rid         string          `json:"rid" validate:"required"`
	Name        string          `json:"name" validate:"required"`
	Sale        int             `json:"sale" validate:"min=0"`
	Size        string          `json:"size" validate:"required"`
	TotalPrice  decimal.Decimal `json:"total_price" validate:"min=0"`
	NmId        int             `json:"nm_id" validate:"required"`
	Brand       string          `json:"brand" validate:"required"`
	Status      int             `json:"status" validate:"required"`
}
