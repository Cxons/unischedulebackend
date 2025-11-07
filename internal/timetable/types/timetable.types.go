package types

import (
	"time"

	"github.com/google/uuid"
)



type CustomSessionPlacement struct {
	SessionIdx int32
	CourseId uuid.UUID
	VenueId uuid.UUID
	Day string
	SessionTime time.Time
	UniversityId uuid.UUID
}