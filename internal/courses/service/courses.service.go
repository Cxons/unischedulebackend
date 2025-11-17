package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Cxons/unischedulebackend/internal/courses/dto"
	"github.com/Cxons/unischedulebackend/internal/courses/repository"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	sharedDto "github.com/Cxons/unischedulebackend/internal/shared/dto"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	"github.com/google/uuid"
)



type courseRepository repository.CourseRepository
type createCourseDto = dto.CreateCourseDto
type retrieveCourseForDeptDto = dto.RetrieveCoursesForDeptDto
type updateCourseDto = dto.UpdateCourseDto
type setStudentCourseDto = dto.SetStudentCourseDto
type CourseResponse = sharedDto.ResponseDto
type SetCourseLecturersDto = dto.SetCourseLecturersDto
type UpdateCourseLecturersDto = dto.UpdateCourseLecturersDto




type CourseService interface{
	CreateCourse(ctx context.Context,courseInfo createCourseDto)(CourseResponse,string,error)
	UpdateCourse(ctx context.Context,courseInfo updateCourseDto)(CourseResponse,string,error)
	SetStudentCourses(ctx context.Context,studCourseParam []setStudentCourseDto)(CourseResponse,string,error)
	RetrieveCoursesForADepartment(ctx context.Context, courseParam retrieveCourseForDeptDto)(CourseResponse,string,error)
	DeleteCourse(ctx context.Context,courseId string)(CourseResponse,string,error)
	SetCourseLecturers(ctx context.Context,param SetCourseLecturersDto)(CourseResponse,string,error)
	UpdateCourseLecturers(ctx context.Context, param UpdateCourseLecturersDto)(CourseResponse,string,error)
	SetCoursePossibleVenues(ctx context.Context,courseVenueData dto.SetCoursePossibleVenuesDto)(CourseResponse,string,error)
	DeleteCoursePossibleVenue(ctx context.Context, courseVenueParam sqlc.DeleteCoursePossibleVenueParams)(CourseResponse,string,error)
	FetchCoursePossibleVenues(ctx context.Context,courseId uuid.UUID)(CourseResponse,string,error)
	RetrieveCoursesForACohort(ctx context.Context,cohortId uuid.UUID)(CourseResponse,string,error)
	SetCoursesForACohort(ctx context.Context,params dto.SetCohortCoursesDto)(CourseResponse,string,error)
	RetrieveAllCourses(ctx context.Context, uniId uuid.UUID)(CourseResponse,string,error)
}

type courseService struct {
	repo courseRepository
	logger *slog.Logger
}



func NewCourseService(repo courseRepository,logger *slog.Logger)*courseService{
	return &courseService{
		repo: repo,
		logger: logger,
	}
}



