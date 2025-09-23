package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)



type LecturerQueries struct {
	q *sqlc.Queries
}

func NewLecturerQueries(q *sqlc.Queries)*LecturerQueries{
	return &LecturerQueries{
		q:q,
	}
}
func (lq *LecturerQueries) RegisterLecturer(ctx context.Context, lecturer sqlc.RegisterLecturerParams)(sqlc.Lecturer,error){
	return lq.q.RegisterLecturer(ctx,lecturer)
}
func (lq *LecturerQueries) RetrieveLecturerEmail(ctx context.Context, email string)(sqlc.RetrieveLecturerEmailRow,error){
	return lq.q.RetrieveLecturerEmail(ctx,email)
}

func (lq *LecturerQueries) RequestLecturerConfirmation(ctx context.Context, lecturer sqlc.RequestLecturerConfirmationParams)(sqlc.LecturerWaitingList,error){
	return lq.q.RequestLecturerConfirmation(ctx,lecturer)
}

func (lq *LecturerQueries) CheckLecturerConfirmation(ctx context.Context,waitId uuid.UUID)(sqlc.CheckLecturerConfirmationRow,error){
	return lq.q.CheckLecturerConfirmation(ctx,waitId)
}

func (lq *LecturerQueries) RetrieveTotalLecturersUnavailability(ctx context.Context, universityId uuid.NullUUID)([]sqlc.RetrieveTotalLecturerUnavailabilityRow,error){
	return lq.q.RetrieveTotalLecturerUnavailability(ctx,universityId)
}

func (lq *LecturerQueries) RetrieveTotalLecturersCount(ctx context.Context,universityId uuid.NullUUID)(int64,error){
	return lq.q.CountTotalLecturers(ctx,universityId)
}