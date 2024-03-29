package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deestarks/infiniti/adapters/client/restapi/constants"
	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/utils"
	"github.com/gorilla/mux"
)

func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accts, err := h.appPort.NewAccountService().ListAccounts()
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
		w.Write([]byte(res))
		return
	}

	res, _ := templates.Template(w, http.StatusOK, "Accounts successfully retrieved", accts)
	w.Write([]byte(res))
}

func (h *Handler) SingleAccount(w http.ResponseWriter, r *http.Request) {
	var idVar string

	// Get route group (user, staff, admin)
	routeGroup := utils.GetRouteGroup(r)

	switch routeGroup {
	case "user":
		// Get the user ID from the request context
		id := r.Context().Value(constants.CTXKey("user_id"))
		if id == nil {
			res, _ := templates.Template(w, http.StatusUnauthorized, "User not logged in", nil)
			w.Write([]byte(res))
			return
		}
		idVar = fmt.Sprintf("%v", id)
	case "admin", "staff":
		urlVars := mux.Vars(r) // Get the URL variables
		id, ok := urlVars["id"] // Get the id variable
		if !ok {
			res, _ := templates.Template(w, http.StatusBadRequest, "Missing \"id\" variable", nil)
			w.Write([]byte(res))
			return
		}
		idVar = id
	default:
		res, _ := templates.Template(w, http.StatusBadRequest, "Invalid url", nil)
		w.Write([]byte(res))
		return
	}

	// Cast the id variable to an int
	id, err := strconv.Atoi(idVar)
	if err != nil {
		res, _ := templates.Template(w, http.StatusBadRequest, "Invalid \"id\" variable. Must be an integer", nil)
		w.Write([]byte(res))
		return
	}

	acctService := h.appPort.NewAccountService()
	switch r.Method {
	case http.MethodGet:
		acct, err := acctService.GetAccount("id", id, []string{})
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
	
		res, _ := templates.Template(w, http.StatusOK, "Account successfully retrieved", acct)
		w.Write([]byte(res))
	case http.MethodPut:
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			res, _ := templates.Template(w, http.StatusBadRequest, "Invalid JSON data", nil)
			w.Write([]byte(res))
			return
		}
		acct, err := acctService.UpdateAccount("id", id, data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
		res, _ := templates.Template(w, http.StatusOK, "Account successfully updated", acct)
		w.Write([]byte(res))
	default:
		res, _ := templates.Template(w, http.StatusMethodNotAllowed, "Accepts only GET and PUT requests", nil)
		w.Write([]byte(res))
	}
}