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

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	users, err := h.appPort.NewUserService().ListUsers()
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
		w.Write([]byte(res))
		return
	}

	res, _ := templates.Template(w, http.StatusOK, "Users successfully retrieved", users)
	w.Write([]byte(res))
}

func (h *Handler) Staff(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
	case http.MethodGet:
		staff, err := h.appPort.NewStaffService().ListStaff()
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
	
		res, _ := templates.Template(w, http.StatusOK, "Staff successfully retrieved", staff)
		w.Write([]byte(res))
	case http.MethodPost:
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			res, _ := templates.Template(w, http.StatusBadRequest, "Invalid JSON data", nil)
			w.Write([]byte(res))
			return
		}

		admin, err := h.appPort.NewStaffService().CreateStaff(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		res, _ := templates.Template(w, http.StatusCreated, "Staff successfully created", admin)
		w.Write([]byte(res))
	default:
		res, _ := templates.Template(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		w.Write([]byte(res))
	}
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		admins, err := h.appPort.NewAdminService().ListAdmins()
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		res, _ := templates.Template(w, http.StatusOK, "Admins successfully retrieved", admins)
		w.Write([]byte(res))
	case http.MethodPost:
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			res, _ := templates.Template(w, http.StatusBadRequest, "Invalid JSON data", nil)
			w.Write([]byte(res))
			return
		}

		admin, err := h.appPort.NewAdminService().CreateAdmin(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		res, _ := templates.Template(w, http.StatusCreated, "Admin successfully created", admin)
		w.Write([]byte(res))
	default:
		res, _ := templates.Template(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		w.Write([]byte(res))
	}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	var idVar float64

	// Get route group (user, staff, admin)
	routeGroup := utils.GetRouteGroup(r)
	groups := map[string]bool{
		"user":  true,
		"staff": true,
		"admin": true,
	}

	if !groups[routeGroup] {
		res, _ := templates.Template(w, http.StatusBadRequest, "Invalid url", nil)
		w.Write([]byte(res))
		return
	}

	// Get the user ID from the request context
	ctxId := r.Context().Value(constants.CTXKey("user_id"))
	if ctxId == nil {
		res, _ := templates.Template(w, http.StatusBadRequest, "Unauthenticated", nil)
		w.Write([]byte(res))
		return
	}

	idVar, ok := ctxId.(float64)
	if !ok {
		res, _ := templates.Template(w, http.StatusBadRequest, "Invalid user id. Must be a number", nil)
		w.Write([]byte(res))
		return
	}
	
	id, _ := strconv.Atoi(fmt.Sprintf("%v", idVar)) // Convert the id variable to an int

	switch r.Method {
	case http.MethodGet:
		switch routeGroup {
		case "user":
			service := h.appPort.NewUserService()
			profile, err := service.GetUser("id", id)
			if err, ok := err.(*utils.RequestError); ok {
				res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
				w.Write([]byte(res))
				return
			}
		
			res, _ := templates.Template(w, http.StatusOK, "User Profile successfully retrieved", profile)
			w.Write([]byte(res))
		case "staff":
			service := h.appPort.NewStaffService()
			profile, err := service.GetStaff("id", id)
			if err, ok := err.(*utils.RequestError); ok {
				res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
				w.Write([]byte(res))
				return
			}

			res, _ := templates.Template(w, http.StatusOK, "Staff Profile successfully retrieved", profile)
			w.Write([]byte(res))
		case "admin":
			service := h.appPort.NewAdminService()
			profile, err := service.GetAdmin("id", id)
			if err, ok := err.(*utils.RequestError); ok {
				res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
				w.Write([]byte(res))
				return
			}

			res, _ := templates.Template(w, http.StatusOK, "Admin Profile successfully retrieved", profile)
			w.Write([]byte(res))
		}
	case http.MethodPut:
		data := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			res, _ := templates.Template(w, http.StatusBadRequest, "Invalid JSON data", nil)
			w.Write([]byte(res))
			return
		}

		switch routeGroup {
		case "user":
			service := h.appPort.NewUserService()
			user, err := service.UpdateUser("id", id, data)
			if err, ok := err.(*utils.RequestError); ok {
				res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
				w.Write([]byte(res))
				return
			}
			
			res, _ := templates.Template(w, http.StatusOK, "User Profile successfully updated", user)
			w.Write([]byte(res))
		case "staff":
			service := h.appPort.NewStaffService()
			staff, err := service.UpdateStaff("id", id, data)
			if err, ok := err.(*utils.RequestError); ok {
				res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
				w.Write([]byte(res))
				return
			}

			res, _ := templates.Template(w, http.StatusOK, "Staff Profile successfully updated", staff)
			w.Write([]byte(res))
		case "admin":
			service := h.appPort.NewAdminService()
			admin, err := service.UpdateAdmin("id", id, data)
			if err, ok := err.(*utils.RequestError); ok {
				res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
				w.Write([]byte(res))
				return
			}

			res, _ := templates.Template(w, http.StatusOK, "Admin Profile successfully updated", admin)
			w.Write([]byte(res))
		}
	default:
		res, _ := templates.Template(w, http.StatusMethodNotAllowed, "Accepts only GET and PUT requests", nil)
		w.Write([]byte(res))
	}
}