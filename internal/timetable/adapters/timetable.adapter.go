package adapters

import (
	"time"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)


type CohortSessionRow sqlc.GetCohortSessionsInCurrentTimetableRow
type StudentSessionRow sqlc.GetStudentTimetableSessionsRow



func (r CohortSessionRow) GetDay() string        { return r.Day }
func (r CohortSessionRow) GetSessionTime() time.Time { return r.SessionTime }
func (r CohortSessionRow) GetCourseID() uuid.UUID    { return r.CourseID }
func (r CohortSessionRow) GetVenueID() uuid.UUID     { return r.VenueID}
func (r CohortSessionRow) GetSessionID() uuid.UUID   { return r.SessionID }

func (r StudentSessionRow) GetDay() string        { return r.Day }
func (r StudentSessionRow) GetSessionTime() time.Time { return r.SessionTime }
func (r StudentSessionRow) GetCourseID() uuid.UUID    { return r.CourseID }
func (r StudentSessionRow) GetVenueID() uuid.UUID     { return r.VenueID }
func (r StudentSessionRow) GetSessionID() uuid.UUID   { return r.SessionID }
