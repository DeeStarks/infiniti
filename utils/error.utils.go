package utils

type (
	appError struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}
	AppError interface {
		Error() 		string
		StatusCode()	int
	}
)

func NewAppError(message string, status int) AppError {
	return appError{
		Message: message,
		Status:  status,
	}
}

func (e appError) Error() string {
	return e.Message
}

func (e appError) StatusCode() int {
	return e.Status
}