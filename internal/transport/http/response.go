package transport

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Message: msg,
	}
}
