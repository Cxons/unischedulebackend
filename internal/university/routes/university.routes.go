package routes

import (
	uniHandler "github.com/Cxons/unischedulebackend/internal/university/handler"
	"github.com/go-chi/chi/v5"
)





func Routes (uniHandler uniHandler.UniversityHandler)chi.Router{
	r := chi.NewRouter()

	r.Get("/all",uniHandler.RetrieveAllUniversities)
	r.Get("/faculties",uniHandler.RetrieveAllFaculties)
	r.Get("/departments",uniHandler.RetrieveAllDepartments)
	r.Get("/all/departments",uniHandler.FetchAllDepartmentsForAUni)
	r.Get("/department/lecturers",uniHandler.RetrieveDepartmentLecturers)
	r.Post("/venue",uniHandler.CreateVenue)
	r.Get("/department/cohorts",uniHandler.FetchCohortsForADepartment)
	r.Get("/venues",uniHandler.RetrieveAllVenues)
	return r
}
