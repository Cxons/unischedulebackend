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

func (dq *DeanQueries) RetrieveDean(ctx context.Context,deanId uuid.UUID)(sqlc.RetrieveDeanRow,error){
	return dq.q.RetrieveDean(ctx,deanId)
}

func (dq *DeanQueries) CreateDean(ctx context.Context,deanInfo sqlc.CreateDeanParams)(sqlc.CurrentDean,error){
	return dq.q.CreateDean(ctx,deanInfo)
}

func (dq *DeanQueries) UpdateDean(ctx context.Context,deanParam sqlc.UpdateDeanParams)(sqlc.CurrentDean,error){
	return dq.q.UpdateDean(ctx,deanParam)
}
func (dq *DeanQueries) RequestDeanConfirmation(ctx context.Context,dean sqlc.RequestDeanConfirmationParams)(sqlc.DeanWaitingList,error){
	return dq.q.RequestDeanConfirmation(ctx,dean)
}

func (dq *DeanQueries) CheckDeanConfirmation(ctx context.Context,waitId uuid.UUID)(sqlc.CheckDeanConfirmationRow,error){
	return dq.q.CheckDeanConfirmation(ctx,waitId)
}

func (dq *DeanQueries) CheckDeanConfirmationWithLecturerId(ctx context.Context, lecturerId uuid.UUID)(sqlc.CheckDeanConfirmationWithLecturerIdRow,error){
	return dq.q.CheckDeanConfirmationWithLecturerId(ctx,lecturerId)
}

func (dq *DeanQueries) RetrievePendingHods(ctx context.Context,deanInfo sqlc.RetrievePendingHodsParams)([]sqlc.HodWaitingList,error){
	return dq.q.RetrievePendingHods(ctx,deanInfo)
}

func (dq *DeanQueries) ApproveHod(ctx context.Context, waitId uuid.UUID)(sqlc.HodWaitingList,error){
	return dq.q.ApproveHod(ctx,waitId)
}