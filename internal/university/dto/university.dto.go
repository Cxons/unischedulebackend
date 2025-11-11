package dto

type RetrieveAllDepartmentsDto struct{
	FacultyId string `json:"FacultyId" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
}
