package subscriber

import (
	"context"
	"database/sql"
	"log"
	"subscribe/pkg/db"
)

// Repository interface
type Repository interface {
	// CreateOrder method to create Order
	CreateOrder(message MessageOrder, ctx context.Context) error
}

// repository
type repository struct {
	db *db.SDatabase
}

// CreateOrder method to create Order
func (r repository) CreateOrder(message MessageOrder, ctx context.Context) error {
	conn, err := r.db.ConnWith(ctx)
	defer conn.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `INSERT INTO orders 
    (orderUid, trackNumber, entry, locale, internalSignature, customerId, deliveryService, shardkey, smId, dateCreated, oofShard)
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	rows, err := tx.QueryContext(ctx, query,
		message.OrderUid, message.TrackNumber, message.Entry, message.Locale, message.InternalSignature, message.CustomerId,
		message.DeliveryService, message.Shardkey, message.SmId, message.DateCreated, message.OofShard,
	)
	if err != nil {
		return err
	}
	rows.Close()
	query = `INSERT INTO payments 
    (orderid, transaction, requestid, currency, provider, amount, paymentdt, bank, deliverycost, goodstotal, customfee) 
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	rows, err = tx.QueryContext(ctx, query,
		message.OrderUid, message.Payment.Transaction, message.Payment.RequestId, message.Payment.Currency,
		message.Payment.Provider, message.Payment.Amount, message.Payment.PaymentDt, message.Payment.Bank,
		message.Payment.DeliveryCost, message.Payment.GoodsTotal, message.Payment.CustomFee,
	)
	if err != nil {
		return err
	}
	rows.Close()
	query = `INSERT into delivery
	(orderid, name, phone, zip, city, address, region, email)
	values ($1, $2, $3, $4, $5, $6, $7, $8)`
	rows, err = tx.QueryContext(ctx, query,
		message.OrderUid, message.Delivery.Name, message.Delivery.Phone, message.Delivery.Zip,
		message.Delivery.City, message.Delivery.Address, message.Delivery.Region, message.Delivery.Email,
	)
	if err != nil {
		return err
	}
	rows.Close()
	for _, item := range message.Items {
		query = `insert into item
		(orderid, chrtid, tracknumber, price, rid, name, sale, size, totalprice, nmid, brand, status) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		rows, err = tx.QueryContext(ctx, query,
			message.OrderUid, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
		rows.Close()
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func NewRepository() Repository {
	return repository{db: db.NewSDatabase()}
}
