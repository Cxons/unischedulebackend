package server

import (
	authmiddleware "github.com/Cxons/unischedulebackend/internal/auth/middleware"
	authRoutes "github.com/Cxons/unischedulebackend/internal/auth/routes"
	courseRoutes "github.com/Cxons/unischedulebackend/internal/courses/routes"
	regRoutes "github.com/Cxons/unischedulebackend/internal/registration/routes"
	timetableRoutes "github.com/Cxons/unischedulebackend/internal/timetable/routes"
	uniRoutes "github.com/Cxons/unischedulebackend/internal/university/routes"
	supRoutes "github.com/Cxons/unischedulebackend/pkg/supabase/routes"
	"github.com/go-chi/chi/v5"
)


func (s *Server) mountRoutes(){
	s.Router.Use(authmiddleware.CORSMiddleware)

	s.Router.Route("/api/v1",func(r chi.Router) {
		r.Mount("/auth",authRoutes.Routes(*s.Auth))
		r.Mount("/registration",regRoutes.Routes(*s.Reg))
		r.Mount("/supabase",supRoutes.Routes(*s.Supabase))
		r.Mount("/university",uniRoutes.Routes(*s.Uni))
		r.Mount("/course",courseRoutes.Routes(*s.Course,*s.Reg))
		r.Mount("/timetable",timetableRoutes.Routes(*s.Timetable))
		// r.Mount("/dean",dean.Routes())
		// r.Mount("/hod",hod.Routes())
		// r.Mount("/lecturer", lecturer.Routes())

	})
}