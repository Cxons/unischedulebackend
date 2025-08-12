package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Cxons/unischedulebackend/internal/registration/dto"
	"github.com/Cxons/unischedulebackend/internal/registration/repository"
	"github.com/Cxons/unischedulebackend/internal/registration/service"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	"github.com/Cxons/unischedulebackend/pkg/validator"
)




type RegHandlerInterface interface{
	UpdateAdmin(res http.ResponseWriter , req *http.Request)
	CreateUniversity(res http.ResponseWriter, req *http.Request)
}

var ctx  = context.Background()

type RegHandler struct {
	regService service.RegService
}


func NewRegPackage(logger *slog.Logger,db *sql.DB )*RegHandler{
	query := sqlc.New(db)

	//initializes queries
	studentQueries := queries.NewStudentQueries(query)
	lecturerQueries := queries.NewLecturerQueries(query)
	adminQueries := queries.NewAdminQueries(query)
	uniQueries := queries.NewUniQueries(query)
	deanQueries := queries.NewDeanQueries(query)
	hodQueries := queries.NewHodQueries(query)

	// initializes repository
	repo := repository.NewRegRepository(adminQueries,studentQueries,lecturerQueries,uniQueries,deanQueries,hodQueries)

	// initializes service
	service := service.NewRegService(repo,logger)

	// initializes handler
	handler := NewRegHandler(service)

	return handler
}


func NewRegHandler(service service.RegService)*RegHandler{
	return &RegHandler{
		regService: service,
	}
}


func (rh *RegHandler) UpdateAdmin(res http.ResponseWriter, req *http.Request){
	var body dto.UpdateAdminDto

	if err:= json.NewDecoder(req.Body).Decode(&body); err!=nil{
		http.Error(res,"Invalid Request Body",status.BadRequest.Code)
		return
	}
	if err := validator.ValidateStruct(body); err!= nil{
		http.Error(res,"Validation Error: " + err.Error(),status.BadRequest.Code)
		return
	}
	resp,errMsg,err := rh.regService.UpdateAdmin(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (rh *RegHandler) CreateUniversity(res http.ResponseWriter, req *http.Request){
	var body dto.CreateUniversityDto

	if err:= json.NewDecoder(req.Body).Decode(&body); err!=nil{
		http.Error(res,"Invalid Request Body",status.BadRequest.Code)
		return
	}
	if err := validator.ValidateStruct(body); err!= nil{
		http.Error(res,"Validation Error: " + err.Error(),status.BadRequest.Code)
		return
	}
	resp,errMsg,err := rh.regService.CreateUniversity(ctx,body)

	// if there is no error then set the cookie
	if err == nil{
		cookie := &http.Cookie{
			Name: "university_id",
			Value: resp.Data.(sqlc.University).UniversityID.String(),
			Path: "/admin",
			HttpOnly: true,
			Secure: utils.IsSecure(),
			SameSite: http.SameSiteNoneMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
		http.SetCookie(res,cookie)
	}
	
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) RetrievePendingDeans(res http.ResponseWriter,req *http.Request){
	cookie,err := req.Cookie("university_id")
	if err != nil{
		slog.Error("Error retrieving university id","err:",err)
		http.Error(res,"Error retrieving university id",status.InternalServerError.Code)
		return
	}
	uni_id := cookie.Value
	resp,errMsg,err := rh.regService.RetrievePendingDeans(ctx,uni_id)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


// func (rh *regHandler)