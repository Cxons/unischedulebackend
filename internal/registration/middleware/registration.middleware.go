package middleware

import (
	"log/slog"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/registration/service"
	"github.com/Cxons/unischedulebackend/internal/shared/constants"
	"github.com/Cxons/unischedulebackend/pkg/auth/jwt"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
)




var ADMIN = constants.ADMIN
var LECTURER = constants.LECTURER



func AdminMiddleware(regService service.RegService)func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request){
			ctx := req.Context()
		
			claims := ctx.Value(constants.UserInfoKey)
			  if claims == nil {
                http.Error(res, "Unauthorized", status.Unauthorized.Code)
                return
            }
			userInfo := claims.(*jwt.CustomClaims)
			adminId := userInfo.User_id

		if userInfo.Role != ADMIN {
				http.Error(res,"Unauthorized",status.Unauthorized.Code)
				return
			}

			confirmed,err := regService.CheckCurrentAdmin(ctx,adminId)
			if err != nil{
					http.Error(res,"Internal Server Error" + err.Error(),status.InternalServerError.Code)
					return
			}
			if !confirmed{
					http.Error(res,"Unauthorized",status.Unauthorized.Code)
					return
			}
			next.ServeHTTP(res,req.WithContext(ctx))
		})
	}
}

func DeanMiddleware(regService service.RegService)func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request){
			ctx := req.Context()
			claims := ctx.Value(constants.UserInfoKey)
			  if claims == nil {
                http.Error(res, "Unauthorized", status.Unauthorized.Code)
                return
            }
			userInfo := claims.(*jwt.CustomClaims)
			// adminId := userInfo.User_id

			cookieValue,err := req.Cookie("current_dean_id")
			if err != nil{
				slog.Error("Error retrieving current dean id","err:",err)
				http.Error(res,"Error retrieving current dean id",status.InternalServerError.Code)
				return
	}
			deanId := cookieValue.Value

		if userInfo.Role != LECTURER {
				http.Error(res,"Unauthorized",status.Unauthorized.Code)
				return
			}

			confirmed,err := regService.CheckCurrentDean(ctx,deanId)
			if err != nil{
					http.Error(res,"Internal Server Error" + err.Error(),status.InternalServerError.Code)
					return
			}
			if !confirmed{
					http.Error(res,"Unauthorized",status.Unauthorized.Code)
					return
			}
			next.ServeHTTP(res,req.WithContext(ctx))
		})
	}
}




func HodMiddleware(regService service.RegService)func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request){
			ctx := req.Context()
			claims := ctx.Value(constants.UserInfoKey)
			  if claims == nil {
                http.Error(res, "Unauthorized", status.Unauthorized.Code)
                return
            }
			userInfo := claims.(*jwt.CustomClaims)
			// adminId := userInfo.User_id


			cookieValue,err := req.Cookie("current_hod_id")
			if err != nil{
				slog.Error("Error retrieving current hod id","err:",err)
				http.Error(res,"Error retrieving current hod id",status.InternalServerError.Code)
				return
	}
			hodId := cookieValue.Value


		if userInfo.Role != LECTURER {
				http.Error(res,"Unauthorized",status.Unauthorized.Code)
				return
			}

			confirmed,err := regService.CheckCurrentHod(ctx,hodId)
			if err != nil{
					http.Error(res,"Internal Server Error" + err.Error(),status.InternalServerError.Code)
					return
			}
			if !confirmed{
					http.Error(res,"Unauthorized",status.Unauthorized.Code)
					return
			}
			next.ServeHTTP(res,req.WithContext(ctx))
		})
	}
}