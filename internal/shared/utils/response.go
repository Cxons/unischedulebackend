package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Cxons/unischedulebackend/internal/shared/dto"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	"github.com/Cxons/unischedulebackend/pkg/validator"
)

type Response struct {
  	Message string
	Data interface{}
	StatusCode int
	StatusCodeMessage string
}
func HandleAuthResponse(resp dto.ResponseDto, err error, errMsg string, res http.ResponseWriter) {
    // Set content type first
    res.Header().Set("Content-Type", "application/json")

    if err != nil {
        code := status.RetrieveCodeFromStatusMessage(errMsg)
        if code == 0 {
            code = status.InternalServerError.Code
            errMsg = status.InternalServerError.Message
        }
        
        slog.Info("error response", "code", code, "error", err.Error())
        
        // Safe write - only if not already written
        writeSafeJSON(res, code, map[string]interface{}{
            "message": errMsg,
            "error":   err.Error(),
        })
        return
    }

    if resp.Data == nil {
        resp.Data = map[string]interface{}{}
    }

    writeSafeJSON(res, http.StatusCreated, resp)
}

func writeSafeJSON(res http.ResponseWriter, code int, data interface{}) {
    // Try to write, but recover from panic if headers already written
    defer func() {
        if r := recover(); r != nil {
            slog.Warn("headers already written", "recover", r)
        }
    }()
    
    res.WriteHeader(code)
    json.NewEncoder(res).Encode(data)
}

func HandleBodyParsing(req *http.Request, res http.ResponseWriter, body interface{}) error {
    if err := json.NewDecoder(req.Body).Decode(body); err != nil {
        sendError(res, status.BadRequest.Code, "Invalid Request Body", err)
        return err
    }
    
    if err := validator.ValidateStruct(body); err != nil {
        sendError(res, status.BadRequest.Code, "Validation Error", err)
        return err
    }
    
    return nil
}

func sendError(res http.ResponseWriter, code int, message string, err error) {
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(code)
    json.NewEncoder(res).Encode(Response{
        Message: message,
        StatusCodeMessage:   err.Error(),
    })
}