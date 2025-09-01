package dto




type CreateCourseDto struct{
	CourseCode string `json:"courseCode" validate:"required"`
	CourseTitle string `json:"courseTitle" validate:"required"`
	CourseCreditUnit int `json:"courseCreditUnit" validate:"required"`
	CourseDuration int `json:"courseDuration" validate:"required"`
	DepartmentId string `json:"departmentId" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
	LecturerId string `json:"universtyId" validate:"omitempty"`
	SessionsPerWeek int `json:"sessionsPerWeek" validate:"required"`
	Level int `json:"level" validate:"required"`
	Semester string `json:"semester" validate:"required"`
}

type RetrieveCoursesForDeptDto struct {
	DepartmentId string `json:"departmentId" validate:"required"`
	UniversityId string `json:"universityId" validate:"required"`
}


type UpdateCourseDto struct{
	CourseCode string `json:"courseCode" validate:"required"`
	CourseTitle string `json:"courseTitle" validate:"required"`
	CourseCreditUnit int `json:"courseCreditUnit" validate:"required"`
	CourseDuration int `json:"courseDuration" validate:"required"`
	LecturerId string `json:"universtyId" validate:"omitempty"`
	SessionsPerWeek int `json:"sessionsPerWeek" validate:"required"`
	Level int `json:"level" validate:"required"`
	Semester string `json:"semester" validate:"required"`
	CourseId string `json:"courseId" validate:"required"`
}

type SetStudentCourseDto struct {
	CourseId string `json:"courseId" validate:"required"`
	StudentId string `json:"studentId" validate:"required"`
}


type SetCourseLecturersDto struct {
	CourseId string `json:"courseId" validate:"required"`
	LecturerId string `json:"lecturerIc" validate:"required"`
}

type UpdateCourseLecturersDto struct {
	LecturerId string `json:"lecturerId" validate:"required"`
	CourseId string `json:"courseId" validate:"required"`
	LecturerId2 string `json:"lecturerId2" validate:"required"`
}