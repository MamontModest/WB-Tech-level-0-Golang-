package db

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"time"
)

const (
	dsn    = "postgres://localhost:8001/l0?sslmode=disable&user=userL0&password=userL0"
	driver = "postgres"
)

type SDatabase struct {
	dsn    string
	driver string
}

func NewSDatabase() *SDatabase {
	return &SDatabase{
		dsn:    dsn,
		driver: driver,
	}
}

func (database *SDatabase) conn() (*sql.DB, error) {
	db, err := sql.Open(database.driver, database.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (database *SDatabase) ConnWith(ctx context.Context) (*sql.DB, error) {
	ct, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	resp := make(chan ResponseConn)
	go func() {
		conn, err := database.conn()
		resp <- ResponseConn{conn, err}
	}()
	for {
		select {
		case <-ct.Done():
			return nil, errors.New("timeout connect")
		case val := <-resp:
			return val.conn, val.err
		}
	}
}

type ResponseConn struct {
	conn *sql.DB
	err  error
}
