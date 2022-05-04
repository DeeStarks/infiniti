package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/utils"
	"github.com/gorilla/mux"
)

func (h *Handler) TransactionTypes(w http.ResponseWriter, r *http.Request) {
	service := h.appPort.NewTransactionTypesService()

	switch r.Method {
	case http.MethodGet:
		transTypes, err := service.ListTransactionTypes()
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
	
		res, _ := templates.Template(w, http.StatusOK, "Transaction types successfully retrieved", transTypes)
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
		
		transType, err := service.CreateTransactionType(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		res, _ := templates.Template(w, http.StatusCreated, "Transaction type successfully created", transType)
		w.Write([]byte(res))
	default:
		tmp, _ := templates.Template(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		w.Write([]byte(tmp))
	}
}

func (h *Handler) SingleTransactionType(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		tmp, _ := templates.Template(w, http.StatusBadRequest, "Missing \"id\" parameter", nil)
		w.Write([]byte(tmp))
		return
	}

	transTypeService := h.appPort.NewTransactionTypesService()
	switch r.Method {
	case http.MethodGet:
		transType, err := transTypeService.GetTransactionType("id", id)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Transaction type successfully retrieved", transType)
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
		transType, err := transTypeService.UpdateTransactionType("id", id, data)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Transaction type successfully updated", transType)
		w.Write([]byte(res))
	case http.MethodDelete:
		rg := utils.GetRouteGroup(r)
		if rg != "admin" {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Admin only endpoint", nil)
			w.Write([]byte(res))
			return
		}

		err := transTypeService.DeleteTransactionType("id", id)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Transaction type successfully deleted", nil)
		w.Write([]byte(res))
	}
}