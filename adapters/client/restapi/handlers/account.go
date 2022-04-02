package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/utils"
	"github.com/gorilla/mux"
)

func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accts, err := h.appPort.NewAccountService().ListAccounts()
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(err.StatusCode(), err.Error(), nil) // If error occurs here, it's a JSON marshal error
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(res))
		return
	}

	res, _ := templates.Template(http.StatusOK, "Accounts successfully retrieved", accts) // If error occurs here, it's a JSON marshal error
	w.Write([]byte(res))
}

func (h *Handler) SingleAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the URL variables
	idVar, ok := vars["id"] // Get the id variable
	if !ok {
		res, _ := templates.Template(http.StatusBadRequest, "Missing \"id\" variable", nil)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(res))
		return
	}

	// Cast the id variable to an int
	id, err := strconv.Atoi(idVar)
	if err != nil {
		res, _ := templates.Template(http.StatusBadRequest, "Invalid \"id\" variable. Must be an integer", nil)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(res))
		return
	}

	switch r.Method {
	case "GET":
		acct, err := h.appPort.NewAccountService().GetAccount("id", id, []string{})
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(err.StatusCode(), err.Error(), nil)
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(res))
			return
		}
	
		res, _ := templates.Template(http.StatusOK, "Account successfully retrieved", acct)
		w.Write([]byte(res))
	case "PUT":
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			res, _ := templates.Template(http.StatusBadRequest, "Invalid JSON data", nil)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(res))
			return
		}
		fmt.Println(data)
		// TODO: Validate data
		
	default:
		res, _ := templates.Template(http.StatusMethodNotAllowed, "Accepts only GET and PUT requests", nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(res))
	}
}