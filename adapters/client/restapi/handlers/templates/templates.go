package templates

import (
	"encoding/json"
	"net/http"
)

// Code: Status code (e.g. 404, 200, 500)
// Message: Message to display to the user
// Data: Data to display to the user (applicable only to success responses)
func Template(code int, message string, data interface{}) ([]byte, error) {
	switch code {
	case http.StatusOK:
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusOK),
			"data":   data,
			"message": message,
		}
		return json.Marshal(res)
	case http.StatusNotFound:
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusNotFound),
			"message": message,
		}
		return json.Marshal(res)
	case http.StatusBadRequest:
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusBadRequest),
			"message": message,
		}
		return json.Marshal(res)
	case http.StatusMethodNotAllowed:
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusMethodNotAllowed),
			"message": message,
		}
		return json.Marshal(res)
	case http.StatusUnauthorized:
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusUnauthorized),
			"message": message,
		}
		return json.Marshal(res)
	// TODO: Add more cases here


	default:
		res := map[string]interface{}{
			"status": http.StatusText(http.StatusInternalServerError),
			"message": message,
		}
		return json.Marshal(res)
	}
}