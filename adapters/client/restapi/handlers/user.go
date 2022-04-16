package handlers

import (
	"net/http"
	// "strconv"

	// "github.com/deestarks/infiniti/adapters/client/restapi/constants"
	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/utils"
	// "github.com/gorilla/mux"
)

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.appPort.NewUserService().ListUsers()
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
		w.Write([]byte(res))
		return
	}

	res, _ := templates.Template(w, http.StatusOK, "Users successfully retrieved", users)
	w.Write([]byte(res))
}

func (h *Handler) ListStaff(w http.ResponseWriter, r *http.Request) {
	staff, err := h.appPort.NewStaffService().ListStaff()
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
		w.Write([]byte(res))
		return
	}

	res, _ := templates.Template(w, http.StatusOK, "Staff successfully retrieved", staff)
	w.Write([]byte(res))
}

func (h *Handler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	admins, err := h.appPort.NewAdminService().ListAdmins()
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
		w.Write([]byte(res))
		return
	}

	res, _ := templates.Template(w, http.StatusOK, "Admins successfully retrieved", admins)
	w.Write([]byte(res))
}