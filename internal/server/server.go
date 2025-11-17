package server

import (
	"log/slog"
	"net/http"
	"time"

	authHandler "github.com/Cxons/unischedulebackend/internal/auth/handler"
	courseHandler "github.com/Cxons/unischedulebackend/internal/courses/handler"
	regHandler "github.com/Cxons/unischedulebackend/internal/registration/handler"
	"github.com/Cxons/unischedulebackend/internal/shared/db"
	timetableHandler "github.com/Cxons/unischedulebackend/internal/timetable/handler"
	uniHandler "github.com/Cxons/unischedulebackend/internal/university/handler"
	"github.com/Cxons/unischedulebackend/pkg/caching"

	supHandler "github.com/Cxons/unischedulebackend/pkg/supabase/handler"
	"github.com/go-chi/chi/v5"
)


type Server struct{
	Config Config
	Router chi.Router
	Server *http.Server
	Logger *slog.Logger
	Cache  *caching.RedisClient
	Database *db.Db
	Auth *authHandler.AuthHandler
	Reg *regHandler.RegHandler
	Supabase *supHandler.SupabaseHandler
	Uni *uniHandler.UniversityHandler
	Course *courseHandler.CourseHandler
	Timetable *timetableHandler.TimetableHandler
}


func NewServer(cfg Config, logger *slog.Logger, cache *caching.RedisClient, db *db.Db, auth *authHandler.AuthHandler, reg *regHandler.RegHandler, supabase *supHandler.SupabaseHandler,uniHandler *uniHandler.UniversityHandler, courseHandler *courseHandler.CourseHandler, timetableHandler *timetableHandler.TimetableHandler) *Server{
	r:= chi.NewRouter();
	srv := &http.Server{
    Addr:         ":5000",
    Handler:      r,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
}

	return &Server{
		Config: cfg,
		Router: r,
		Server: srv,
		Logger: logger,
		Cache: cache,
		Database: db,
		Auth: auth,
		Reg: reg,
		Supabase: supabase,
		Uni: uniHandler,
		Course: courseHandler,
		Timetable: timetableHandler,
	}
}


func (s *Server) Start() error{
	s.mountRoutes()
	return s.Server.ListenAndServe()
}