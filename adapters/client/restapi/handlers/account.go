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
	vars := mux.Vars(r) // Get the URL variables
	var idVar string

	// Get route group (user, staff, admin)
	routeGroup := utils.GetRouteGroup(r)

	switch routeGroup {
	case "user", "staff":
		// Get the user ID from the request context
		id := r.Context().Value(constants.CTXKey("user_id"))
		if id == nil {
			res, _ := templates.Template(w, http.StatusUnauthorized, "User not logged in", nil)
			w.Write([]byte(res))
			return
		}
		idVar = fmt.Sprintf("%v", id)
	case "admin":
		id, ok := vars["id"] // Get the id variable
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

	switch r.Method {
	case "GET":
		acct, err := h.appPort.NewAccountService().GetAccount("id", id, []string{})
		fmt.Println(err)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
	
		res, _ := templates.Template(w, http.StatusOK, "Account successfully retrieved", acct)
		w.Write([]byte(res))
	case "PUT":
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			res, _ := templates.Template(w, http.StatusBadRequest, "Invalid JSON data", nil)
			w.Write([]byte(res))
			return
		}
		fmt.Println(data)
		// TODO: Validate data
		
	default:
		res, _ := templates.Template(w, http.StatusMethodNotAllowed, "Accepts only GET and PUT requests", nil)
		w.Write([]byte(res))
	}
}