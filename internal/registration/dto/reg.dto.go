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

type RequestDeanConfirmationDto struct {
	LecturerId string `json:"lecturerId" validate:"required"`
	PotentialFaculty string `json:"potentialFaculty" validate:"required"`
	AdditionalMessage string `json:"addtionalMessage" validate:"omitempty"`
	UniversityId string `json:"universityId" validate:"required"`
}

type RequestHodConfirmationDto struct {
	LecturerId string `json:"lecturerId" validate:"required"`
	PotentialDepartment string `json:"potentialDepartment" validate:"required"`
	AdditionalMessage string `json:"additionalMessage" validate:"omitempty"`
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
}


type RequestLecturerConfirmationDto struct {
	LecturerId string `json:"lecturerId" validate:"required"`
	AdditionalMessage string `json:"addtionalMessage" validate:"omitempty"`
	UniversityId string `json:"universityId" validate:"required"`
	FacultyId string `json:"facultyId" validate:"required"`
	DepartmentId string `json:"departmentId" validate:"required"`
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