package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	// register pgx driver name
	_ "github.com/jackc/pgx/stdlib"

	_ "github.com/lib/pq"
)

const (
	ErrInvalidTransaction   = "no valid transaction"
	ErrCantStartTransaction = "can't start transaction"
	ErrCantCloseTransaction = "can't close transaction"
)

type SqlStore interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Begin(ctx context.Context, opts *sql.TxOptions) (*SqlTx, error)
	Commit() error
	Rollback() error
	Close() error
}

type SqlDB struct {
	db *sqlx.DB
}

type SqlTx struct {
	tx *sqlx.Tx
}

func New(dbConfig Config) (*SqlDB, error) {
	db, err := sqlx.Open(dbConfig.Driver, dbConfig.URL())
	if err != nil {
		return nil, err
	}

	sqlDb := &SqlDB{db: db}
	return sqlDb, nil
}

func (s *SqlDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dest, query, args...)
}

func (s *SqlDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}

func (s *SqlDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *SqlDB) Begin(ctx context.Context, opts *sql.TxOptions) (*SqlTx, error) {
	tx, err := s.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	sqlTx := &SqlTx{tx}
	return sqlTx, nil
}

func (s *SqlDB) Commit() error {
	return errors.New(ErrInvalidTransaction)
}

func (s *SqlDB) Rollback() error {
	return errors.New(ErrInvalidTransaction)
}

func (s *SqlDB) Close() error {
	return s.db.Close()
}

func (s *SqlTx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.tx.GetContext(ctx, dest, query, args...)
}

func (s *SqlTx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.tx.SelectContext(ctx, dest, query, args...)
}

func (s *SqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.tx.ExecContext(ctx, query, args...)
}

func (s *SqlTx) Begin(ctx context.Context, opts *sql.TxOptions) (*SqlTx, error) {
	return nil, errors.New(ErrCantStartTransaction)
}

func (s *SqlTx) Commit() error {
	return s.tx.Commit()
}

func (s *SqlTx) Rollback() error {
	return s.tx.Rollback()
}

func (s *SqlTx) Close() error {
	return errors.New(ErrCantCloseTransaction)
}
