package middleware

import (
	"net/http"
	"strings"

	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/config"
	"github.com/golang-jwt/jwt"
)

func UserGuard(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rToken := r.Header.Get("Authorization")
		if rToken == "" {
			res, _ := templates.Template(http.StatusUnauthorized, "Missing token", nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		bearerToken := strings.Split(rToken, " ")
		if len(bearerToken) != 2 {
			res, _ := templates.Template(http.StatusUnauthorized, "Invalid token", nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.GetTokenSecret()), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			groupName := claims["group_name"]
			if groupName != "admin" && groupName != "user" && groupName != "staff" {
				res, _ := templates.Template(http.StatusUnauthorized, "You are not authorized to access this resource", nil)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(res))
				return
			}
		} else {
			res, _ := templates.Template(http.StatusUnauthorized, err.Error(), nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		handler.ServeHTTP(w, r)
    })
}

func StaffGuard(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rToken := r.Header.Get("Authorization")
		if rToken == "" {
			res, _ := templates.Template(http.StatusUnauthorized, "Missing token", nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		bearerToken := strings.Split(rToken, " ")
		if len(bearerToken) != 2 {
			res, _ := templates.Template(http.StatusUnauthorized, "Invalid token", nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.GetTokenSecret()), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			groupName := claims["group_name"]
			if groupName != "staff" {
				res, _ := templates.Template(http.StatusUnauthorized, "You are not authorized to access this resource", nil)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(res))
				return
			}
		} else {
			res, _ := templates.Template(http.StatusUnauthorized, err.Error(), nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		handler.ServeHTTP(w, r)
    })
}

func AdminGuard(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rToken := r.Header.Get("Authorization")
		if rToken == "" {
			res, _ := templates.Template(http.StatusUnauthorized, "Missing token", nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		bearerToken := strings.Split(rToken, " ")
		if len(bearerToken) != 2 {
			res, _ := templates.Template(http.StatusUnauthorized, "Invalid token", nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.GetTokenSecret()), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			groupName := claims["group_name"]
			if groupName != "admin" {
				res, _ := templates.Template(http.StatusUnauthorized, "You are not authorized to access this resource", nil)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(res))
				return
			}
		} else {
			res, _ := templates.Template(http.StatusUnauthorized, err.Error(), nil)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(res))
			return
		}

		handler.ServeHTTP(w, r)
    })
}