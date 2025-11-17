package repository

import (
	"context"
	"fmt"
	"log/slog"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/timetable/types"
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
	CreateACandidateTimeTable(ctx context.Context, candidateData sqlc.CreateCandidateParams, sessionPlacements []types.CustomSessionPlacement)error
	DeprecateLatestCandidate(ctx context.Context,uniId uuid.UUID)error
	RestoreCurrentCandidate(ctx context.Context,uniId uuid.UUID)error
	FetchSessionsForACohort(ctx context.Context,params sqlc.GetCohortSessionsInCurrentTimetableParams)([]sqlc.GetCohortSessionsInCurrentTimetableRow,error)
	FetchSessionsForAStudent(ctx context.Context,studentId uuid.UUID)([]sqlc.GetStudentTimetableSessionsRow,error)
}
type timetableRepository struct {
	vq *queries.VenueQueries
	lq *queries.LecturerQueries
	cohq *queries.CohortQueries
	cq *queries.CoursesQueries
	tmtq *queries.TimeTableQueries
	store sqlc.Store
}



func NewtimeTableRepository(vq *queries.VenueQueries, lq *queries.LecturerQueries,cohq *queries.CohortQueries, cq *queries.CoursesQueries, tmtq *queries.TimeTableQueries,store sqlc.Store)*timetableRepository{
	return &timetableRepository{
		vq: vq,
		lq: lq,
		cohq:cohq,
		cq:cq,
		tmtq: tmtq,
		store: store,
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
// func parseTimeFromString(timeStr string) (time.Time, error) {
//     if timeStr == "" {
//         return time.Time{}, fmt.Errorf("empty time string")
//     }
    
//     // Remove any timezone information if present
//     timeStr = strings.Split(timeStr, "+")[0]
//     timeStr = strings.Split(timeStr, "-")[0]
//     timeStr = strings.TrimSpace(timeStr)
    
//     // Try different time formats
//     formats := []string{
//         "15:04:05",
//         "15:04:05.999999",
//         "15:04",
//         "2006-01-02 15:04:05",
//         "2006-01-02T15:04:05",
//         time.RFC3339,
//     }
    
//     for _, format := range formats {
//         if t, err := time.Parse(format, timeStr); err == nil {
//             // Extract only the time part (ignore date)
//             return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, time.UTC), nil
//         }
//     }
    
//     return time.Time{}, fmt.Errorf("unable to parse time string: %s", timeStr)
// }
func (ttrp *timetableRepository) CreateACandidateTimeTable(ctx context.Context, candidateData sqlc.CreateCandidateParams, sessionPlacements []types.CustomSessionPlacement) error {
    return ttrp.store.ExecTx(ctx, func(q *sqlc.Queries) error {
        val, createCandidateErr := q.CreateCandidate(ctx, candidateData)
        if createCandidateErr != nil {
            return createCandidateErr
        }

        slog.Info("Creating session placements", "count", len(sessionPlacements))

        for i, placement := range sessionPlacements {
            params := sqlc.CreateSessionPlacementsParams{
                CandidateID:  val.ID,
                SessionIdx:   placement.SessionIdx,
                CourseID:     placement.CourseId,
                VenueID:      placement.VenueId,
                Day:          placement.Day,
                SessionTime:  placement.SessionTime,
                UniversityID: placement.UniversityId,
            }

            createSessionPlacementsErr := q.CreateSessionPlacements(ctx, params)
            if createSessionPlacementsErr != nil {
                slog.Error("Failed to create session placement", 
                    "sessionIdx", placement.SessionIdx,
                    "index", i,
                    "error", createSessionPlacementsErr,
                    "sessionTime", placement.SessionTime,
                    "sessionTimeType", fmt.Sprintf("%T", placement.SessionTime))
                return fmt.Errorf("failed to create session placement for session %d: %w", placement.SessionIdx, createSessionPlacementsErr)
            }
        }
        
        slog.Info("Successfully created all session placements", "count", len(sessionPlacements))
        return nil
    })
}

func (ttrp *timetableRepository) DeprecateLatestCandidate(ctx context.Context,uniId uuid.UUID)error{
	return ttrp.tmtq.DeprecateLatestCandidate(ctx,uniId)
}

func (ttrp *timetableRepository) RestoreCurrentCandidate(ctx context.Context,uniId uuid.UUID)error{
	return ttrp.tmtq.RestoreCurrentCandidate(ctx,uniId)
}

func (ttrp *timetableRepository) FetchSessionsForACohort(ctx context.Context,params sqlc.GetCohortSessionsInCurrentTimetableParams)([]sqlc.GetCohortSessionsInCurrentTimetableRow,error){
	return ttrp.cq.FetchSessionsForACohort(ctx,params)
}

func (ttrp *timetableRepository) FetchSessionsForAStudent(ctx context.Context,studentId uuid.UUID)([]sqlc.GetStudentTimetableSessionsRow,error){
	return ttrp.cq.FetchSessionsForAStudent(ctx,studentId)
}