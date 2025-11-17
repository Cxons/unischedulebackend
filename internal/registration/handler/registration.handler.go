package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	regDto "github.com/Cxons/unischedulebackend/internal/registration/dto"
	"github.com/Cxons/unischedulebackend/internal/registration/repository"
	"github.com/Cxons/unischedulebackend/internal/registration/service"
	"github.com/Cxons/unischedulebackend/internal/shared/constants"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/dto"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/Cxons/unischedulebackend/pkg/auth/jwt"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
)



type cookieData struct{
		UniversityId string
		FacultyId string
		DepartmentId string
		DeanId string
}
type HodCookieData struct{
		UniversityId string
		FacultyId string
		DepartmentId string
		HodId string
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
	FetchDeanWaitDetails(res http.ResponseWriter,req *http.Request)
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

	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := rh.RegService.CreateUniversity(ctx,body)

	// if there is no error then set the cookie
	if err == nil{
		cookie := &http.Cookie{
			Name: "university_id",
			Value: resp.Data.(sqlc.University).UniversityID.String(),
			Path: "/",
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteLaxMode,
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

func (rh *RegHandler) FetchDeanWaitDetails(res http.ResponseWriter,req *http.Request){
	queryParams := req.URL.Query()

	waitId := queryParams.Get("waitId")
	resp,errMsg,err := rh.RegService.FetchDeanWaitDetails(ctx,waitId)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) FetchHodWaitDetails(res http.ResponseWriter,req *http.Request){
	queryParams := req.URL.Query()

	waitId := queryParams.Get("waitId")
	resp,errMsg,err := rh.RegService.FetchHodWaitDetails(ctx,waitId)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) FetchLecturerWaitDetails(res http.ResponseWriter,req *http.Request){
	queryParams := req.URL.Query()

	waitId := queryParams.Get("waitId")
	resp,errMsg,err := rh.RegService.FetchLecturerWaitDetails(ctx,waitId)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}


func (rh *RegHandler) RequestDeanConfirmation(res http.ResponseWriter, req *http.Request){
	var body regDto.RequestDeanConfirmationDto

	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := rh.RegService.RequestDeanConfirmation(req.Context(),body)
	if err == nil{
		cookie := &http.Cookie{
			Name: "dean_wait_id",
			Value: resp.Data.(sqlc.DeanWaitingList).WaitID.String(),
			Path: "/",
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteLaxMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
	http.SetCookie(res,cookie)
}
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) ApproveDean(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()

	waitId := queryParams.Get("wait_Id")
	lecturerId := queryParams.Get("lecturer_id")
	resp,errMsg,err := rh.RegService.ApproveDean(ctx,waitId,lecturerId)
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
	var body regDto.CreateFacultyDtoResponse

	utils.HandleBodyParsing(req,res,&body)
	var lectuererId string
	claims := req.Context().Value(constants.UserInfoKey)
	if claims != nil{
		lectuererId = claims.(*jwt.CustomClaims).User_id
	}else{
		res.Header().Set("Content-Type","application/json")
		res.WriteHeader(status.InternalServerError.Code)
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message":"error validating dean authenticity",
			"error": errors.New("error validating dean authenticity"),
		})
	}
	resp,errMsg,err := rh.RegService.CreateFaculty(ctx,regDto.CreateFacultyDto{
		UniversityId: body.UniversityId,
		FacultyName: body.FacultyName,
		FacultyCode: body.FacultyCode,
		LecturerId: lectuererId,
	},body.StartDate,body.EndDate)
	facultyData,ok := resp.Data.(regDto.CreateFacultyResponse)

	if !ok {
    slog.Error("creating faculty response data is invalid", "resp.Data", resp.Data)

    res.Header().Set("Content-Type", "application/json")

    // If the error corresponds to unauthorized (e.g., status.Unauthorized)
    if errMsg == status.Forbidden.Message {
        res.WriteHeader(http.StatusForbidden)
        json.NewEncoder(res).Encode(map[string]interface{}{
            "message": "Dean not confirmed",
            "error":   "Unauthorized",
        })
        return
    }

    // Fallback for other cases
    res.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(res).Encode(map[string]interface{}{
        "message": "Problem creating faculty",
        "error":   "Dean is probably not confirmed",
    })
    return
}

	cookieValue := &cookieData{
		UniversityId: facultyData.UniversityID.String(),
		FacultyId: facultyData.FacultyID.String(),
		DeanId: facultyData.DeanId.String(),
	}

	cookieJsonData,jsonErr := json.Marshal(cookieValue)
	if jsonErr != nil{
		slog.Error("Error marshalling json data","err:",err)
		http.Error(res,"Problem marshalling json data",status.InternalServerError.Code)
		return
	}

	slog.Info("cookiejsondata","val",cookieJsonData)
	// if there is no error then set the cookie
	if err == nil{
		cookie := &http.Cookie{
			Name: "faculty_info",
			Value: string(cookieJsonData),
			Path: "/",
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteLaxMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
		http.SetCookie(res,cookie)
	}
	deanCookie := &http.Cookie{
			Name: "current_dean_id",
			Value: facultyData.DeanId.String(),
			Path: "/",
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteLaxMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
		http.SetCookie(res,deanCookie)
	
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) RetrievePendingHods(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// Retrieve the cookie
	cookie, err := req.Cookie("faculty_info")
	if err != nil {
		slog.Error("Error retrieving faculty info", "err", err)
		http.Error(res, "Error retrieving faculty info", status.InternalServerError.Code)
		return
	}

	facInfo := cookie.Value
	slog.Info("Raw cookie value", "fac_info", facInfo)

	var cookieValue cookieData

	// Try to unmarshal as JSON first
	err = json.Unmarshal([]byte(facInfo), &cookieValue)
	if err != nil {
		slog.Warn("Cookie not valid JSON, trying fallback parser", "err", err)

		// Try to manually parse the non-JSON format
		// Expected: {UniversityId:xxx,FacultyId:yyy,DepartmentId:zzz,DeanId:aaa}
		facInfo = strings.Trim(facInfo, "{}")
		parts := strings.Split(facInfo, ",")

		for _, part := range parts {
			pair := strings.SplitN(part, ":", 2)
			if len(pair) != 2 {
				continue
			}
			key := strings.TrimSpace(pair[0])
			val := strings.TrimSpace(pair[1])
			switch key {
			case "UniversityId":
				cookieValue.UniversityId = val
			case "FacultyId":
				cookieValue.FacultyId = val
			case "DepartmentId":
				cookieValue.DepartmentId = val
			case "DeanId":
				cookieValue.DeanId = val
			}
		}
	}

	// Log parsed values
	slog.Info("Parsed cookie data",
		"UniversityId", cookieValue.UniversityId,
		"FacultyId", cookieValue.FacultyId,
	)

	if cookieValue.UniversityId == "" || cookieValue.FacultyId == "" {
		http.Error(res, "Missing faculty info in cookie", http.StatusBadRequest)
		return
	}

	// Pass correct data to service
	resp, errMsg, err := rh.RegService.RetrievePendingHods(ctx, service.PendingHodDto{
		FacultyId:    cookieValue.FacultyId,
		UniversityId: cookieValue.UniversityId,
	})

	utils.HandleAuthResponse(resp, err, errMsg, res)
}



func (rh *RegHandler) RequestHodConfirmation(res http.ResponseWriter, req *http.Request){
	var body regDto.RequestHodConfirmationDto

	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := rh.RegService.RequestHodConfirmation(req.Context(),body)
	if err == nil{
		cookie := &http.Cookie{
			Name: "hod_wait_id",
			Value: resp.Data.(sqlc.HodWaitingList).WaitID.String(),
			Path: "/",
			HttpOnly: true,
			Secure:false,
			SameSite: http.SameSiteLaxMode,
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
	var body regDto.CreateDepartmentDtoResponse
	utils.HandleBodyParsing(req,res,&body)
	var lecturerId string
	claims := req.Context().Value(constants.UserInfoKey)
	if claims != nil{
		lecturerId = claims.(*jwt.CustomClaims).User_id
	}else{
		res.Header().Set("Content-Type","application/json")
		res.WriteHeader(status.InternalServerError.Code)
		json.NewEncoder(res).Encode(map[string]interface{}{
			"message":"error validating hod authenticity",
			"error": errors.New("error validating hod authenticity"),
		})
	}
	resp,errMsg,err := rh.RegService.CreateDepartment(ctx,regDto.CreateDepartmentDto{
		DepartmentName: body.DepartmentName,
		DepartmentCode: body.DepartmentCode,
		UniversityId: body.UniversityId,
		FacultyId: body.FacultyId,
		NumberOfLevels: body.NumberOfLevels,
	},utils.StringToUUID(lecturerId),body.StartDate,body.EndDate)
	deptInfo,ok := resp.Data.(regDto.CreateDepartmentResponse)
	slog.Info("the dept info","dept",deptInfo)
	if !ok {
    slog.Error("creating department response data is invalid", "resp.Data", resp.Data)

    res.Header().Set("Content-Type", "application/json")

    // If the error corresponds to unauthorized (e.g., status.Unauthorized)
    if errMsg == status.Forbidden.Message {
        res.WriteHeader(http.StatusForbidden)
        json.NewEncoder(res).Encode(map[string]interface{}{
            "message": "Hod not confirmed",
            "error":   "Unauthorized",
        })
        return
    }
	slog.Error("error","err",errMsg)

    // Fallback for other cases
    res.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(res).Encode(map[string]interface{}{
        "message": "Problem creating department",
        "error":   "Hod is probably not confirmed",
    })
    return
}

	cookieValue := &cookieData{
		UniversityId: deptInfo.UniversityID.String(),
		FacultyId: deptInfo.FacultyID.String(),
		DepartmentId: deptInfo.DepartmentID.String(),
	}

	cookieJsonData,jsonErr := json.Marshal(cookieValue)
	if jsonErr != nil{
		slog.Error("Error marshalling json data","err:",err)
		http.Error(res,"Problem marshalling json data",status.InternalServerError.Code)
		return
	}

	slog.Info("cookiejsondata","val",cookieJsonData)
	// if there is no error then set the cookie
	if err == nil{
		cookie := &http.Cookie{
			Name: "department_info",
			Value: string(cookieJsonData),
			Path: "/",
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteLaxMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
		http.SetCookie(res,cookie)
	}
	hodCookie := &http.Cookie{
			Name: "current_hod_id",
			Value: deptInfo.HodId.String(),
			Path: "/",
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteLaxMode,
			Expires: time.Now().AddDate(10, 0, 0),
	}
		http.SetCookie(res,hodCookie)
	
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) RetrievePendingLecturers(res http.ResponseWriter, req *http.Request){
	cookie,err := req.Cookie("department_info")
	var cookieValue HodCookieData

	// handle cookie error
	if err != nil{
		slog.Error("Error retrieving department info","err:",err)
		http.Error(res,"Error retrieving department info",status.InternalServerError.Code)
		return
	}
	dept_info := cookie.Value

	// Try to unmarshal as JSON first
	err = json.Unmarshal([]byte(dept_info), &cookieValue)
	if err != nil {
		slog.Warn("Cookie not valid JSON, trying fallback parser", "err", err)
		// Try to manually parse the non-JSON format
		// Expected: {UniversityId:xxx,FacultyId:yyy,DepartmentId:zzz,DeanId:aaa}
		dept_info = strings.Trim(dept_info, "{}")
		parts := strings.Split(dept_info, ",")

		for _, part := range parts {
			pair := strings.SplitN(part, ":", 2)
			if len(pair) != 2 {
				continue
			}
			key := strings.TrimSpace(pair[0])
			val := strings.TrimSpace(pair[1])
			switch key {
			case "UniversityId":
				cookieValue.UniversityId = val
			case "FacultyId":
				cookieValue.FacultyId = val
			case "DepartmentId":
				cookieValue.DepartmentId = val
			case "HodId":
				cookieValue.HodId = val
			}
		}
	}

	// Log parsed values
	slog.Info("Parsed cookie data",
		"UniversityId", cookieValue.UniversityId,
		"FacultyId", cookieValue.FacultyId,
	)

	if cookieValue.UniversityId == "" || cookieValue.FacultyId == "" || cookieValue.DepartmentId == "" {
		http.Error(res, "Missing department info in cookie", http.StatusBadRequest)
		return
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

	utils.HandleBodyParsing(req,res,&body)
	resp,errMsg,err := rh.RegService.RequestLecturerConfirmation(req.Context(),body)
	if err == nil{
		cookie := &http.Cookie{
			Name: "lecturer_wait_id",
			Value: resp.Data.(sqlc.LecturerWaitingList).WaitID.String(),
			Path: "/",
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteLaxMode,
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
			Secure: false,
			SameSite: http.SameSiteLaxMode,
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
			Secure: false,
			SameSite: http.SameSiteLaxMode,
			Expires: expires,
	}
		http.SetCookie(res,cookie)
	}
	
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (rh *RegHandler) CreateLecturerUnavailability(res http.ResponseWriter, req *http.Request) {
    var body struct {
        Unavailability []struct {
            Reason    string `json:"unavailabilityReason" validate:"omitempty"`
            Day       string `json:"unavailabilityDay" validate:"required"`
            StartTime string `json:"unavailabilityStartTime" validate:"required"` // Change to string
            EndTime   string `json:"unavailabilityEndtime" validate:"required"`   // Change to string
        } `json:"unavailability"`
    }
    
    if err := utils.HandleBodyParsing(req, res, &body); err != nil {
        return
    }

    var lecturerId string
    claims := req.Context().Value(constants.UserInfoKey)
    if claims != nil {
        lecturerId = claims.(*jwt.CustomClaims).User_id
    } else {
        res.Header().Set("Content-Type", "application/json")
        res.WriteHeader(status.InternalServerError.Code)
        json.NewEncoder(res).Encode(map[string]interface{}{
            "message": "error validating dean authenticity",
            "error":   errors.New("error validating dean authenticity"),
        })
        return
    }

    // Convert to service DTO with proper time parsing
    serviceDTO := regDto.CreateLecturerUnavailability{
        Unavailability: make([]regDto.LecturerUnavailability, 0, len(body.Unavailability)),
    }

    for _, slot := range body.Unavailability {
        // Parse ISO time strings to time.Time
        startTime, err := parseISOTime(slot.StartTime)
        if err != nil {
            utils.HandleAuthResponse(dto.ResponseDto{}, err, "Invalid start time format: "+slot.StartTime, res)
            return
        }

        endTime, err := parseISOTime(slot.EndTime)
        if err != nil {
            utils.HandleAuthResponse(dto.ResponseDto{}, err, "Invalid end time format: "+slot.EndTime, res)
            return
        }

        // Extract just the time part (ignore date)
        startTimeOnly := time.Date(0, 1, 1, startTime.Hour(), startTime.Minute(), startTime.Second(), 0, time.UTC)
        endTimeOnly := time.Date(0, 1, 1, endTime.Hour(), endTime.Minute(), endTime.Second(), 0, time.UTC)

        serviceDTO.Unavailability = append(serviceDTO.Unavailability, regDto.LecturerUnavailability{
            Reason:    slot.Reason,
            Day:       slot.Day,
            StartTime: startTimeOnly,
            EndTime:   endTimeOnly,
        })
    }

    resp, errMsg, err := rh.RegService.CreateLecturerUnavailability(req.Context(), serviceDTO, utils.StringToUUID(lecturerId))
    utils.HandleAuthResponse(resp, err, errMsg, res)
}

// Helper function to parse ISO time strings
func parseISOTime(timeStr string) (time.Time, error) {
    // Try RFC3339 format (e.g., "1970-01-01T08:00:00Z")
    t, err := time.Parse(time.RFC3339, timeStr)
    if err == nil {
        return t, nil
    }

    // Try RFC3339 with nanoseconds
    t, err = time.Parse(time.RFC3339Nano, timeStr)
    if err == nil {
        return t, nil
    }

    // Try simple time format (HH:MM:SS)
    t, err = time.Parse("15:04:05", timeStr)
    if err == nil {
        return t, nil
    }

    // Try time format without seconds (HH:MM)
    t, err = time.Parse("15:04", timeStr)
    if err == nil {
        return t, nil
    }

    return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
}