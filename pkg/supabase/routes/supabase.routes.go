package routes

import (
	supHandler "github.com/Cxons/unischedulebackend/pkg/supabase/handler"
	"github.com/go-chi/chi/v5"
)




func Routes(supHandler supHandler.SupabaseHandler) chi.Router{
	r := chi.NewRouter()
	r.Get("/signedUrl",supHandler.GetSignedUrl)
	return r
}





