package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/utils"
	"github.com/gorilla/mux"
)

func (h *Handler) Currencies(w http.ResponseWriter, r *http.Request) {
	service := h.appPort.NewCurrencyService()

	switch r.Method {
	case http.MethodGet:
		currencies, err := service.ListCurrencies()
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
	
		res, _ := templates.Template(w, http.StatusOK, "Currencies successfully retrieved", currencies)
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
		
		currency, err := service.CreateCurrency(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		res, _ := templates.Template(w, http.StatusCreated, "Currency successfully created", currency)
		w.Write([]byte(res))
	default:
		tmp, _ := templates.Template(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		w.Write([]byte(tmp))
	}
}

func (h *Handler) SingleCurrency(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		tmp, _ := templates.Template(w, http.StatusBadRequest, "Missing \"id\" parameter", nil)
		w.Write([]byte(tmp))
		return
	}

	currService := h.appPort.NewCurrencyService()
	switch r.Method {
	case http.MethodGet:
		acctT, err := currService.GetCurrency("id", id)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Currency successfully retrieved", acctT)
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
		currency, err := currService.UpdateCurrency("id", id, data)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Currency successfully updated", currency)
		w.Write([]byte(res))
	case http.MethodDelete:
		rg := utils.GetRouteGroup(r)
		if rg != "admin" {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Admin only endpoint", nil)
			w.Write([]byte(res))
			return
		}

		err := currService.DeleteCurrency("id", id)
		if err, ok := err.(*utils.RequestError); ok {
			tmp, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(tmp))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Currency successfully deleted", nil)
		w.Write([]byte(res))
	}
}