package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deestarks/infiniti/adapters/client/restapi/constants"
	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/utils"
)

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.Context().Value(constants.CTXKey("group_name")) == "user" {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Unauthorized to perform transaction", nil)
			w.Write([]byte(res))
			return
		}

		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			tmp, _ := templates.Template(w, http.StatusBadRequest, "Invalid JSON data", nil)
			w.Write([]byte(tmp))
			return
		}

		service := h.appPort.NewTransactionsService()
		transaction, err := service.Deposit(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Deposit successful", transaction)
		w.Write([]byte(res))
	default:
		tmp, _ := templates.Template(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		w.Write([]byte(tmp))
	}
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		tmp, _ := templates.Template(w, http.StatusBadRequest, "Invalid JSON data", nil)
		w.Write([]byte(tmp))
		return
	}

	if utils.GetRouteGroup(r) == "user" {
		userIdCtx := r.Context().Value(constants.CTXKey("user_id"))
		userID, err := strconv.Atoi(fmt.Sprintf("%v", userIdCtx))
		if err != nil {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Invalid user id", nil)
			w.Write([]byte(res))
			return
		}
		data["sender_id"] = userID // Override user_id with the user_id from the context
	}

	switch r.Method {
	case http.MethodPost:
		service := h.appPort.NewTransactionsService()
		transaction, err := service.Transfer(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Transfer successful", transaction)
		w.Write([]byte(res))
	default:
		tmp, _ := templates.Template(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		w.Write([]byte(tmp))
	}
}