package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/config"
	"github.com/deestarks/infiniti/utils"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId  			int		`json:"user_id"`
	UserGroupName 		string	`json:"group_name"`
	jwt.StandardClaims
}

func GenerateToken(userId int, userGroupName string, expiresAt int64) (string, error) {
	claims := Claims{
		UserId: 		userId,
		UserGroupName: 	userGroupName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 	expiresAt, // Token expires in 48 hours
			Issuer:   	"Infiniti",
			IssuedAt:  	time.Now().Unix(),
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
	expiresAt := time.Now().Add(time.Hour * 48).Unix()
	token, err := GenerateToken(user.Id, user.Group.Name, expiresAt)
	if err != nil {
		res, _ := templates.Template(w, http.StatusInternalServerError, "Error generating token", nil)
		w.Write([]byte(res))
		return
	}
	newData := make(map[string]interface{})
	newData["token"] = token
	newData["token_expires_at"] = expiresAt
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
		expiresAt := time.Now().Add(time.Hour * 48)
		token, err := GenerateToken(user.Id, "user", expiresAt.Unix())
		if err != nil {
			res, _ := templates.Template(w, http.StatusInternalServerError, "Error generating token", nil)
			w.Write([]byte(res))
			return
		}
		newData := make(map[string]interface{})
		newData["token"] = token
		newData["token_expires_at"] = expiresAt
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
		expiresAt := time.Now().Add(time.Hour * 48)
		token, err := GenerateToken(staff.Id, "staff", expiresAt.Unix())
		if err != nil {
			res, _ := templates.Template(w, http.StatusInternalServerError, "Error generating token", nil)
			w.Write([]byte(res))
			return
		}
		newData := make(map[string]interface{})
		newData["token"] = token
		newData["token_expires_at"] = expiresAt
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
		expiresAt := time.Now().Add(time.Hour * 48)
		token, err := GenerateToken(admin.Id, "admin", expiresAt.Unix())
		if err != nil {
			res, _ := templates.Template(w, http.StatusInternalServerError, "Error generating token", nil)
			w.Write([]byte(res))
			return
		}
		newData := make(map[string]interface{})
		newData["token"] = token
		newData["token_expires_at"] = expiresAt
		newData["admin"] = admin
		res, _ := templates.Template(w, http.StatusOK, "Admin successfully logged in", newData)
		w.Write([]byte(res))
	default:
		res, _ := templates.Template(w, http.StatusBadRequest, "Invalid route", nil)
		w.Write([]byte(res))
		return
	}
}