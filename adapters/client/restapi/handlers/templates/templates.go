package templates

import "net/http"

func Success(data interface{}, message string) map[string]interface{} {
	return map[string]interface{}{
		"status": http.StatusText(http.StatusOK),
		"data":   data,
		"message": message,
	}
}

func ErrorNotFound(message string) map[string]interface{} {
	return map[string]interface{}{
		"status": http.StatusText(http.StatusNotFound),
		"message": message,
	}
}

func ErrorBadRequest(message string) map[string]interface{} {
	return map[string]interface{}{
		"status": http.StatusText(http.StatusBadRequest),
		"message": message,
	}
}

func ErrorMethodNotAllowed(message string) map[string]interface{} {
	return map[string]interface{}{
		"status": http.StatusText(http.StatusMethodNotAllowed),
		"message": message,
	}
}