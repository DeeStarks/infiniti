package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/utils"
	"github.com/gorilla/mux"
)

func (h *Handler) AccountTypes(w http.ResponseWriter, r *http.Request) {
	service := h.appPort.NewAccountTypeService()

	switch r.Method {
	case http.MethodGet:
		acctTypes, err := service.ListAccountTypes()
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
	
		res, _ := templates.Template(w, http.StatusOK, "Account types successfully retrieved", acctTypes)
		w.Write([]byte(res))
	case http.MethodPost:
		routeGroup := utils.GetRouteGroup(r)
		if routeGroup != "admin" {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Admin only endpoint", nil)
			w.Write([]byte(res))
			return
		}

		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			res, _ := templates.Template(w, http.StatusBadRequest, "Invalid JSON data", nil)
			w.Write([]byte(res))
			return
		}
		
		acctType, err := service.CreateAccountType(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		res, _ := templates.Template(w, http.StatusCreated, "Account type successfully created", acctType)
		w.Write([]byte(res))
	default:
		tmp, _ := templates.Template(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		w.Write([]byte(tmp))
	}
}

func (h *Handler) SingleAccountType(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		tmp, _ := templates.Template(w, http.StatusBadRequest, "Missing \"id\" parameter", nil)
		w.Write([]byte(tmp))
		return
	}

	acctTService := h.appPort.NewAccountTypeService()
	switch r.Method {
	case http.MethodGet:
		acctT, err := acctTService.GetAccountType("id", id)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Account type successfully retrieved", acctT)
		w.Write([]byte(res))
	case http.MethodPut:
		rg := utils.GetRouteGroup(r)
		if rg != "admin" {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Admin only endpoint", nil)
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
		acctT, err := acctTService.UpdateAccountType("id", id, data)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Account type successfully updated", acctT)
		w.Write([]byte(res))
	case http.MethodDelete:
		rg := utils.GetRouteGroup(r)
		if rg != "admin" {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Admin only endpoint", nil)
			w.Write([]byte(res))
			return
		}

		err := acctTService.DeleteAccountType("id", id)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Account type successfully deleted", nil)
		w.Write([]byte(res))
	}
}