package dto

import (
	"time"

	"github.com/google/uuid"
)





type CreateATimeTableDto struct{
	StartTime time.Time `json:"startTime" validate:"required"`
	EndTime time.Time `json:"endTime" validate:"required"`
	UniversityId uuid.UUID `json:"universityId" validate:"required"`
}



