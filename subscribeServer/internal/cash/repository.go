package cash

import (
	"context"
	"database/sql"
	"subscribe/internal/entity"
	"subscribe/pkg/db"
)

// Repository interface
type Repository interface {
	// UpdateCash method to UpdateCash
	UpdateCash(ctx context.Context) ([]entity.Order, error)
	// GetOrderById method to get order if order not exists in cash
	GetOrderById(orderId string, ctx context.Context) (entity.Order, error)
}

// repository
type repository struct {
	db *db.SDatabase
}

func (r repository) UpdateCash(ctx context.Context) ([]entity.Order, error) {
	orders := make(map[string]*entity.Order, 0)
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return []entity.Order{}, err
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{sql.LevelReadCommitted, true})
	if err != nil {
		return []entity.Order{}, err
	}
	defer tx.Rollback()
	query := `SELECT 
    orderUid, trackNumber, entry, locale, internalSignature, customerId, deliveryService, shardkey, smId, dateCreated, oofShard
	FROM orders`
	rows, err := tx.QueryContext(ctx, query)
	for rows.Next() {
		var id string
		order := new(entity.Order)
		rows.Scan(&id, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId,
			&order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)
		orders[id] = order
	}
	query = `SELECT 
    orderid, transaction, requestid, currency, provider, amount, paymentdt, bank, deliverycost, goodstotal, customfee
	FROM payments`
	rows, err = tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id string
		Payment := new(entity.Payment)
		rows.Scan(&id, &Payment.Transaction, &Payment.RequestId, &Payment.Currency,
			&Payment.Provider, &Payment.Amount, &Payment.PaymentDt, &Payment.Bank,
			&Payment.DeliveryCost, &Payment.GoodsTotal, &Payment.CustomFee)
		orders[id].Payment = *Payment
	}
	query = `SELECT 
    orderid, name, phone, zip, city, address, region, email
	FROM delivery`
	rows, err = tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id string
		Delivery := new(entity.Delivery)
		rows.Scan(&id, &Delivery.Name, &Delivery.Phone, &Delivery.Zip,
			&Delivery.City, &Delivery.Address, &Delivery.Region, &Delivery.Email)
		orders[id].Delivery = *Delivery
	}
	query = `SELECT 
    orderid, chrtid, tracknumber, price, rid, name, sale, size, totalprice, nmid, brand, status
	FROM item`
	rows, err = tx.QueryContext(ctx, query)
	if err != nil {
		return []entity.Order{}, err
	}
	for rows.Next() {
		var id string
		item := new(entity.Item)
		rows.Scan(&id, &item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
		)
		orders[id].Items = append(orders[id].Items, *item)
	}

	ord := make([]entity.Order, 0)
	for i, v := range orders {
		v.OrderUid = i
		ord = append(ord, *v)
	}
	return ord, nil
}
func (r repository) GetOrderById(orderId string, ctx context.Context) (entity.Order, error) {
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return entity.Order{}, err
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{sql.LevelReadCommitted, true})
	if err != nil {
		return entity.Order{}, err
	}
	defer tx.Rollback()
	order := new(entity.Order)
	query := `SELECT 
    orderUid, trackNumber, entry, locale, internalSignature, customerId, deliveryService, shardkey, smId, dateCreated, oofShard
	FROM orders where orderuid=$1`
	err = tx.QueryRowContext(ctx, query, orderId).Scan(
		&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId,
		&order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		return entity.Order{}, err
	}

	query = `SELECT 
    orderid, transaction, requestid, currency, provider, amount, paymentdt, bank, deliverycost, goodstotal, customfee
	FROM payments where orderid=$1`
	err = tx.QueryRowContext(ctx, query, orderId).Scan(
		&order.Payment.OrderUid, &order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank,
		&order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)
	if err != nil {
		return entity.Order{}, err
	}
	query = `SELECT 
    orderid, name, phone, zip, city, address, region, email
	FROM delivery where orderid=$1`
	err = tx.QueryRowContext(ctx, query, orderId).Scan(
		&order.Delivery.OrderUid, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
	)
	if err != nil {
		return entity.Order{}, err
	}
	query = `SELECT 
    orderid, chrtid, tracknumber, price, rid, name, sale, size, totalprice, nmid, brand, status
	FROM item where orderid=$1`
	rows, err := tx.QueryContext(ctx, query, orderId)
	if err != nil {
		return entity.Order{}, err
	}
	for rows.Next() {
		item := new(entity.Item)
		rows.Scan(&item.OrderUid, &item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
		)
		order.Items = append(order.Items, *item)
	}

	return *order, nil
}

func NewRepository() Repository {
	return repository{db: db.NewSDatabase()}
}
