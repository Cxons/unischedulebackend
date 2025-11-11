package authHandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	// "database/sql"

	"log/slog"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/auth/dto"
	"github.com/Cxons/unischedulebackend/internal/auth/repository"
	"github.com/Cxons/unischedulebackend/internal/auth/service"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
)

type AuthHandlerInterface interface {
 Register(res http.ResponseWriter, req *http.Request)
 Login(res http.ResponseWriter,req *http.Request)
 RefreshAccessToken(res http.ResponseWriter, req *http.Request)
 SendOtp(res http.ResponseWriter, req *http.Request)
 VerifyOtp(res http.ResponseWriter, req *http.Request)
}

var ctx  = context.Background()


type AuthHandler struct {
	service service.AuthService
}

func NewAuthPackage(logger *slog.Logger, db *sql.DB)*AuthHandler{
	query := sqlc.New(db)

	//initializes queries
	studentQueries := queries.NewStudentQueries(query)
	lecturerQueries := queries.NewLecturerQueries(query)
	adminQueries := queries.NewAdminQueries(query)
	tokenQueries := queries.NewTokenQueries(query)
	otpQueries := queries.NewOtpQueries(query)

	// initializes repository
	repo := repository.NewAuthRepository(studentQueries,lecturerQueries,adminQueries,tokenQueries,otpQueries,logger)

	
	// initializes service
	service := service.NewAuthService(repo,logger)

	// initializes handler
	handler := NewAuthHandler(service)

	return handler

}

func NewAuthHandler(service service.AuthService)*AuthHandler{
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(res http.ResponseWriter, req *http.Request){
	var body service.RegisterDto;
	utils.HandleBodyParsing(req,res,&body)
	// ctx := context.Background()
	 resp,errMsg,err := h.service.Register(ctx,body)
	 utils.HandleAuthResponse(resp,err,errMsg,res)
	
}

func (h *AuthHandler) Login(res http.ResponseWriter, req *http.Request){
     var body service.LoginDto;
	 
	 utils.HandleBodyParsing(req,res,&body)
	// passes data into service layer for handling
	resp,errMsg,err:= h.service.Login(ctx,body)
	data,ok := resp.Data.(dto.LoginResponseData)
	if !ok {
    slog.Error("Login response data is invalid", "resp.Data", resp.Data)

    res.Header().Set("Content-Type", "application/json")

    // If the error corresponds to unauthorized (e.g., status.Unauthorized)
    if errMsg == status.Unauthorized.Message {
        res.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(res).Encode(map[string]interface{}{
            "message": "Invalid email or password",
            "error":   "Unauthorized",
        })
        return
    }

    // Fallback for other cases
    res.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(res).Encode(map[string]interface{}{
        "message": "Invalid login response data",
        "error":   "Login response data is nil or of unexpected type",
    })
    return
}

	
	if err == nil{
		// create cookie to be added to user request
			cookie := &http.Cookie{
			Name: "refresh_token",
			Value: data.RefreshToken,
			Path: "/",
			HttpOnly: true,
			Secure: false,
			SameSite: http.SameSiteLaxMode,
			Expires: time.Now().Add(time.Hour * 24 * 30),
	}
		slog.Info("cookie","key",cookie)
		http.SetCookie(res,cookie)
	}
	modifiedResp := service.AuthResponse{
		Message: resp.Message,
		Data: dto.LoginResponseData{
		AccessToken: resp.Data.(dto.LoginResponseData).AccessToken,
	},
	StatusCode: resp.StatusCode,
	StatusCodeMessage: resp.StatusCodeMessage,
	}
	

	// return appropriate response message to user
	utils.HandleAuthResponse(modifiedResp,err,errMsg,res)
	
}

func (h *AuthHandler) SendOtp(res http.ResponseWriter,req *http.Request){
	var body dto.SendOtpDto
	 
	utils.HandleBodyParsing(req,res,&body)
	 resp,errMsg,err := h.service.SendOtp(ctx,body.Email,body.UserType)
	 utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (h *AuthHandler) VerifyOtp(res http.ResponseWriter, req *http.Request){
	var body dto.VerifyOtpDto
	utils.HandleBodyParsing(req,res,&body)
	 resp,errMsg,err := h.service.VerifyOtp(ctx,body.Email,body.Otp)
	 utils.HandleAuthResponse(resp,err,errMsg,res)
}

func (h *AuthHandler) RefreshAccessToken(res http.ResponseWriter, req *http.Request){
	cookie,err := req.Cookie("refresh_token")
	if err != nil{
		slog.Error("Error retreiving refresh token","err:",err)
		http.Error(res,"Error retrieving refresh token",status.InternalServerError.Code)
		return
	}
	cookieValue := &cookie.Value
	resp,errMsg,err := h.service.RefreshAccessToken(ctx,cookieValue)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

// func (h *AuthHandler) handleAuthResponse(resp service.AuthResponse,err error,errMsg string, res http.ResponseWriter){
// 	 if err != nil{
// 		code := status.RetrieveCodeFromStatusMessage(errMsg)
// 		if code == 0 {
// 			http.Error(res,status.InternalServerError.Message,status.InternalServerError.Code)
// 			return
// 		}
// 		http.Error(res,err.Error(),code)
// 		return
// 	 }
// 	// res.WriteHeader(http.StatusCreated)
// 	if err = json.NewEncoder(res).Encode(resp); err!= nil{
//     http.Error(res, status.InternalServerError.Message, status.InternalServerError.Code)
//     return
// }
// }
