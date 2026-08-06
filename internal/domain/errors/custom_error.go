package errors

// CustomException mirrors DIGIT's standard error envelope:
// {"Errors":[{"code":"...","message":"..."}]}
type ErrorField struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []ErrorField `json:"Errors"`
}

type CustomException struct {
	Code    string
	Message string
}

func (e *CustomException) Error() string {
	return e.Message
}

func New(code, message string) *CustomException {
	return &CustomException{Code: code, Message: message}
}

func (e *CustomException) Response() ErrorResponse {
	return ErrorResponse{Errors: []ErrorField{{Code: e.Code, Message: e.Message}}}
}
