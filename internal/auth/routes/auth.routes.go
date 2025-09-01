package routes

import (
	authHandler "github.com/Cxons/unischedulebackend/internal/auth/handler"
	"github.com/go-chi/chi/v5"
)




func Routes(authHandler authHandler.AuthHandler) chi.Router {
	r:= chi.NewRouter()
	r.Post("/register",authHandler.Register)
	r.Post("/login",authHandler.Login)
	r.Get("/accesstoken",authHandler.RefreshAccessToken)
	r.Post("/otp",authHandler.SendOtp)
	r.Post("/otp/verify",authHandler.VerifyOtp)
	return r
}




