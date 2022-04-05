package templates

import (
	"encoding/json"
	"net/http"
)

// Code: Status code (e.g. 404, 200, 500)
// Message: Message to display to the user
// Data: Data to display to the user (applicable only to success responses)
func Template(w http.ResponseWriter, code int, message string, data interface{}) ([]byte, error) {
	switch code {
	case http.StatusOK:
		w.WriteHeader(http.StatusOK)
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusOK),
			"data":   data,
			"message": message,
		}
		return json.Marshal(res)
	case http.StatusNotFound:
		w.WriteHeader(http.StatusNotFound)
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusNotFound),
			"message": message,
		}
		return json.Marshal(res)
	case http.StatusBadRequest:
		w.WriteHeader(http.StatusBadRequest)
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusBadRequest),
			"message": message,
		}
		return json.Marshal(res)
	case http.StatusMethodNotAllowed:
		w.WriteHeader(http.StatusMethodNotAllowed)
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusMethodNotAllowed),
			"message": message,
		}
		return json.Marshal(res)
	case http.StatusUnauthorized:
		w.WriteHeader(http.StatusUnauthorized)
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusUnauthorized),
			"message": message,
		}
		return json.Marshal(res)
	// TODO: Add more cases here


	default:
		w.WriteHeader(http.StatusInternalServerError)
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusInternalServerError),
			"message": message,
		}
		return json.Marshal(res)
	}
}