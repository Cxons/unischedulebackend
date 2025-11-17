package dto

import (
	"time"

	"github.com/google/uuid"
)

type RetrieveAllDepartmentsDto struct{
	FacultyId string `json:"FacultyId" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
}

type UniversityResponse struct{
	Id uuid.UUID
	Name string
	Abbr string
	Email string
	Website string
	PhoneNumber string
	Address string
	CurrentSession string
	Logo string
}

type FacultyResponse struct{
	Id uuid.UUID
	UniversityId uuid.UUID
	Name string
	Code string
}


type DepartmentResponse struct{
	Id uuid.UUID
	UniversityId uuid.UUID
	FacultyId uuid.UUID
	Name string
	Code string
	NumberOfLevels string
}

type FetchApprovedLecturersInDepartmentResponse struct {
	LecturerId uuid.UUID
	LecturerFirstName string
	LecturerLastName string
	LecturerMiddleName string
	LecturerEmail string
	LecturerProfilePic string
	WaitId uuid.UUID
	AdditionalMessage string
	Approved bool
}

type CreateVenueDto struct {
	VenueName string `json:"venueName" validate:"required"`
	VenueLongitude float64 `json:"venueLongitude" validate:"omitempty"`
	VenueLatitude float64 `json:"venueLatitude" validate:"omitempty"`
	Location string `json:"location" validate:"omitempty"`
	VenueImage string `json:"venueImage" validate:"omitempty"`
	Capacity int32 `json:"capacity" validate:"required"`
	UniversityId uuid.UUID `json:"universityId" validate:"required"`
	VenueType string `json:"venueType" validate:"required"`
	TypeId uuid.UUID `json:"typeId" validate:"required"`
	UnavailabilityDay string `json:"unavailabilityDay" validate:"required"`
	UnavailabilityReason string `json:"unavailabilityReason" validate:"omitempty"`
	UnavailabilityStartTime time.Time `json:"unavailabilityStartTime" validate:"required"`
	UnavailabilityEndTime time.Time `json:"unavailabilityEndtime" validate:"required"`
}

type VenueUnavailability struct {
	Reason string
	Day string
	StartTime time.Time
	EndTime time.Time
}
