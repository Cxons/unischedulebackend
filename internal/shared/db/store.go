package db

import (
	"context"
	"database/sql"
)



type Store interface{
	ExecTx(ctx context.Context,fn func(*Queries)error)error
}

type store struct{
	db *sql.DB
	*Queries
}



func NewStore(db *sql.DB) *store{
	return &store{
		db: db,
		Queries: New(db),
	}
}


func (s *store) ExecTx(ctx context.Context, fn func(*Queries)error)error{
	tx,err := s.db.BeginTx(ctx,nil)
	if err != nil{
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil{
			return rbErr
		}
		return err
	}
	return tx.Commit()
}