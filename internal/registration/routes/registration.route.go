package routes

import (
	authmiddleware "github.com/Cxons/unischedulebackend/internal/auth/middleware"
	regHandler "github.com/Cxons/unischedulebackend/internal/registration/handler"
	"github.com/go-chi/chi/v5"
)


func Routes(regHandler regHandler.RegHandler) chi.Router {
	r := chi.NewRouter()

	r.Use(authmiddleware.JwtMiddleware())

	r.Post("/admin/update",regHandler.UpdateAdmin)
	r.Post("/admin/university",regHandler.CreateUniversity)
	r.Get("/admin/pendingdeans",regHandler.RetrievePendingDeans)

	return r;
}