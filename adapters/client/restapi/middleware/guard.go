package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/deestarks/infiniti/adapters/client/restapi/constants"
	"github.com/deestarks/infiniti/adapters/client/restapi/handlers/templates"
	"github.com/deestarks/infiniti/config"
	"github.com/golang-jwt/jwt"
)

func UserGuard(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rToken := r.Header.Get("Authorization")
		if rToken == "" {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Missing token", nil)
			w.Write([]byte(res))
			return
		}

		bearerToken := strings.Split(rToken, " ")
		if len(bearerToken) != 2 {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Invalid token", nil)
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
			if groupName != "user" {
				res, _ := templates.Template(w, http.StatusUnauthorized, "You are not authorized to access this resource", nil)
				w.Write([]byte(res))
				return
			}

			// Add user id and group name to request

			// There's a warning - "should not use built-in type string as key for value; define your own type to avoid collisions"
			// hence the use of "constants.CTXKey" (which is a custom type)
			ctx := context.WithValue(r.Context(), constants.CTXKey("user_id"), claims["user_id"])
			ctx = context.WithValue(ctx, constants.CTXKey("group_name"), claims["group_name"])
			r = r.WithContext(ctx)
		} else {
			res, _ := templates.Template(w, http.StatusUnauthorized, err.Error(), nil)
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
			res, _ := templates.Template(w, http.StatusUnauthorized, "Missing token", nil)
			w.Write([]byte(res))
			return
		}

		bearerToken := strings.Split(rToken, " ")
		if len(bearerToken) != 2 {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Invalid token", nil)
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
				res, _ := templates.Template(w, http.StatusUnauthorized, "You are not authorized to access this resource", nil)
				w.Write([]byte(res))
				return
			}

			// Add staff id and group name to request

			// There's a warning - "should not use built-in type string as key for value; define your own type to avoid collisions"
			// hence the use of constants.CTXKey (which is a custom type)
			ctx := context.WithValue(r.Context(), constants.CTXKey("user_id"), claims["user_id"])
			ctx = context.WithValue(ctx, constants.CTXKey("group_name"), claims["group_name"])
			r = r.WithContext(ctx)
		} else {
			res, _ := templates.Template(w, http.StatusUnauthorized, err.Error(), nil)
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
			res, _ := templates.Template(w, http.StatusUnauthorized, "Missing token", nil)
			w.Write([]byte(res))
			return
		}

		bearerToken := strings.Split(rToken, " ")
		if len(bearerToken) != 2 {
			res, _ := templates.Template(w, http.StatusUnauthorized, "Invalid token", nil)
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
				res, _ := templates.Template(w, http.StatusUnauthorized, "You are not authorized to access this resource", nil)
				w.Write([]byte(res))
				return
			}

			// Add admin id and group name to request

			// There's a warning - "should not use built-in type string as key for value; define your own type to avoid collisions"
			// hence the use of constants.CTXKey (which is a custom type)
			ctx := context.WithValue(r.Context(), constants.CTXKey("user_id"), claims["user_id"])
			ctx = context.WithValue(ctx, constants.CTXKey("group_name"), claims["group_name"])
			r = r.WithContext(ctx)
		} else {
			res, _ := templates.Template(w, http.StatusUnauthorized, err.Error(), nil)
			w.Write([]byte(res))
			return
		}

		handler.ServeHTTP(w, r)
    })
}