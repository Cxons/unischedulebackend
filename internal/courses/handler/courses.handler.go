package courses

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/courses/dto"
	"github.com/Cxons/unischedulebackend/internal/courses/repository"
	"github.com/Cxons/unischedulebackend/internal/courses/service"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
)

var ctx  = context.Background()


type CourseHandlerInterface interface {

}


type CourseHandler struct {
	CourseService service.CourseService
}

func NewCourseHandler(courseService service.CourseService)*CourseHandler{
	return &CourseHandler{
		CourseService: courseService,
	}
}


func NewCoursePackage(logger *slog.Logger,db *sql.DB) *CourseHandler{
	query := sqlc.New(db)

	// initializing queries
	courseQueries := queries.NewCoursesQueries(query)


	// initialize store for transactions
	store := sqlc.NewStore(db)


	// initializing repository
	courseRepository := repository.NewCourseRepository(courseQueries,store)


	// initialzing service
	courseService := service.NewCourseService(courseRepository,logger)


	// initializing handler
	courseHandler := NewCourseHandler(courseService)


	return courseHandler

}



func (ch *CourseHandler) CreateCourse( res http.ResponseWriter,req *http.Request){
	var body dto.CreateCourseDto
	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := ch.CourseService.CreateCourse(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (ch *CourseHandler) UpdateCourse(res http.ResponseWriter,req *http.Request){
	var body dto.UpdateCourseDto
	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := ch.CourseService.UpdateCourse(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (ch *CourseHandler) SetStudentCourses(res http.ResponseWriter,req *http.Request){
	var body []dto.SetStudentCourseDto
	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := ch.CourseService.SetStudentCourses(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (ch *CourseHandler) RetrieveCoursesForADepartment(res http.ResponseWriter,req *http.Request){
	var body dto.RetrieveCoursesForDeptDto
	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := ch.CourseService.RetrieveCoursesForADepartment(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (ch *CourseHandler) SetCoursePossibleVenues(res http.ResponseWriter,req *http.Request){
	var body dto.SetCoursePossibleVenuesDto
	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := ch.CourseService.SetCoursePossibleVenues(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (ch *CourseHandler) DeleteCoursePossibleVenue(res http.ResponseWriter,req *http.Request){
	queryParams := req.URL.Query()
	courseId := queryParams.Get("courseId")
	venueId := queryParams.Get("venueId")
	resp,errMsg,err := ch.CourseService.DeleteCoursePossibleVenue(ctx,sqlc.DeleteCoursePossibleVenueParams{
		CourseID: utils.StringToUUID(courseId),
		VenueID: utils.StringToUUID(venueId),
	})
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (ch *CourseHandler) FetchCoursePossibleVenues(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	courseId := queryParams.Get("courseId")
	resp,errMsg,err := ch.CourseService.FetchCoursePossibleVenues(ctx,utils.StringToUUID(courseId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (ch *CourseHandler) RetrieveCoursesForACohort(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	cohortId := queryParams.Get("cohortId")
	resp,errMsg,err := ch.CourseService.RetrieveCoursesForACohort(ctx,utils.StringToUUID(cohortId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (ch *CourseHandler) SetCoursesForACohort(res http.ResponseWriter, req *http.Request){
	var body dto.SetCohortCoursesDto
	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := ch.CourseService.SetCoursesForACohort(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (ch *CourseHandler) FetchAllCourses(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	universityId := queryParams.Get("uniId")
	resp,errMsg,err := ch.CourseService.RetrieveAllCourses(ctx,utils.StringToUUID(universityId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}
// func (ch *CourseHandler) DeleteCourse(req *http.Request, res http.ResponseWriter){
	
// }