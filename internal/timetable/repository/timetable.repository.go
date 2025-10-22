package repository

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/google/uuid"
)
















type timetableRepository struct {
	vq *queries.VenueQueries
	lq *queries.LecturerQueries
	cohq *queries.CohortQueries
	cq *queries.CoursesQueries
}



func NewTimeTableRepository(){}


func (ttrp *timetableRepository) CountNumVenues(ctx context.Context,uniId uuid.UUID)(int64,error){
	return ttrp.vq.RetrieveTotalVenueCount(ctx,uniId)
}

func (ttrp *timetableRepository) CountNumLecturers(ctx context.Context, uniId uuid.NullUUID)(int64,error){
	return ttrp.lq.RetrieveTotalLecturersCount(ctx,uniId)

}

func (ttrp *timetableRepository) CountNumCohorts(ctx context.Context,uniId uuid.UUID)(int64,error){
	return ttrp.cohq.CountCohortsForAUni(ctx,uniId)
}

func (ttrp *timetableRepository) CountNumCourses(ctx context.Context, uniId uuid.UUID)(int64,error){
	return ttrp.cq.CountCoursesForAUni(ctx,uniId)
}

func (ttrp *timetableRepository) RetrieveTotalVenueUnavailability(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveTotalVenueUnavailabilityRow,error){
	return ttrp.vq.RetrieveTotalVenueUnavailability(ctx,uniId)
}

func (ttrp *timetableRepository) RetrieveTotalLecturerUnavailability(ctx context.Context,uniId uuid.NullUUID)([]sqlc.RetrieveTotalLecturerUnavailabilityRow,error){
	return ttrp.lq.RetrieveTotalLecturersUnavailability(ctx,uniId)
}