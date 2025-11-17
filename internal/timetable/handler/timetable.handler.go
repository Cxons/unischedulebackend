package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/shared/constants"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/Cxons/unischedulebackend/internal/timetable/dto"
	"github.com/Cxons/unischedulebackend/internal/timetable/repository"
	"github.com/Cxons/unischedulebackend/internal/timetable/service"
	"github.com/Cxons/unischedulebackend/pkg/auth/jwt"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
)


var ctx = context.Background()

type TimetableHandler struct{
	TimeTableService service.TimeTableService
}


func NewTimetablePackage(logger *slog.Logger, db *sql.DB) *TimetableHandler{
	query := sqlc.New(db)
	venueQueries := queries.NewVenueQueries(query)
	lecturerQueries := queries.NewLecturerQueries(query)
	cohortQueries := queries.NewCohortQueries(query)
	courseQueries := queries.NewCoursesQueries(query)
	timetableQueries := queries.NewTimeTableQueries(query)
	store := sqlc.NewStore(db)

	repo := repository.NewtimeTableRepository(venueQueries,lecturerQueries,cohortQueries,courseQueries,timetableQueries,store)
	service := service.NewTimetableService(repo,logger)

	handler := NewTimetableHandler(service)
	return handler

}



func NewTimetableHandler(service service.TimeTableService)*TimetableHandler{
	return &TimetableHandler{
		TimeTableService: service,
	}
}



func (tth *TimetableHandler) CreateATimeTable(res http.ResponseWriter, req *http.Request){
	var body dto.CreateATimeTableDto
	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := tth.TimeTableService.CreateATimeTable(ctx,body.StartTime,body.EndTime,body.UniversityId)
	slog.Info("resp","val",resp)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (tth *TimetableHandler) FetchTimetableForCohort(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	cohortId := queryParams.Get("cohortId")
	uniId := queryParams.Get("uniId")
	slog.Info("cohortid","val",cohortId)
	slog.Info("universityId","val",uniId)
	resp,errMsg,err := tth.TimeTableService.RetrieveTimetableForACohort(ctx,utils.StringToUUID(cohortId),utils.StringToUUID(uniId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (tth *TimetableHandler) FetchTimetableForAStudent(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	uniId := queryParams.Get("uniId")
	var studentId string
	claims := req.Context().Value(constants.UserInfoKey)
	if claims != nil{
		studentId = claims.(*jwt.CustomClaims).User_id
	}else{
		res.Header().Set("Content-Type","application/json")
		res.WriteHeader(status.InternalServerError.Code)
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message":"error validating student authenticity",
			"error": errors.New("error validating student authenticity"),
		})
	}
	resp,errMsg,err := tth.TimeTableService.RetrieveTimetableForAStudent(ctx,utils.StringToUUID(studentId),utils.StringToUUID(uniId))
	utils.HandleAuthResponse(resp,err,errMsg,res)
}