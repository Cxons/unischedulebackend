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
		r.Get("/pendingdean",regHandler.FetchDeanWaitDetails)
		r.Get("/pendingdeans",regHandler.RetrievePendingDeans)
		r.Post("/approvedean",regHandler.ApproveDean)
	})

	
	r.Post("/dean/confirm",regHandler.RequestDeanConfirmation)
	r.Get("/dean/confirm",regHandler.CheckDeanConfirmation)
	r.Post("/dean/faculty",regHandler.CreateFaculty)
	r.Route("/dean",func(r chi.Router) {
		r.Use(deanMiddleware)
		r.Post("/",regHandler.CreateDean)
		r.Get("/pendinghod",regHandler.FetchHodWaitDetails)
		r.Get("/pendinghods",regHandler.RetrievePendingHods)
		r.Post("/approvehod",regHandler.ApproveHod)
	})

	

	r.Post("/hod/confirm",regHandler.RequestHodConfirmation)
	r.Get("/hod/confirm",regHandler.CheckHodConfirmation)
	r.Post("/hod/department",regHandler.CreateDeparment)
	r.Route("/hod",func(r chi.Router) {
		r.Use(hodMiddleware)
		r.Get("/",regHandler.CreateHod)
		r.Get("/pendinglecturers",regHandler.RetrievePendingLecturers)
		r.Post("/approvelecturer",regHandler.ApproveLecturer)
		r.Get("/pendinglecturer",regHandler.FetchLecturerWaitDetails)
	})


	r.Post("/lecturer/confirm",regHandler.RequestLecturerConfirmation)
	r.Get("/lecturer/confirm",regHandler.CheckLecturerConfirmation)

	r.Post("/lecturer/unavailability",regHandler.CreateLecturerUnavailability)


	return r;
}