func(cs *courseService) CreateCourse(ctx context.Context,courseInfo createCourseDto)(CourseResponse,string,error){
	course := sqlc.CreateCourseParams{
		CourseCode: courseInfo.CourseCode,
		CourseTitle: courseInfo.CourseTitle,
		CourseDuration:int32(courseInfo.CourseDuration),
		CourseCreditUnit: int32(courseInfo.CourseCreditUnit),
		DepartmentID: utils.StringToUUID(courseInfo.DepartmentId),
		UniversityID: utils.StringToUUID(courseInfo.UniversityId),
		LecturerID: utils.StringToNullUUID(courseInfo.LecturerId),
		SessionsPerWeek: int32(courseInfo.SessionsPerWeek),
		Level: int32(courseInfo.Level),
		Semester: courseInfo.Semester,
	}
	_,err := cs.repo.CreateCourse(ctx,course)

	if err != nil{
		cs.logger.Error("error creating course","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Course created successfully",
		Data: nil,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}


func (cs *courseService) UpdateCourse(ctx context.Context,courseInfo updateCourseDto)(CourseResponse,string,error){
	course := sqlc.UpdateCourseParams{
		CourseCode: courseInfo.CourseCode,
		CourseTitle: courseInfo.CourseTitle,
		CourseDuration: int32(courseInfo.CourseDuration),
		CourseCreditUnit: int32(courseInfo.CourseCreditUnit),
		SessionsPerWeek: int32(courseInfo.SessionsPerWeek),
		LecturerID: utils.StringToNullUUID(courseInfo.LecturerId),
		Level: int32(courseInfo.Level),
		Semester: courseInfo.Semester,
		CourseID: utils.StringToUUID(courseInfo.CourseId),
	}

	_,err := cs.repo.UpdateCourse(ctx,course)

		if err != nil{
		cs.logger.Error("error updating course","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Course created successfully",
		Data: nil,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil

}

func (cs *courseService) RetrieveCoursesForADepartment(ctx context.Context, courseParam retrieveCourseForDeptDto)(CourseResponse,string,error){
	course := sqlc.RetrieveCoursesForADepartmentParams{
		DepartmentID: utils.StringToUUID(courseParam.DepartmentId),
		UniversityID: utils.StringToUUID(courseParam.UniversityId),
	}
	deptExists,courses,err := cs.repo.RetrieveCoursesForADepartment(ctx,course)
	if err != nil{
		if !deptExists{
			return CourseResponse{},status.NotFound.Message,errors.New("no courses created for this department yet")
		}
		cs.logger.Error("error retrieving courses for the department","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "The courses for this department",
		Data: courses,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}


func (cs *courseService) SetStudentCourses(ctx context.Context,studCourseParam []setStudentCourseDto)(CourseResponse,string,error){
	var studCourses []sqlc.SetStudentCourseParams

	for i := 0; i < len(studCourseParam); i++ {
		studCourse := sqlc.SetStudentCourseParams{
			CourseID: utils.StringToUUID(studCourseParam[i].CourseId),
			StudentID: utils.StringToUUID(studCourseParam[i].StudentId),
	}
		studCourses = append(studCourses, studCourse)
	}


	err := cs.repo.SetStudentCourses(ctx,studCourses)

	if err != nil{
		cs.logger.Error("error setting student courses","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Student courses registered",
		Data: nil,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}


func (cs *courseService) DeleteCourse(ctx context.Context,courseId string)(CourseResponse,string,error){
	if err := cs.repo.DeleteCourse(ctx,utils.StringToUUID(courseId)); err != nil{
		cs.logger.Error("error deleting course","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Course deleted successfully",
		Data: nil,
		StatusCode: status.NoContent.Code,
		StatusCodeMessage: status.NoContent.Message,
	},status.NoContent.Message,nil
	
}


func (cs *courseService) SetCourseLecturers(ctx context.Context,param SetCourseLecturersDto)(CourseResponse,string,error){
	courseLecturer := sqlc.SetCourseLecturersParams{
		LecturerID: utils.StringToUUID(param.LecturerId),
		CourseID: utils.StringToUUID(param.CourseId),
	}
	_,err := cs.repo.SetCourseLecturers(ctx,courseLecturer)


	if err != nil{
		cs.logger.Error("error setting course lecturers","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Lecturer successfully assigned course",
		Data: nil,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}


func (cs *courseService) UpdateCourseLecturers(ctx context.Context, param UpdateCourseLecturersDto)(CourseResponse,string,error){
	courseLecturer := sqlc.UpdateCourseLecturersParams{
		LecturerID: utils.StringToUUID(param.LecturerId),
		LecturerID_2: utils.StringToUUID(param.LecturerId2),
	}
	_,err := cs.repo.UpdateCourseLecturers(ctx,courseLecturer)

	if err != nil{
		cs.logger.Error("error updating course lecturers","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Updated course lecturers successfully",
		Data: nil,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (cs *courseService) SetCoursePossibleVenues(ctx context.Context,courseVenueData dto.SetCoursePossibleVenuesDto)(CourseResponse,string,error){
	actualCourseVenueData := make([]sqlc.SetCoursePossibleVenueParams,0)
	for _,val := range courseVenueData.Venues{
		actualCourseVenueData = append(actualCourseVenueData,sqlc.SetCoursePossibleVenueParams{
			CourseID: utils.StringToUUID(courseVenueData.CourseId),
			VenueID: val,
			UniversityID: utils.StringToNullUUID(courseVenueData.UniversityId),
	} ) 
	}
	err := cs.repo.SetCoursePossibleVenues(ctx,actualCourseVenueData)
	if err != nil{
		cs.logger.Error("error setting course possible venues","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Course possible venues set successfully",
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}

func (cs *courseService) DeleteCoursePossibleVenue(ctx context.Context, courseVenueParam sqlc.DeleteCoursePossibleVenueParams)(CourseResponse,string,error){
	err := cs.repo.DeleteCoursePossibleVenue(ctx,courseVenueParam)
	if err != nil{
		cs.logger.Error("error deleting course possible venue","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Course possible venue deleted successfully",
		StatusCode: status.NoContent.Code,
		StatusCodeMessage: status.NoContent.Message,
	},status.NoContent.Message,nil
}

func (cs *courseService) FetchCoursePossibleVenues(ctx context.Context,courseId uuid.UUID)(CourseResponse,string,error){
	data,err := cs.repo.FetchCoursePossibleVenues(ctx,courseId)
	if err != nil{
		cs.logger.Error("error fetching course possible venues","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Course and their possible venues",
		Data: data,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (cs *courseService) RetrieveCoursesForACohort(ctx context.Context,cohortId uuid.UUID)(CourseResponse,string,error){
	data,err := cs.repo.RetrieveCoursesForACohort(ctx,cohortId)
	if err != nil{
		cs.logger.Error("error retrieving courses for the cohort","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Courses by the cohorts",
		Data: data,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (cs *courseService) SetCoursesForACohort(ctx context.Context,params dto.SetCohortCoursesDto)(CourseResponse,string,error){
	err := cs.repo.SetCoursesForACohort(ctx,utils.StringToUUID(params.UniversityId),utils.StringToUUID(params.CohortId),params.Courses)
	if err != nil{
		cs.logger.Error("error setting courses for a cohort","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "Courses set successfully",
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}

func (cs *courseService) RetrieveAllCourses(ctx context.Context, uniId uuid.UUID)(CourseResponse,string,error){
	data,err := cs.repo.RetrieveAllCourses(ctx,uniId)
	if err != nil{
		cs.logger.Error("error retrieving all courses","err:",err)
		return CourseResponse{},status.InternalServerError.Message,err
	}
	return CourseResponse{
		Message: "The courses",
		Data: data,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}
