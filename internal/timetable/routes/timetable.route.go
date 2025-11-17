package routes

import (
	timetableHandler "github.com/Cxons/unischedulebackend/internal/timetable/handler"
	"github.com/go-chi/chi/v5"
)





func Routes (timetableHandler timetableHandler.TimetableHandler)chi.Router{
	r := chi.NewRouter()

	r.Post("/",timetableHandler.CreateATimeTable)
	r.Get("/cohort",timetableHandler.FetchTimetableForCohort)

	return r
}
