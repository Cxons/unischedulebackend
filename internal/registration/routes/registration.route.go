package routes

import (
	authMiddleware "github.com/Cxons/unischedulebackend/internal/auth/middleware"
	regHandler "github.com/Cxons/unischedulebackend/internal/registration/handler"
	regMiddleware "github.com/Cxons/unischedulebackend/internal/registration/middleware"
	"github.com/go-chi/chi/v5"
)


func Routes(regHandler regHandler.RegHandler) chi.Router {
	r := chi.NewRouter()

	r.Use(authMiddleware.JwtMiddleware())
	adminMiddleware := regMiddleware.AdminMiddleware(regHandler.RegService)
	deanMiddleware := regMiddleware.DeanMiddleware(regHandler.RegService)
	hodMiddleware := regMiddleware.HodMiddleware(regHandler.RegService)

	r.Route("/admin",func(r chi.Router) {
		r.Use(adminMiddleware)
		r.Post("/update",regHandler.UpdateAdmin)
		r.Post("/university",regHandler.CreateUniversity)
		r.Get("/pendingdeans",regHandler.RetrievePendingDeans)
		r.Get("/approvedean",regHandler.ApproveDean)
	})

	
	r.Post("/dean/confirm",regHandler.RequestDeanConfirmation)
	r.Get("/dean/confirm",regHandler.CheckDeanConfirmation)

	r.Route("/dean",func(r chi.Router) {
		r.Use(deanMiddleware)
		r.Post("/",regHandler.CreateDean)
		r.Post("/faculty",regHandler.CreateFaculty)
		r.Get("/pendinghods",regHandler.RetrievePendingHods)
		r.Get("/approvehod",regHandler.ApproveHod)
	})

	

	r.Post("/confirm",regHandler.RequestHodConfirmation)
	r.Get("/confirm",regHandler.CheckHodConfirmation)
	r.Route("/hod",func(r chi.Router) {
		r.Use(hodMiddleware)
		r.Get("/",regHandler.CreateHod)
		r.Post("/department",regHandler.CreateDeparment)
		r.Get("/pendinglecturers",regHandler.RetrievePendingLecturers)
		r.Get("/approvelecturers",regHandler.ApproveLecturer)
	})


	r.Post("/lecturer/confirm",regHandler.RequestLecturerConfirmation)
	r.Get("/lecturer/confirm",regHandler.CheckLecturerConfirmation)


	return r;
}