package repository

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/google/uuid"
)

type TimetableRepository interface{
	CountNumVenues(ctx context.Context,uniId uuid.UUID)(int64,error)
	CountNumLecturers(ctx context.Context, uniId uuid.NullUUID)(int64,error)
	CountNumCohorts(ctx context.Context,uniId uuid.UUID)(int64,error)
	CountNumCourses(ctx context.Context, uniId uuid.UUID)(int64,error)
	RetrieveTotalVenueUnavailability(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveTotalVenueUnavailabilityRow,error)
	RetrieveTotalLecturerUnavailability(ctx context.Context,uniId uuid.NullUUID)([]sqlc.RetrieveTotalLecturerUnavailabilityRow,error)
	RetrieveAllVenues(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllVenuesRow,error)
	RetrieveAllCourses(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllCoursesRow,error)
	RetrieveAllCoursesAndVenues(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllCoursesAndTheirVenueIdsRow,error)
	RetrieveAllCohorts(ctx context.Context,uniId uuid.UUID)([]sqlc.Cohort,error)
	RetrieveTotalLecturers(ctx context.Context, uniId uuid.NullUUID)([]sqlc.RetrieveTotalLecturersRow,error)
	RetrieveCohortsForAllCourses(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveCohortsForAllCoursesRow,error)

}
type timetableRepository struct {
	vq *queries.VenueQueries
	lq *queries.LecturerQueries
	cohq *queries.CohortQueries
	cq *queries.CoursesQueries
}



func NewtimeTableRepository(vq *queries.VenueQueries, lq *queries.LecturerQueries,cohq *queries.CohortQueries, cq *queries.CoursesQueries)*timetableRepository{
	return &timetableRepository{
		vq: vq,
		lq: lq,
		cohq:cohq,
		cq:cq,
	}

}


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

func (ttrp *timetableRepository) RetrieveAllVenues(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllVenuesRow,error){
	return ttrp.vq.RetrieveAllVenues(ctx,uniId)
}

func (ttrp *timetableRepository) RetrieveAllCourses(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllCoursesRow,error){
	return ttrp.cq.RetrieveAllCourses(ctx,uniId)
}

func (ttrp *timetableRepository) RetrieveAllCoursesAndVenues(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllCoursesAndTheirVenueIdsRow,error){
	return ttrp.cq.RetrieveAllCoursesAndVenues(ctx,uniId)
}

func (ttrp *timetableRepository) RetrieveAllCohorts(ctx context.Context,uniId uuid.UUID)([]sqlc.Cohort,error){
	return ttrp.cohq.RetrieveAllCohorts(ctx,uniId)
}

func (ttrp *timetableRepository) RetrieveTotalLecturers(ctx context.Context, uniId uuid.NullUUID)([]sqlc.RetrieveTotalLecturersRow,error){
	return ttrp.lq.RetrieveAllLecturers(ctx,uniId)
}

func (ttrp *timetableRepository) RetrieveTotalCohortCourses(ctx context.Context,uniId uuid.UUID)([]sqlc.RetrieveCohortsForAllCoursesRow,error){
	return ttrp.cohq.RetrieveTotalCohortCourses(ctx,uniId)
}

func (ttrp *timetableRepository) RetrieveCohortsForAllCourses(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveCohortsForAllCoursesRow,error){
	return ttrp.cohq.RetrieveTotalCohortCourses(ctx,uniId)
}