package dto

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)


type UpdateAdminDto struct{
	MiddleName string `json:"middleName" validate:"omitempty"`
	PhoneNumber string `json:"phoneNumber" validate:"omitempty,max=15"`
	StaffCard string `json:"staffCard" validate:"omitempty"`
	AdminNumber string `json:"adminNumber" validate:"omitempty"`
	UniversityId string `json:"universityId" validate:"required"`
	AdminId string `json:"adminId" validate:"required"`
}

type CreateUniversityDto struct {
	UniName string `json:"uniName" validate:"required"`
	UniLogo string `json:"uniLogo" validate:"required"`
	UniAbbr string `json:"uniAbbr" validate:"omitempty"`
	UniEmail string `json:"uniEmail" validate:"required"`
	UniWebsite string `json:"uniWebsite" validate:"omitempty"`
	UniPhoneNumber string `json:"uniPhoneNumber" validate:"required"`
	UniversityAddr string `json:"uniAddr" validate:"required"`
	CurrentSession string `json:"currentSession" validate:"required"`
}

type CreateFacultyDto struct {
	FacultyName string `json:"facultyName" validate:"required"`
	FacultyCode string `json:"facultyCode" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
	LecturerId string `json:"lecturerId" validate:"required"`
}
type CreateFacultyResponse struct{
	FacultyID uuid.UUID
	FacultyName string
	FacultyCode sql.NullString
	UniversityID uuid.UUID
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	DeanId uuid.UUID
}


type CreateFacultyDtoResponse struct {
	FacultyName  string    `json:"facultyName" validate:"required"`
	FacultyCode  string    `json:"facultyCode" validate:"required"`
	UniversityId string    `json:"universityId" validate:"required"`
	StartDate    time.Time `json:"startDate" validate:"required"`
	EndDate      time.Time `json:"endDate"` // optional, if not yet ended
}

type CreateDepartmentDto struct {
	DepartmentName string `json:"departmentName" validate:"required"`
	DepartmentCode string `json:"departmentCode" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
	NumberOfLevels int `json:"numberOfLevels" validate:"required"`
}

type CreateDepartmentResponse struct{
	FacultyID uuid.UUID
	DepartmentName string
	DepartmentCode sql.NullString
	UniversityID uuid.UUID
	DepartmentID uuid.UUID
	HodId uuid.UUID
	NumberOfLevels int
}

type CreateDepartmentDtoResponse struct {
	DepartmentName string `json:"departmentName" validate:"required"`
	DepartmentCode string `json:"departmentCode" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
	NumberOfLevels int `json:"numberOfLevels" validate:"required"`
	StartDate    time.Time `json:"startDate" validate:"required"`
	EndDate      time.Time `json:"endDate"` 
}
type RequestDeanConfirmationDto struct {
	PotentialFaculty string `json:"potentialFaculty" validate:"required"`
	AdditionalMessage string `json:"additionalMessage" validate:"omitempty"`
	UniversityId string `json:"universityId" validate:"required"`
}

type RequestHodConfirmationDto struct {
	PotentialDepartment string `json:"potentialDepartment" validate:"required"`
	AdditionalMessage string `json:"additionalMessage" validate:"omitempty"`
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
}

type PendingHodDto struct {
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
}

type PendingLecturerDto struct {
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
	DepartmentId string `json:"departmentId" validate:"required"`
}


type RequestLecturerConfirmationDto struct {
	AdditionalMessage string `json:"addtionalMessage" validate:"omitempty"`
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
	DepartmentId string `json:"departmentId" validate:"required"`
}


type CreateDeanDto struct {
	LecturerId string `json:"lecturerId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
	StartDate string `json:"startDate" validate:"required"`
	EndDate string `json:"endDate" validate:"required"`
}
type CreateHodDto struct {
	LecturerId string `json:"lecturerId" validate:"required"`
	DepartmentId string `json:"facultyId" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
	StartDate string `json:"startDate" validate:"required"`
	EndDate string `json:"endDate" validate:"required"`
}


type DeanWaitingList struct{
	WaitID uuid.UUID
	LecturerID uuid.UUID
	PotentialFaculty string
	UniversityID string
	AdditionalMessage string
	Approved bool
}
type CreateLecturerUnavailability struct {
	Unavailability []LecturerUnavailability
}

type LecturerUnavailability struct{
	Reason string `json:"unavailabilityReason" validate:"omitempty"`
	Day string `json:"unavailabilityDay" validate:"required"`
	StartTime time.Time `json:"unavailabilityStartTime" validate:"required"`
	EndTime time.Time `json:"unavailabilityEndtime" validate:"required"`
}
// type CreateUniversityParams struct {
//     UniversityName string
//     UniversityLogo sql.NullString
//     UniversityAbbr sql.NullString
//     Email          string
//     Website        sql.NullString
//     PhoneNumber    string
//     UniversityAddr sql.NullString
//     CurrentSession sql.NullString
// }