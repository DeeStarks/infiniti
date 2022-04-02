package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/config"
	"github.com/deestarks/infiniti/utils"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId  			interface{}	`json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userId interface{}) (string, error) {
	claims := Claims{
		UserId: 		userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 	15000,
			Issuer:   	"Infiniti",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetTokenSecret()))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { // Only accept POST requests
		res, _ := templates.Template(http.StatusMethodNotAllowed, "Accepts only POST requests", nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(res))
		return
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		res, _ := templates.Template(http.StatusBadRequest, "Invalid JSON data", nil)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(res))
		return
	}

	user, err := h.appPort.NewUserService().CreateUser(data)
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(err.StatusCode(), err.Error(), nil)
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(res))
		return
	}

	// Create a JWT token
	token, err := GenerateToken(user["id"])
	if err != nil {
		res, _ := templates.Template(http.StatusInternalServerError, "Error generating token", nil)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(res))
		return
	}
	newData := make(map[string]interface{})
	newData["token"] = token
	newData["user"] = user
	res, _ := templates.Template(http.StatusOK, "User successfully created", newData)
	w.Write([]byte(res))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { // Only accept POST requests
		res, _ := templates.Template(http.StatusMethodNotAllowed, "Accepts only POST requests", nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(res))
		return
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		res, _ := templates.Template(http.StatusBadRequest, "Invalid JSON data", nil)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(res))
		return
	}

	user, err := h.appPort.NewUserAuthService().AuthenticateUser(data)
	if err, ok := err.(*utils.RequestError); ok {
		res, _ := templates.Template(err.StatusCode(), err.Error(), nil)
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(res))
		return
	}

	// Create a JWT token
	token, err := GenerateToken(user.Id)
	if err != nil {
		res, _ := templates.Template(http.StatusInternalServerError, "Error generating token", nil)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(res))
		return
	}
	newData := make(map[string]interface{})
	newData["token"] = token
	newData["user"] = user
	res, _ := templates.Template(http.StatusOK, "User successfully logged in", newData)
	w.Write([]byte(res))
}