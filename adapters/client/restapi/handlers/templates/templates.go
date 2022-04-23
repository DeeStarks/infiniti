package templates

import (
	"encoding/json"
	"net/http"
)

type withData struct {
	Status  	string      `json:"status"`
	Message 	string      `json:"message"`
	Data    	interface{} `json:"data"`
}

type withoutData struct {
	Status  	string      `json:"status"`
	Message 	string      `json:"message"`
}

// Code: Status code (e.g. 404, 200, 500)
// Message: Message to display to the user
// Data: Data to display to the user (applicable only to success responses)
func Template(w http.ResponseWriter, code int, message string, data interface{}) ([]byte, error) {
	switch code {
	case http.StatusOK:
		w.WriteHeader(http.StatusOK)
		res := withData{
			Status:   http.StatusText(http.StatusOK),
			Message:  message,
			Data:     data,
		}
		return json.Marshal(res)
	case http.StatusCreated:
		w.WriteHeader(http.StatusCreated)
		res := withData{
			Status:   http.StatusText(http.StatusCreated),
			Message:  message,
			Data:     data,
		}
		return json.Marshal(res)
	case http.StatusNotFound:
		w.WriteHeader(http.StatusNotFound)
		res := withoutData{
			Status:   http.StatusText(http.StatusNotFound),
			Message:  message,
		}
		return json.Marshal(res)
	case http.StatusBadRequest:
		w.WriteHeader(http.StatusBadRequest)
		res := withoutData{
			Status:   http.StatusText(http.StatusBadRequest),
			Message:  message,
		}
		return json.Marshal(res)
	case http.StatusMethodNotAllowed:
		w.WriteHeader(http.StatusMethodNotAllowed)
		res := withoutData{
			Status:   http.StatusText(http.StatusMethodNotAllowed),
			Message:  message,
		}
		return json.Marshal(res)
	case http.StatusUnauthorized:
		w.WriteHeader(http.StatusUnauthorized)
		res := withoutData{
			Status:   http.StatusText(http.StatusUnauthorized),
			Message:  message,
		}
		return json.Marshal(res)
	// TODO: Add more cases here


	default:
		w.WriteHeader(http.StatusInternalServerError)
		res := withoutData{
			Status:   http.StatusText(http.StatusInternalServerError),
			Message:  message,
		}
		return json.Marshal(res)
	}
}