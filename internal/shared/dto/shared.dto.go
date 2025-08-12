package dto



type ResponseDto struct {
	Message string
	Data interface{}
	StatusCode int
	StatusCodeMessage string
}