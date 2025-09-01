package utils

import (
	"encoding/json"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/shared/dto"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	"github.com/Cxons/unischedulebackend/pkg/validator"
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



func HandleBodyParsing(req *http.Request, res http.ResponseWriter, body interface{}){
	if err:= json.NewDecoder(req.Body).Decode(&body); err!=nil{
		http.Error(res,"Invalid Request Body",status.BadRequest.Code)
		return
	}
	if err := validator.ValidateStruct(body); err!= nil{
		http.Error(res,"Validation Error: " + err.Error(),status.BadRequest.Code)
		return
	}
}