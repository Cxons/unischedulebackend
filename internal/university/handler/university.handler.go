package handler

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/Cxons/unischedulebackend/internal/university/dto"
	"github.com/Cxons/unischedulebackend/internal/university/repository"
	"github.com/Cxons/unischedulebackend/internal/university/service"
)


var ctx  = context.Background()


type UniversityHandlerInterface interface{

}


type UniversityHandler struct{
	service service.UniService
}


func NewUniversityPackage(logger *slog.Logger, db *sql.DB) *UniversityHandler{
	query := sqlc.New(db)
	store := sqlc.NewStore(db)
	//initializes queries
	studentQueries := queries.NewStudentQueries(query)
	lecturerQueries := queries.NewLecturerQueries(query)
	adminQueries := queries.NewAdminQueries(query)
	uniQueries := queries.NewUniQueries(query)
	deanQueries := queries.NewDeanQueries(query)
	hodQueries := queries.NewHodQueries(query)
	facQueries := queries.NewFacQueries(query)
	deptQueries := queries.NewDeptQueries(query)
	cohortQueries :=queries.NewCohortQueries(query)
	venueQueries := queries.NewVenueQueries(query)

	repo := repository.NewUniRepository(adminQueries,studentQueries,lecturerQueries,uniQueries,deanQueries,hodQueries,facQueries,deptQueries,cohortQueries,venueQueries,store)

	service := service.NewUniService(repo,logger)

	handler :=  NewUniversityHandler(service)

	return handler
}


func NewUniversityHandler(service service.UniService)*UniversityHandler{
	return &UniversityHandler{
		service: service,
	}
}



func (uh *UniversityHandler) RetrieveAllUniversities(res http.ResponseWriter,req *http.Request){
	resp,errMsg,err := uh.service.RetrieveAllUniversities(ctx)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (uh *UniversityHandler) RetrieveAllFaculties(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	uniId := queryParams.Get("uniId")
	resp,errMsg,err := uh.service.RetrieveAllFaculties(ctx,utils.StringToUUID(uniId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (uh *UniversityHandler) RetrieveAllDepartments(res http.ResponseWriter,req *http.Request){
	queryParams := req.URL.Query()
	uniId := queryParams.Get("uniId")
	facId := queryParams.Get("facId")
	resp,errMsg,err := uh.service.RetrieveAllDepartments(ctx,dto.RetrieveAllDepartmentsDto{
		UniversityId: uniId,
		FacultyId: facId,
	})
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (uh *UniversityHandler) RetrieveDepartmentLecturers(res http.ResponseWriter, req *http.Request){
	slog.Info("i was reached")
	queryParams := req.URL.Query()
	deptId := queryParams.Get("deptId")
	resp,errMsg,err := uh.service.FetchApprovedLecturersInDepartment(ctx,deptId)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (uh *UniversityHandler) CreateVenue(res http.ResponseWriter, req *http.Request){
	var body dto.CreateVenueDto
	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := uh.service.CreateVenue(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (uh *UniversityHandler) FetchCohortsForADepartment(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	cohortId := queryParams.Get("deptId")
	resp,errMsg,err := uh.service.FetchCohortsForADepartment(ctx,utils.StringToUUID(cohortId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (uh *UniversityHandler) RetrieveAllVenues(res http.ResponseWriter,req *http.Request){
	queryParams := req.URL.Query()
	uniId := queryParams.Get("uniId")
	resp,errMsg,err := uh.service.RetrieveAllVenues(ctx,utils.StringToUUID(uniId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (uh *UniversityHandler) FetchAllDepartmentsForAUni(res http.ResponseWriter,req *http.Request){
	queryParams := req.URL.Query()
	uniId := queryParams.Get("uniId")
	resp,errMsg,err := uh.service.FetchAllDepartmentsForAUni(ctx,utils.StringToUUID(uniId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}