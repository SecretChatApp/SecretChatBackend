package middlewares

import (
	"backend/app/config"
	"backend/app/helpers"
	"context"

	// "context"

	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{
					"message": "User unauthorized",
				}

				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		tokenString := c.Value
		claims := &config.JWTClaim{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			var response = make(map[string]string)
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				response["message"] = "User unauthorized"
			case jwt.ValidationErrorExpired:
				response["message"] = "Token expired"
			default:
				response["message"] = "User unauthorized"
			}

			if len(response) > 0 {
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		if !token.Valid {
			response := map[string]string{
				"message": "User unauthorized",
			}

			helpers.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		ctx := context.WithValue(r.Context(), "userInfo", claims.Email)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
