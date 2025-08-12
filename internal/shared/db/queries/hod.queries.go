package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)



type HodQueries struct{
	q *sqlc.Queries
}


func NewHodQueries(q *sqlc.Queries) *HodQueries{
	return &HodQueries{
		q:q,
	}
}

func (hq *HodQueries) RequestHodConfirmation(ctx context.Context,hod sqlc.RequestHodConfirmationParams)(sqlc.HodWaitingList,error){
	return hq.q.RequestHodConfirmation(ctx,hod)
}
func (hq *HodQueries) RetrievePendingLecturers(ctx context.Context,hodInfo sqlc.RetrievePendingLecturersParams)([]sqlc.LecturerWaitingList,error){
	return hq.q.RetrievePendingLecturers(ctx,hodInfo)
}

func (hq *HodQueries) ApproveLecturer(ctx context.Context, waitId uuid.UUID)(sqlc.LecturerWaitingList,error){
	return hq.q.ApproveLecturer(ctx,waitId)
}