package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	regDto "github.com/Cxons/unischedulebackend/internal/registration/dto"
	"github.com/Cxons/unischedulebackend/internal/registration/repository"
	"github.com/Cxons/unischedulebackend/internal/registration/service"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
)



type cookieData struct{
		UniversityId string
		FacultyId string
		DepartmentId string
}

type RegHandlerInterface interface{
	UpdateAdmin(res http.ResponseWriter , req *http.Request)
	CreateUniversity(res http.ResponseWriter, req *http.Request)
	RetrievePendingDeans(res http.ResponseWriter,req *http.Request)
	CheckDeanConfirmation(res http.ResponseWriter,req *http.Request)
	RequestDeanConfirmation(res http.ResponseWriter, req *http.Request)
	ApproveDean(res http.ResponseWriter, req *http.Request)
	CreateFaculty(res http.ResponseWriter,req *http.Request)
	RetrievePendingHods(res http.ResponseWriter, req *http.Request)
	RequestHodConfirmation(res http.ResponseWriter, req *http.Request)
	ApproveHod(res http.ResponseWriter, req *http.Request)
	CheckHodConfirmation(res http.ResponseWriter,req *http.Request )
	CreateDeparment(res http.ResponseWriter, req *http.Request)
	RetrievePendingLecturers(res http.ResponseWriter, req *http.Request)
	RequestLecturerConfirmation(res http.ResponseWriter, req *http.Request)
	ApproveLecturer(res http.ResponseWriter, req *http.Request)
	CheckLecturerConfirmation(res http.ResponseWriter, req *http.Request)
}

var ctx  = context.Background()

type RegHandler struct {
	RegService service.RegService
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
	facQueries := queries.NewFacQueries(query)
	deptQueries := queries.NewDeptQueries(query)
	store := sqlc.NewStore(db)

	// initializes repository
	repo := repository.NewRegRepository(adminQueries,studentQueries,lecturerQueries,uniQueries,deanQueries,hodQueries,facQueries,deptQueries,store)

	// initializes service
	service := service.NewRegService(repo,logger)

	// initializes handler
	handler := NewRegHandler(service)

	return handler
}


func NewRegHandler(service service.RegService)*RegHandler{
	return &RegHandler{
		RegService: service,
	}
}


