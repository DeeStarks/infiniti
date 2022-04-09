package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/config"
	"github.com/deestarks/infiniti/utils"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId  			interface{}	`json:"user_id"`
	UserGroupName 		interface{}	`json:"group_name"`
	jwt.StandardClaims
}

func GenerateToken(userId, userGroupName interface{}) (string, error) {
	claims := Claims{
		UserId: 		userId,
		UserGroupName: 	userGroupName,
		StandardClaims: jwt.StandardClaims{
			Issuer:   	"Infiniti",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetTokenSecret()))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Only accept POST requests
		res, _ := templates.Template(w, http.StatusMethodNotAllowed, "Accepts only POST requests", nil)
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

	user, err := h.appPort.NewUserService().CreateUser(data)
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
		w.Write([]byte(res))
		return
	}

	// Create a JWT token
	token, err := GenerateToken(user.Id, user.Group.Name)
	if err != nil {
		res, _ := templates.Template(w, http.StatusInternalServerError, "Error generating token", nil)
		w.Write([]byte(res))
		return
	}
	newData := make(map[string]interface{})
	newData["token"] = token
	newData["user"] = user
	res, _ := templates.Template(w, http.StatusOK, "User successfully created", newData)
	w.Write([]byte(res))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Only accept POST requests
		res, _ := templates.Template(w, http.StatusMethodNotAllowed, "Accepts only POST requests", nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
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

	// Check path to get the subroute (user, staff, admin)
	routeGroup := strings.Split(r.URL.Path, "/")[3]

	switch routeGroup {
	case "user":
		user, err := h.appPort.NewUserAuthService().AuthenticateUser(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
	
		// Create a JWT token
		token, err := GenerateToken(user.Id, "user")
		if err != nil {
			res, _ := templates.Template(w, http.StatusInternalServerError, "Error generating token", nil)
			w.Write([]byte(res))
			return
		}
		newData := make(map[string]interface{})
		newData["token"] = token
		newData["user"] = user
		res, _ := templates.Template(w, http.StatusOK, "User successfully logged in", newData)
		w.Write([]byte(res))
		
	case "staff":
		staff, err := h.appPort.NewUserAuthService().AuthenticateStaff(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}
	
		// Create a JWT token
		token, err := GenerateToken(staff.Id, "staff")
		if err != nil {
			res, _ := templates.Template(w, http.StatusInternalServerError, "Error generating token", nil)
			w.Write([]byte(res))
			return
		}
		newData := make(map[string]interface{})
		newData["token"] = token
		newData["staff"] = staff
		res, _ := templates.Template(w, http.StatusOK, "User successfully logged in", newData)
		w.Write([]byte(res))

	case "admin":
		admin, err := h.appPort.NewUserAuthService().AuthenticateAdmin(data)
		if err, ok := err.(*utils.RequestError); ok {
			res, _ := templates.Template(w, err.StatusCode(), err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		// Create a JWT token
		token, err := GenerateToken(admin.Id, "admin")
		if err != nil {
			res, _ := templates.Template(w, http.StatusInternalServerError, "Error generating token", nil)
			w.Write([]byte(res))
			return
		}
		newData := make(map[string]interface{})
		newData["token"] = token
		newData["admin"] = admin
		res, _ := templates.Template(w, http.StatusOK, "Admin successfully logged in", newData)
		w.Write([]byte(res))
	default:
		res, _ := templates.Template(w, http.StatusBadRequest, "Invalid route", nil)
		w.Write([]byte(res))
		return
	}
}