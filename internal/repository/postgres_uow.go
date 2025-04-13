package repository

import (
	"context"
	"database/sql"
)

type PgUnitOfWork struct {
	db *sql.DB
}

func NewPgUnitOfWork(db *sql.DB) UnitOfWork {
	return &PgUnitOfWork{db: db}
}

func (u *PgUnitOfWork) Begin(ctx context.Context) (UnitOfWorkTx, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &PgUnitOfWorkTx{tx: tx}, nil
}

type PgUnitOfWorkTx struct {
	tx *sql.Tx
}

func (t *PgUnitOfWorkTx) Commit() error {
	return t.tx.Commit()
}
func (t *PgUnitOfWorkTx) Rollback() error {
	return t.tx.Rollback()
}
