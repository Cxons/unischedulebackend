package server

import (
	"log/slog"
	"net/http"
	"time"

	authHandler "github.com/Cxons/unischedulebackend/internal/auth/handler"
	regHandler "github.com/Cxons/unischedulebackend/internal/registration/handler"
	"github.com/Cxons/unischedulebackend/internal/shared/db"
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
}


func NewServer(cfg Config, logger *slog.Logger, cache *caching.RedisClient, db *db.Db, auth *authHandler.AuthHandler, reg *regHandler.RegHandler, supabase *supHandler.SupabaseHandler) *Server{
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
	}
}


func (s *Server) Start() error{
	s.mountRoutes()
	return s.Server.ListenAndServe()
}