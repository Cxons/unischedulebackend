package routes

import (
	authMiddleware "github.com/Cxons/unischedulebackend/internal/auth/middleware"
	courseHandler "github.com/Cxons/unischedulebackend/internal/courses/handler"
	regHandler "github.com/Cxons/unischedulebackend/internal/registration/handler"
	regMiddleware "github.com/Cxons/unischedulebackend/internal/registration/middleware"

	"github.com/go-chi/chi/v5"
)


func Routes(courseHandler courseHandler.CourseHandler,regHandler regHandler.RegHandler) chi.Router {
	r := chi.NewRouter()

	r.Use(authMiddleware.JwtMiddleware())
	// adminMiddleware := regMiddleware.AdminMiddleware(regHandler.RegService)
	// deanMiddleware := regMiddleware.DeanMiddleware(regHandler.RegService)
	hodMiddleware := regMiddleware.HodMiddleware(regHandler.RegService)

	r.Route("/create",func(r chi.Router) {
		r.Use(hodMiddleware)
		r.Post("/",courseHandler.CreateCourse)
		r.Post("/possiblevenues",courseHandler.SetCoursePossibleVenues)
	})
	r.Delete("/possiblevenue",courseHandler.DeleteCoursePossibleVenue)
	r.Get("/possiblevenues",courseHandler.FetchCoursePossibleVenues)
	r.Post("/department/all",courseHandler.RetrieveCoursesForADepartment)
	r.Get("/cohort",courseHandler.RetrieveCoursesForACohort)
	r.Post("/cohort",courseHandler.SetCoursesForACohort)
	r.Get("/all",courseHandler.FetchAllCourses)


	return r;
}