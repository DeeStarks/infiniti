package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/gorilla/mux"
)

func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accts, err := h.appPort.NewAccountService().ListAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonRpr, err := json.Marshal(
		templates.Success(accts, "Accounts successfully retrieved"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(jsonRpr))
}

func (h *Handler) SingleAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the URL variables
	idVar, ok := vars["id"] // Get the id variable
	if !ok {
		res := templates.ErrorBadRequest("Account ID is required")
		jsonRpr, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonRpr))
		return
	}

	// Cast the id variable to an int
	id, err := strconv.Atoi(idVar)
	if err != nil {
		res := templates.ErrorBadRequest("Account ID must be an integer")
		jsonRpr, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonRpr))
		return
	}

	switch r.Method {
	case "GET":
		acct, err := h.appPort.NewAccountService().GetAccount("id", id, []string{})
		if err != nil {
			res := templates.ErrorNotFound("Account not found")
			jsonRpr, err := json.Marshal(res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(jsonRpr))
			return
		}
	
		jsonRpr, err := json.Marshal(
			templates.Success(acct, "Account successfully retrieved"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(jsonRpr))
	case "PUT":
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			res := templates.ErrorBadRequest("Invalid Data")
			jsonRpr, err := json.Marshal(res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(jsonRpr))
			return
		}
		// TODO: Validate data
		
	case "DELETE":
		// TODO: Implement DELETE

		jsonRpr, err := json.Marshal(
			templates.Success(nil, "Account successfully deleted"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(jsonRpr))
	default:
		res := templates.ErrorMethodNotAllowed(
			fmt.Sprintf("Method %s not allowed", r.Method))
		jsonRpr, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(jsonRpr))
	}
}