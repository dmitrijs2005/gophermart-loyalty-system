package service

import "github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"

type BaseService struct{}

func (s *BaseService) EndTransaction(tx repository.UnitOfWorkTx, errPtr *error) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if *errPtr != nil {
		_ = tx.Rollback()
	} else {
		*errPtr = tx.Commit()
	}
}
