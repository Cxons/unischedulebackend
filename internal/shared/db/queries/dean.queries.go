package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)



type DeanQueries struct {
	q *sqlc.Queries
}

func NewDeanQueries(q *sqlc.Queries) *DeanQueries{
	return &DeanQueries{
		q:q,
	}
}

func (dq *DeanQueries) RequestDeanConfirmation(ctx context.Context,dean sqlc.RequestDeanConfirmationParams)(sqlc.DeanWaitingList,error){
	return dq.q.RequestDeanConfirmation(ctx,dean)
}

func (dq *DeanQueries) RetrievePendingHods(ctx context.Context,deanInfo sqlc.RetrievePendingHodsParams)([]sqlc.HodWaitingList,error){
	return dq.q.RetrievePendingHods(ctx,deanInfo)
}

func (dq *DeanQueries) ApproveHod(ctx context.Context, waitId uuid.UUID)(sqlc.HodWaitingList,error){
	return dq.q.ApproveHod(ctx,waitId)
}