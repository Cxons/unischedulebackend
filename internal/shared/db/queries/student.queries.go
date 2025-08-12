package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
)


type StudentQueries struct {
 q *sqlc.Queries
}

func NewStudentQueries(q *sqlc.Queries)*StudentQueries{
	return &StudentQueries{
		q:q,
	}
}
func (sq *StudentQueries) RegisterStudent(ctx context.Context,student sqlc.RegisterStudentParams)(sqlc.Student,error){
	return sq.q.RegisterStudent(ctx,student)
}
func (sq *StudentQueries) RetrieveStudentEmail(ctx context.Context,email string)(sqlc.RetrieveStudentEmailRow,error){
	return sq.q.RetrieveStudentEmail(ctx,email)
}