func (rh *RegHandler) UpdateAdmin(res http.ResponseWriter, req *http.Request){
	var body regDto.UpdateAdminDto

	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.UpdateAdmin(ctx,body)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (rh *RegHandler) CreateUniversity(res http.ResponseWriter, req *http.Request){
	var body regDto.CreateUniversityDto

	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.CreateUniversity(ctx,body)

	// if there is no error then set the cookie
	if err == nil{
		cookie := &http.Cookie{
			Name: "university_id",
			Value: resp.Data.(sqlc.University).UniversityID.String(),
			Path: "/",
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
	resp,errMsg,err := rh.RegService.RetrievePendingDeans(ctx,uni_id)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (rh *RegHandler) RequestDeanConfirmation(res http.ResponseWriter, req *http.Request){
	var body regDto.RequestDeanConfirmationDto

	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.RequestDeanConfirmation(req.Context(),body)
	if err == nil{
		cookie := &http.Cookie{
			Name: "dean_wait_id",
			Value: resp.Data.(sqlc.DeanWaitingList).WaitID.String(),
			Path: "/",
			HttpOnly: true,
			Secure: utils.IsSecure(),
			SameSite: http.SameSiteNoneMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
	http.SetCookie(res,cookie)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}
}

func (rh *RegHandler) ApproveDean(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()

	waitId := queryParams.Get("wait_Id")
	resp,errMsg,err := rh.RegService.ApproveDean(ctx,waitId)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) CheckDeanConfirmation(res http.ResponseWriter,req *http.Request ){
	cookie,err := req.Cookie("dean_wait_id")
	if err != nil{
		slog.Error("Error retrieving dean wait id","err:",err)
		http.Error(res,"Error retrieving dean wait id",status.InternalServerError.Code)
		return
	}
	wait_id := cookie.Value
	resp,errMsg,err := rh.RegService.CheckDeanConfirmation(ctx,wait_id)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) CreateFaculty(res http.ResponseWriter,req *http.Request){
	var body regDto.CreateFacultyDto

	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.CreateFaculty(ctx,body)

	cookieValue := &cookieData{
		UniversityId: resp.Data.(sqlc.Faculty).UniversityID.String(),
		FacultyId: resp.Data.(sqlc.Faculty).FacultyID.String(),
	}

	cookieJsonData,jsonErr := json.Marshal(cookieValue)
	if jsonErr != nil{
		slog.Error("Error marshalling json data","err:",err)
		http.Error(res,"Problem marshalling json data",status.InternalServerError.Code)
		return
	}
	// if there is no error then set the cookie
	if err != nil{
		cookie := &http.Cookie{
			Name: "faculty_info",
			Value: string(cookieJsonData),
			Path: "/",
			HttpOnly: true,
			Secure: utils.IsSecure(),
			SameSite: http.SameSiteNoneMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
		http.SetCookie(res,cookie)
	}
	
	utils.HandleAuthResponse(resp,err,errMsg,res)
}
func (rh *RegHandler) RetrievePendingHods(res http.ResponseWriter, req *http.Request){
	cookie,err := req.Cookie("faculty_info")
	var cookieValue cookieData

	// handle cookie error
	if err != nil{
		slog.Error("Error retrieving faculty info","err:",err)
		http.Error(res,"Error retrieving faculty info",status.InternalServerError.Code)
		return
	}
	fac_info := cookie.Value

	// coverts cookie value of string to proper struct
	var byte_fac_info = []byte(fac_info)
	if err :=json.Unmarshal(byte_fac_info,&cookieValue); err!= nil{
		slog.Error("Error unmarshaling json data","err:",err)
		http.Error(res,"Error unmarshaling json data",status.InternalServerError.Code)
	}

	resp,errMsg,err := rh.RegService.RetrievePendingHods(ctx,service.PendingHodDto{
		FacultyId: cookieValue.FacultyId,
		UniversityId: cookieValue.UniversityId,
	})
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) RequestHodConfirmation(res http.ResponseWriter, req *http.Request){
	var body regDto.RequestHodConfirmationDto

	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.RequestHodConfirmation(req.Context(),body)
	if err == nil{
		cookie := &http.Cookie{
			Name: "hod_wait_id",
			Value: resp.Data.(sqlc.HodWaitingList).WaitID.String(),
			Path: "/",
			HttpOnly: true,
			Secure: utils.IsSecure(),
			SameSite: http.SameSiteNoneMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
	http.SetCookie(res,cookie)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}
}

func (rh *RegHandler) ApproveHod(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()

	waitId := queryParams.Get("wait_Id")
	resp,errMsg,err := rh.RegService.ApproveHod(ctx,waitId)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) CheckHodConfirmation(res http.ResponseWriter,req *http.Request ){
	cookie,err := req.Cookie("hod_wait_id")
	if err != nil{
		slog.Error("Error retrieving hod wait id","err:",err)
		http.Error(res,"Error retrieving hod wait id",status.InternalServerError.Code)
		return
	}
	wait_id := cookie.Value
	resp,errMsg,err := rh.RegService.CheckHodConfirmation(ctx,wait_id)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (rh *RegHandler) CreateDeparment(res http.ResponseWriter, req *http.Request){
	var body regDto.CreateDepartmentDto

	
	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.CreateDepartment(ctx,body)

	cookieValue := &cookieData{
		UniversityId: resp.Data.(sqlc.Department).UniversityID.String(),
		FacultyId: resp.Data.(sqlc.Department).FacultyID.String(),
		DepartmentId: resp.Data.(sqlc.Department).DepartmentID.String(),
	}

	cookieJsonData,jsonErr := json.Marshal(cookieValue)
	if jsonErr != nil{
		slog.Error("Error marshalling json data","err:",err)
		http.Error(res,"Problem marshalling json data",status.InternalServerError.Code)
		return
	}
	// if there is no error then set the cookie
	if err != nil{
		cookie := &http.Cookie{
			Name: "department_info",
			Value: string(cookieJsonData),
			Path: "/",
			HttpOnly: true,
			Secure: utils.IsSecure(),
			SameSite: http.SameSiteNoneMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
		http.SetCookie(res,cookie)
	}
	
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) RetrievePendingLecturers(res http.ResponseWriter, req *http.Request){
	cookie,err := req.Cookie("department_info")
	var cookieValue cookieData

	// handle cookie error
	if err != nil{
		slog.Error("Error retrieving department info","err:",err)
		http.Error(res,"Error retrieving department info",status.InternalServerError.Code)
		return
	}
	dept_info := cookie.Value

	// coverts cookie value of string to proper struct
	var byte_dept_info = []byte(dept_info)
	if err :=json.Unmarshal(byte_dept_info,&cookieValue); err!= nil{
		slog.Error("Error unmarshaling json data","err:",err)
		http.Error(res,"Error unmarshaling json data",status.InternalServerError.Code)
	}

	resp,errMsg,err := rh.RegService.RetrievePendingLecturers(ctx,service.PendingLecturerDto{
		FacultyId: cookieValue.FacultyId,
		UniversityId: cookieValue.UniversityId,
		DepartmentId: cookieValue.DepartmentId,
	})
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) RequestLecturerConfirmation(res http.ResponseWriter, req *http.Request){
	var body regDto.RequestLecturerConfirmationDto

	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.RequestLecturerConfirmation(req.Context(),body)
	if err == nil{
		cookie := &http.Cookie{
			Name: "lecturer_wait_id",
			Value: resp.Data.(sqlc.LecturerWaitingList).WaitID.String(),
			Path: "/",
			HttpOnly: true,
			Secure: utils.IsSecure(),
			SameSite: http.SameSiteNoneMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
	http.SetCookie(res,cookie)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}
}


func (rh *RegHandler) ApproveLecturer(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()

	waitId := queryParams.Get("wait_Id")
	resp,errMsg,err := rh.RegService.ApproveLecturer(ctx,waitId)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (rh *RegHandler) CheckLecturerConfirmation(res http.ResponseWriter, req *http.Request){
	cookie,err := req.Cookie("lecturer_wait_id")
	if err != nil{
		slog.Error("Error retrieving lecturer wait id","err:",err)
		http.Error(res,"Error retrieving lecturer wait id",status.InternalServerError.Code)
		return
	}
	wait_id := cookie.Value
	resp,errMsg,err := rh.RegService.CheckLecturerConfirmation(ctx,wait_id)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) CreateDean(res http.ResponseWriter, req *http.Request){
	
	var body regDto.CreateDeanDto

	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.CreateDean(ctx,body)

	// cookie expires a month after end date for dean tenure
	expires := resp.Data.(sqlc.CurrentDean).EndDate.Time.AddDate(0,1,0)
	
	if err != nil{
		cookie := &http.Cookie{
			Name: "current_dean_id",
			Value: resp.Data.(sqlc.CurrentDean).DeanID.String(),
			Path: "/",
			HttpOnly: true,
			Secure: utils.IsSecure(),
			SameSite: http.SameSiteNoneMode,
			Expires: expires,
	}
		http.SetCookie(res,cookie)
	}
	
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) CreateHod(res http.ResponseWriter, req *http.Request){
	
	var body regDto.CreateHodDto

	utils.HandleBodyParsing(req,res,body)
	resp,errMsg,err := rh.RegService.CreateHod(ctx,body)

	// cookie expires a month after end date of hod tenure
	expires := resp.Data.(sqlc.CurrentHod).EndDate.Time.AddDate(0,1,0)
	
	if err != nil{
		cookie := &http.Cookie{
			Name: "current_hod_id",
			Value: resp.Data.(sqlc.CurrentHod).HodID.String(),
			Path: "/",
			HttpOnly: true,
			Secure: utils.IsSecure(),
			SameSite: http.SameSiteNoneMode,
			Expires: expires,
	}
		http.SetCookie(res,cookie)
	}
	
	utils.HandleAuthResponse(resp,err,errMsg,res)
}