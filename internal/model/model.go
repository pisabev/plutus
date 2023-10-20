package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName    = "transactions"
	queryTimeout = 30 * time.Second
	poolCons     = 10
)

type TransactionType string

// Transaction types
const (
	Sale   TransactionType = "SALE"
	Credit TransactionType = "CREDIT"
	Refund TransactionType = "REFUND"
)

type Currency string

// Transaction types
const (
	Eur Currency = "EUR"
	Usd Currency = "USD"
	Gbp Currency = "GBP"
)

type Transaction struct {
	Id              int             `db:"id"`
	TransactionId   string          `db:"transaction_id"`
	TransactionType TransactionType `db:"transaction_type"`
	OrderId         string          `db:"order_id"`
	Amount          float64         `json:",string" db:"amount"`
	Currency        Currency        `db:"currency"`
	Description     string          `db:"description"`
	AccountId       string          `db:"account_id"`
}

func (m *Transaction) Validate() bool {
	return len(m.OrderId) > 0
}

type DbRepository interface {
	Insert(ent *Transaction) error
	Find(transactionId string) (*Transaction, error)
	FetchAll() ([]*Transaction, error)
	FetchAllByAccount(accountId string) ([]*Transaction, error)
	Init() error
	Empty() error
	Close()
}

type DBRepo struct {
	ctx context.Context
	db  *pgxpool.Conn
}

// NewDBRepo constructs DBRepo object.
func NewDBRepo(ctx context.Context) (DbRepository, error) {
	db, err := dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &DBRepo{
		ctx: ctx,
		db:  db,
	}, nil
}
func (m *DBRepo) Insert(ent *Transaction) error {
	ctx, cancel := context.WithTimeout(m.ctx, queryTimeout)
	defer cancel()
	ins := `
	insert into %s (id, transaction_id, transaction_type, order_id, amount, currency, description, account_id) 
	values (DEFAULT, $1, $2, $3, $4, $5, $6, $7)
	`
	_, err := m.db.Exec(ctx, fmt.Sprintf(ins, tableName),
		ent.TransactionId, ent.TransactionType, ent.OrderId, ent.Amount, ent.Currency, ent.Description, ent.AccountId)
	return err
}

func (m *DBRepo) Find(transactionId string) (*Transaction, error) {
	ctx, cancel := context.WithTimeout(m.ctx, queryTimeout)
	defer cancel()
	rows, err := m.db.Query(ctx, fmt.Sprintf("select * from %s where transaction_id = $1", tableName),
		transactionId)
	if err != nil {
		return nil, err
	}
	p, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Transaction])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (m *DBRepo) FetchAll() ([]*Transaction, error) {
	ctx, cancel := context.WithTimeout(m.ctx, queryTimeout)
	defer cancel()
	rows, err := m.db.Query(ctx, fmt.Sprintf("select * from %s", tableName))
	if err != nil {
		return nil, err
	}
	p, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Transaction])
	return p, err
}

func (m *DBRepo) FetchAllByAccount(accountId string) ([]*Transaction, error) {
	ctx, cancel := context.WithTimeout(m.ctx, queryTimeout)
	defer cancel()
	rows, err := m.db.Query(ctx, fmt.Sprintf("select * from %s where account_id = $1", tableName), accountId)
	if err != nil {
		return nil, err
	}
	p, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Transaction])
	return p, err
}

// Init creates the [tableName] if it does not exist
func (m *DBRepo) Init() error {
	ctx, cancel := context.WithTimeout(m.ctx, queryTimeout)
	defer cancel()

	tableCreate := `
	CREATE TABLE IF NOT EXISTS %s ( 
		id serial PRIMARY KEY,
		transaction_id text NOT NULL UNIQUE,
		transaction_type text NOT NULL,
		order_id text NOT NULL,
		amount decimal NOT NULL,
		currency text NOT NULL,
		description text NOT NULL,
		account_id text NOT NULL
	)`

	_, err := m.db.Exec(ctx, fmt.Sprintf(tableCreate, tableName))
	return err
}

func (m *DBRepo) Close() {
	m.db.Release()
}

func (m *DBRepo) Empty() error {
	ctx, cancel := context.WithTimeout(m.ctx, queryTimeout)
	defer cancel()
	_, err := m.db.Exec(ctx, fmt.Sprintf("truncate table %s restart identity", tableName))
	return err
}

var dbPool *pgxpool.Pool

func InitDBPool(dsn string) error {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}
	config.MaxConns = poolCons
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	dbPool = pool
	if err != nil {
		return err
	}
	return err
}
