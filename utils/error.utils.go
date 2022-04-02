package utils

type RequestError struct {
	Code int
	Err error
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

func (r *RequestError) StatusCode() int {
	return r.Code
}
