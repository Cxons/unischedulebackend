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

func (hq *HodQueries) RetrieveHod(ctx context.Context,hodId uuid.UUID)(sqlc.RetrieveHodRow,error){
	return hq.q.RetrieveHod(ctx,hodId)
}

func (hq *HodQueries) CreateHod(ctx context.Context,hodInfo sqlc.CreateHodParams)(sqlc.CurrentHod,error){
	return hq.q.CreateHod(ctx,hodInfo)
}

func (hq *HodQueries) UpdateHod(ctx context.Context,hodParams sqlc.UpdateHodParams)(sqlc.CurrentHod,error){
	return hq.q.UpdateHod(ctx,hodParams)
}

func (hq *HodQueries) RequestHodConfirmation(ctx context.Context,hod sqlc.RequestHodConfirmationParams)(sqlc.HodWaitingList,error){
	return hq.q.RequestHodConfirmation(ctx,hod)
}

func (hq *HodQueries) CheckHodConfirmation(ctx context.Context, waitId uuid.UUID)(sqlc.CheckHodConfirmationRow,error){
	return hq.q.CheckHodConfirmation(ctx,waitId)
}

func (hq *HodQueries) RetrievePendingLecturers(ctx context.Context,hodInfo sqlc.RetrievePendingLecturersParams)([]sqlc.LecturerWaitingList,error){
	return hq.q.RetrievePendingLecturers(ctx,hodInfo)
}

func (hq *HodQueries) ApproveLecturer(ctx context.Context, waitId uuid.UUID)(sqlc.LecturerWaitingList,error){
	return hq.q.ApproveLecturer(ctx,waitId)
}

func (hq *HodQueries) CheckHodConfirmationWithLecturerId(ctx context.Context, lecturerId uuid.UUID)(sqlc.CheckHodConfirmationWithLecturerIdRow,error){
	return hq.q.CheckHodConfirmationWithLecturerId(ctx,lecturerId)
}
