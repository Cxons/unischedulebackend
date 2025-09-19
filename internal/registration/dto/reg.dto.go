package dto


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
}

type CreateDepartmentDto struct {
	DepartmentName string `json:"departmentName" validate:"required"`
	DepartmentCode string `json:"departmentCode" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
	NumberOfLevels int `json:"numberOfLevels" validate:"required"`
}

type RequestDeanConfirmationDto struct {
	PotentialFaculty string `json:"potentialFaculty" validate:"required"`
	AdditionalMessage string `json:"addtionalMessage" validate:"omitempty"`
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
	UniversityId string `json:"ungiversityId" validate:"required"`
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