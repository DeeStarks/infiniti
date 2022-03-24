package handlers

import (
	"net/http"
	"encoding/json"
)

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) {
	jsonWelcome := map[string]string{
		"status": http.StatusText(http.StatusOK),
		"message": "Welcome to Infiniti Bank API",
	}
	jsonRpr, err := json.Marshal(jsonWelcome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(jsonRpr))
}