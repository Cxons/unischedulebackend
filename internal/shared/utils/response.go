package utils

import (
	"encoding/json"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/shared/dto"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	"github.com/Cxons/unischedulebackend/pkg/validator"
)


func HandleAuthResponse(resp dto.ResponseDto, err error, errMsg string, res http.ResponseWriter) {
    res.Header().Set("Content-Type", "application/json") // always JSON

    if err != nil {
        // Map status message to HTTP code
        code := status.RetrieveCodeFromStatusMessage(errMsg)
        if code == 0 {
            code = status.InternalServerError.Code
            errMsg = status.InternalServerError.Message
        } else if errMsg == "" {
            errMsg = err.Error()
        }

        // Return JSON error
        res.WriteHeader(code)
        json.NewEncoder(res).Encode(map[string]interface{}{
            "message": errMsg,
            "error":   err.Error(),
        })
        return
    }

    // Success response
    res.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(res).Encode(resp); err != nil {
        // fallback in case encoding fails
        res.WriteHeader(status.InternalServerError.Code)
        json.NewEncoder(res).Encode(map[string]interface{}{
            "message": status.InternalServerError.Message,
            "error":   err.Error(),
        })
        return
    }
}




func HandleBodyParsing(req *http.Request, res http.ResponseWriter, body interface{}){
	if err:= json.NewDecoder(req.Body).Decode(body); err!=nil{
		http.Error(res,"Invalid Request Body",status.BadRequest.Code)
		return
	}
	if err := validator.ValidateStruct(body); err!= nil{
		http.Error(res,"Validation Error: " + err.Error(),status.BadRequest.Code)
		return
	}
}