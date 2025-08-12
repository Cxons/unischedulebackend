package utils

import (
	"encoding/json"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/shared/dto"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
)




func HandleAuthResponse(resp dto.ResponseDto,err error,errMsg string, res http.ResponseWriter){
	 if err != nil{
		code := status.RetrieveCodeFromStatusMessage(errMsg)
		if code == 0 {
			http.Error(res,status.InternalServerError.Message,status.InternalServerError.Code)
			return
		}
		http.Error(res,err.Error(),code)
		return
	 }
	// res.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(res).Encode(resp); err!= nil{
    http.Error(res, status.InternalServerError.Message, status.InternalServerError.Code)
    return
}
}