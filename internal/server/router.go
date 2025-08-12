package server

import (
	authRoutes "github.com/Cxons/unischedulebackend/internal/auth/routes"
	regRoutes "github.com/Cxons/unischedulebackend/internal/registration/routes"
	"github.com/go-chi/chi/v5"
)


func (s *Server) mountRoutes(){
	s.Router.Route("/api/v1",func(r chi.Router) {
		r.Mount("/auth",authRoutes.Routes(*s.Auth))
		r.Mount("/registration",regRoutes.Routes(*s.Reg))
		// r.Mount("/dean",dean.Routes())
		// r.Mount("/hod",hod.Routes())
		// r.Mount("/lecturer", lecturer.Routes())

	})
}