package repository

import "context"

type InMemoryUnitOfWork struct {
	repository *InMemoryRepository
}

func NewInMemoryUnitOfWork(r *InMemoryRepository) *InMemoryUnitOfWork {
	return &InMemoryUnitOfWork{repository: r}
}

type InMemoryTx struct {
	repository *InMemoryRepository
}

func (u *InMemoryUnitOfWork) Begin(ctx context.Context) (UnitOfWorkTx, error) {
	u.repository.BeginTransaction()
	return &InMemoryTx{repository: u.repository}, nil
}

func (t *InMemoryTx) Commit() error {
	t.repository.Commit()
	return nil
}

func (t *InMemoryTx) Rollback() error {
	t.repository.Rollback()
	return nil
}
