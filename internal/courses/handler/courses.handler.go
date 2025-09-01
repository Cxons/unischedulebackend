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



func (ch *CourseHandler) CreateCourse(req *http.Request, res http.ResponseWriter){
	var body dto.CreateCourseDto
	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := ch.CourseService.CreateCourse(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (ch *CourseHandler) UpdateCourse(req *http.Request, res http.ResponseWriter){
	var body dto.UpdateCourseDto
	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := ch.CourseService.UpdateCourse(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (ch *CourseHandler) SetStudentCourses(req *http.Request, res http.ResponseWriter){
	var body []dto.SetStudentCourseDto
	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := ch.CourseService.SetStudentCourses(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (ch *CourseHandler) RetrieveCoursesForADepartment(req *http.Request, res http.ResponseWriter){
	var body dto.RetrieveCoursesForDeptDto
	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := ch.CourseService.RetrieveCoursesForADepartment(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (ch *CourseHandler) DeleteCourse(req *http.Request, res http.ResponseWriter){
	
}