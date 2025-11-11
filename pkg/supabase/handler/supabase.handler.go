package handler

import (
	"log/slog"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/Cxons/unischedulebackend/pkg/supabase/service"
)


type SupabaseHandlerInterface interface{
	GetSignedUrl(res http.ResponseWriter,req *http.Request)
}


type SupabaseHandler struct {
	service service.SupabaseService
}



func NewSupabasePackage(logger *slog.Logger,supabaseUrl string,supabaseSecretKey string)*SupabaseHandler{

	supabaseService := service.NewSupabaseService(supabaseUrl,supabaseSecretKey,logger)
	// initialize supabase client
	supabaseService.InitClient()
	handler := NewSupabaseHandler(supabaseService)
	return handler

}


func NewSupabaseHandler(service service.SupabaseService)*SupabaseHandler{
	return &SupabaseHandler{
		service: service,
	}
}


func (h *SupabaseHandler) GetSignedUrl(res http.ResponseWriter,req *http.Request){
	queryParams := req.URL.Query()
	fileName := queryParams.Get("file_name")
	resp,errMsg,err := h.service.CreateSignedUrl(fileName)
	utils.HandleAuthResponse(resp,err,errMsg,res)
}

